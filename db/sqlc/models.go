// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"
)

type Accounts struct {
	Id           int64         `json:"Id"`
	Owner        string        `json:"owner"`
	Balance      int64         `json:"balance"`
	Currency     string        `json:"currency"`
	CreationTime time.Time     `json:"creationTime"`
	CountryCode  sql.NullInt32 `json:"countryCode"`
}

type Entries struct {
	Id           int64         `json:"Id"`
	AccountId    sql.NullInt64 `json:"accountId"`
	Amount       int64         `json:"amount"`
	CreationTime time.Time     `json:"creationTime"`
}

type Transfers struct {
	Id           int64         `json:"Id"`
	SenderId     sql.NullInt64 `json:"senderId"`
	RecipientId  sql.NullInt64 `json:"recipientId"`
	Amount       int64         `json:"amount"`
	CreationTime time.Time     `json:"creationTime"`
}