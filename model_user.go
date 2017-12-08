package main

import (
	"fmt"
	"net"

	"golang.org/x/crypto/bcrypt"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const createUserSql = "" +
	"CREATE TABLE IF NOT EXISTS `users` (" +
	"  `id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
	"  `username` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT ''," +
	"  `password` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT ''," +
	"  PRIMARY KEY (`id`)" +
	") ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;"

type User struct {
	ID         int64
	Username   string //用户名，*代表匿名用户
	Password   string //bcrypt加密密码, 用`htpasswd -nBb 用户名 密码`生成
	Privileges []Privilege
}

func (u *User) Authorize(ip net.IP, scopes []AuthScope) ([]ResourceActions, error) {
	if err := u.LoadPrivileges(); err != nil {
		return nil, fmt.Errorf("Load privileges failed: %s", err)
	}

	resourceActions := []ResourceActions{}
	for _, scope := range scopes {
		ownedActionMap := map[string]bool{}
		for _, p := range u.Privileges {
			actions := p.AuthorizedActions(scope.Category, scope.RepoName, ip)
			for _, a := range actions {
				ownedActionMap[a] = true
			}
		}

		//Action只需要给出用户已有权限和申请权限的交集即可，不需要验证
		//Refer: https://docs.docker.com/registry/spec/auth/token/
		actions := []string{}
		if len(ownedActionMap) >= 0 {
			for _, a := range scope.Actions {
				_, found := ownedActionMap[a]
				if found {
					actions = append(actions, a)
				}
			}
		}

		resourceAction := ResourceActions{
			Actions: actions,
			Name:    scope.RepoName,
			Type:    scope.Category,
		}
		resourceActions = append(resourceActions, resourceAction)
	}
	return resourceActions, nil
}

func GetUser(username, password string) (User, error) {
	u := User{}

	row := dbConn.QueryRow("SELECT `id`,`username`,`password` FROM `users` WHERE `username`=?", username)
	err := row.Scan(&u.ID, &u.Username, &u.Password)
	if err == sql.ErrNoRows {
		return u, fmt.Errorf("User not exists")
	}

	if err != nil {
		return u, fmt.Errorf("Query error: %s", err)
	}

	//非匿名用户需要检查密码
	if username != "*" {
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			return u, fmt.Errorf("Invalid password: %s", err)
		}
	}

	return u, nil
}

func (u *User) LoadPrivileges() error {
	rows, err := dbConn.Query("SELECT `host`,`action`,`repo`,`category` FROM `privileges` WHERE `user_id`=?", u.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	u.Privileges = make([]Privilege, 0)
	for rows.Next() {
		p := Privilege{}
		err := rows.Scan(&p.Host, &p.Action, &p.Repo, &p.Category)
		if err != nil {
			return err
		}
		u.Privileges = append(u.Privileges, p)
	}

	return rows.Err()
}
