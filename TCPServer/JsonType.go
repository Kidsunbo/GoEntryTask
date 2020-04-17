package TCPServer

import "awesomeProject/UserInfo"

type JsonType struct{
	Type string `json:"type"`
}

type JsonLogin struct {
	JsonType
	Username string `json:"username"`
	Password string `json:"password"`
}


type JsonLoginRes struct {
	JsonType
	Result string `json:"result"`
	Reason string `json:"reason"`
	UserInfo UserInfo.UserInfo `json:"user_info"`
}

type JsonUpdate struct {
	JsonType
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Profile string `json:"profile"`
}

type JsonUpdateRes struct {
	JsonType
	Result string `json:"result"`
	Reason string `json:"reason"`
}


