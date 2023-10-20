package model

type UserItem struct {
	id       int64  `reindex:"id,hash,pk"`
	login    string `reindex:"login,hash"`
	password string `reindex:"password,hash"`
	ip       string `reindex:"ip,hash"`
	date     int64  `reindex:"date,tree"`
}
