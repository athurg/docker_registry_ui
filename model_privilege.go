package main

import (
	"database/sql"
	"log"
	"net"
	"regexp"
	"strings"
)

func InitPrivilegeTable(db *sql.DB) error {
	createSql := "CREATE TABLE IF NOT EXISTS `privileges` ("
	createSql += "  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,"
	createSql += "  `user_id` int(255) NOT NULL,"
	createSql += "  `host` varchar(255) NOT NULL DEFAULT '',"
	createSql += "  `action` varchar(255) NOT NULL DEFAULT '',"
	createSql += "  `repo` varchar(255) NOT NULL DEFAULT '',"
	createSql += "  `category` varchar(255) NOT NULL DEFAULT '',"
	createSql += "  PRIMARY KEY (`id`)"
	createSql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8"
	if _, err := db.Exec(createSql); err != nil {
		return err
	}

	initSql := "INSERT IGNORE INTO `privileges` (`id`, `user_id`, `host`, `action`, `repo`, `category`) VALUES"
	initSql += " (1, 2, '127.0.0.1/32', 'pull,push', '.*', '.*'),"
	initSql += " (2, 2, '0.0.0.0/0', 'pull', '.*', '.*')"
	if _, err := db.Exec(initSql); err != nil {
		return err
	}

	return nil
}

type Privilege struct {
	UserId   int64
	Host     string //IP掩码匹配,多网段逗号分隔
	Action   string //逗号分隔,数组交集
	Repo     string //正则匹配
	Category string //正则匹配
}

func (p *Privilege) AuthorizedActions(category, repo string, ip net.IP) []string {
	//检查Host
	_, ipNet, err := net.ParseCIDR(p.Host)
	if err != nil {
		log.Printf("[ERROR]不合法的用户Host(%s): %s", p.Host, err)
		return nil
	} else if !ipNet.Contains(ip) {
		log.Printf("[DEBUG]来源%s不符合%s", ip, ipNet)
		return nil
	}

	//检查Category
	if ok, err := regexp.MatchString(p.Category, category); !ok || err != nil {
		log.Printf("[DEBUG]类别%s不符合%s", category, p.Category)
		return nil
	}

	//检查Repo
	if ok, err := regexp.MatchString(p.Repo, repo); !ok || err != nil {
		log.Printf("[DEBUG]仓库%s不符合%s", repo, p.Repo)
		return nil
	}

	actions := strings.Split(p.Action, ",")

	log.Printf("[DEBUG]匹配成功,授权操作%s", actions)

	return actions
}
