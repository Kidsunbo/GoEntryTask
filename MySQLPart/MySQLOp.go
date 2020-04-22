/*
Handle the operation with MySQL
 */

package MySQLPart

import (
	"awesomeProject/UserInfo"
	"database/sql"
	"fmt"
	"log"
)

import _ "github.com/go-sql-driver/mysql"

type WriteKind int
const (
	INSERT WriteKind = iota
	UPDATE
)


var db *sql.DB =nil
var insert *sql.Stmt
var update *sql.Stmt
func init(){
	db=newMysqlClient("192.168.1.9","Learn","root","root")
	tInsert ,err := db.Prepare("INSERT INTO info VALUES(?,?,?)")
	if err!=nil{
		log.Println(err.Error())
	}
	insert = tInsert
	update,_ = db.Prepare("UPDATE info SET nickname=?,profile=? WHERE username=?")


	rows, err :=db.Query("SELECT * FROM users")
	if err!=nil{
		log.Fatal(err)
	}

	for rows.Next(){
		user :=UserInfo.UserLogin{}
		err :=rows.Scan(&user.Username,&user.Password)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		MqUsers<-user
	}
	close(MqUsers)

	rows,err = db.Query("SELECT * FROM info")
	if err!=nil{
		log.Fatal(err.Error())
	}
	for rows.Next(){
		user:=UserInfo.UserInfo{}
		err :=rows.Scan(&user.Username,&user.Nickname,&user.Profile)
		if err!=nil{
			log.Println(err.Error())
			continue
		}
		Mq<-user
	}

}

var Mq chan UserInfo.UserInfo = make(chan UserInfo.UserInfo,1024)
var MqUsers chan UserInfo.UserLogin = make(chan UserInfo.UserLogin,1024)

func newMysqlClient(ip,db,user,password string) *sql.DB{
	con,err:=sql.Open("mysql",fmt.Sprintf("%s:%s@(%s)/%s",user,password,ip,db))
	if err!=nil{
		log.Println(err.Error())
		return nil
	}
	return con
}


func WriteToMysql(name,nickname,profile string,op WriteKind) error{
	switch op {
	case INSERT:
		_,err := insert.Exec(name,nickname,profile)
		if err!=nil{
			log.Println(err.Error())
			return err
		}
	case UPDATE:
		_,err := update.Exec(nickname,profile,name)
		if err!=nil{
			log.Println(err.Error())
			return err
		}
	default:
		log.Println("Unknown operation")
	}
	Mq<-UserInfo.UserInfo{Username: name, Nickname: nickname, Profile: profile}
	return nil
}
