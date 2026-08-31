package handler

import "encoding/gob"

func init() {
	gob.Register(sessionKey(""))
}

type sessionKey string

const UserIDSessionKey string = "userid"
