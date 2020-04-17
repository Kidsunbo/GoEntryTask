/*
The user information ,which actually doesn't deserve a package
 */
package UserInfo

type UserInfo struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Profile string `json:"profile"`
}

type UserLogin struct {
	Username string
	Password string
}