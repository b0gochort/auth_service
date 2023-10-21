package model

import "github.com/restream/reindexer/v3"

type UserItem struct {
	ID            int64           `json:"id" reindex:"id,hash,pk"`
	Name          string          `json:"name" reindex:"name,hash"`
	Surname       string          `json:"surname" reindex:"surname,hash"`
	Patronymic    string          `json:"patronymic" reindex:"patronymic,hash"`
	Email         string          `json:"email" reindex:"email,hash"`
	Authenticated int32           `json:"authenticated" reindex:"authenticated,hash"`
	Login         string          `json:"login" reindex:"login,hash"`
	Password      string          `json:"password" reindex:"password,hash"`
	IP            string          `json:"ip" reindex:"ip,hash"`
	Birthday      int64           `json:"birthday" reindex:"birthday,tree"`
	City          string          `json:"city" reindex:"city,hash"`
	Position      reindexer.Point `json:"position" reindex:"position,rtree"`
	Date          DateType        `json:"date"`
}

type DateType struct {
	Update int64 `json:"update" reindex:"update,tree"`
	Create int64 `json:"create" reindex:"create,tree"`
}

type EmailItem struct {
	Id    int64  `json:"id" reindex:"id,,pk"`
	Code  string `json:"code" reindex:"code,hash"`
	Email string `json:"email" reindex:"email,hash"`
	Time  int64  `json:"time" reindex:"time,ttl,expire_after=3600"`
}
