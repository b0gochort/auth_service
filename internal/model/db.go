package model

type UserItem struct {
	ID       int64  `json:"id" reindex:"id,hash,pk"`
	Login    string `json:"login" reindex:"login,hash"`
	Password string `json:"password" reindex:"password,hash"`
	IP       string `json:"ip" reindex:"ip,hash"`
	Date     int64  `json:"date" reindex:"date,tree"`
}
