package dto

import db "osvauld/db/sqlc"

type CreateUser struct {
	UserName  string `json:"username"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

type Login struct {
	UserName string `json:"username"`
}

type LoginReturn struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}
