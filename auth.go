package main

import (
	"fmt"
)

type AuthScope struct {
	RepoName string
	Category string
	Actions  []string
}

func (s AuthScope) String() string {
	return fmt.Sprintf("类别%s 仓库%s 操作%s", s.Category, s.RepoName, s.Actions)
}
