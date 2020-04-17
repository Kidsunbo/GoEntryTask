/*
Handle the operation with redis
 */
package RedisPart

import (
	"awesomeProject/MySQLPart"
	"awesomeProject/UserInfo"
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
)


func init(){
	if ok := newRedisClient("<ip and port of my redis server","");!ok{
		log.Println("Redis fails to connect")
	}
	go func() {
		for login := range MySQLPart.MqUsers{
			_ =setUserLogin(login.Username,login.Password)
		}
	}()

}

func Authenticate(username string,password string) (bool,string){
	res,err :=client.Get(fmt.Sprintf("users:%s",username)).Result()
	if err!=nil{
		return false,err.Error()
	}
	if password!=res{
		return false,fmt.Sprintf("Password not match for user %s",username)
	}
	return true,""
}

func GetUserInfo(username string) (ui *UserInfo.UserInfo){
	res,err:=client.HGet(fmt.Sprintf("info:%s",username),"nickname").Result()
	if err!=nil{
		log.Println(err.Error())
		return nil
	}
	pro,err:=client.HGet(fmt.Sprintf("info:%s",username),"profile").Result()
	if err!=nil{
		log.Println(err.Error())
		return nil
	}
	return &UserInfo.UserInfo{Username: username, Nickname: res, Profile: pro}
}

func Exist(username string) bool{
	count,err :=client.Exists(username).Result()
	if err!=nil || count==0{
		return false
	}
	return true
}

func SyncUserInfoWithMySql(){
	for m := range MySQLPart.Mq{
		err:=setUserInfo(m.Username,m.Nickname,m.Profile)
		if err!=nil{
			log.Println(err.Error())
		}
	}
}

func setUserInfo(username, nickname,profile string) error{
	_,err := client.HMSet(fmt.Sprintf("info:%s",username),"username",username,"nickname",nickname,"profile",profile).Result()
	return err
}

func setUserLogin(username,password string)error{
	_,err :=client.Set(fmt.Sprintf("users:%s",username),password,0).Result()
	return err
}

var client *redis.Client = nil

func newRedisClient(addr string,password string) bool{
	client = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: 0,
	})
	_,err :=client.Ping().Result()
	if err!=nil{
		log.Println(err)
		return false
	}
	return true
}