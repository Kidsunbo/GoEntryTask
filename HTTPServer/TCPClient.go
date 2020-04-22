package HTTPServer

import (
	"awesomeProject/TCPServer"
	. "awesomeProject/UserInfo"
	"encoding/json"
	"log"
	"net"
	"strconv"
)

type User struct {
	Username string

}

func sendData(buf []byte,conn net.Conn){
	msg :=string(buf)
	msg = strconv.Itoa(len(msg))+msg
	buf = []byte(msg)
	length := len(buf)
	count :=0

	for count<length{
		temp,err :=conn.Write(buf)
		if err!=nil{
			log.Println(err.Error())
			return
		}
		count +=temp
	}
}

func tryToLogin(username,password string) (bool,string,*UserInfo){
	com :=completer{}
	client,err:=net.Dial("tcp","localhost:12345")
	if err!=nil{
		log.Println(err.Error())
		return false,err.Error(),nil
	}
	defer client.Close()

	l :=TCPServer.JsonLogin{Username: username,Password: password,JsonType:TCPServer.JsonType{Type: "login"}}
	buf,_ := json.Marshal(l)
	sendData(buf,client)
	buf = make([]byte,1024*1024)
	ok :=false
	msg:=""
	for !ok {
		count, err := client.Read(buf)
		if err != nil {
			return false, err.Error(),nil
		}
		com.set(string(buf[:count]))
		msg,ok = com.get()
	}
	ret := &TCPServer.JsonLoginRes{}
	_ =json.Unmarshal([]byte(msg),ret)
	if ret.Result=="fail"{
		return false,ret.Reason,nil
	}
	return true,"",&ret.UserInfo
}


func postUpdate(username,nickname,profile string) bool {
	com :=completer{}
	client,err:=net.Dial("tcp","localhost:12345")
	if err!=nil{
		log.Println(err.Error())
		return false
	}
	defer client.Close()
	l:=TCPServer.JsonUpdate{}
	l.Type="update"
	l.Username = username
	l.Nickname=nickname
	l.Profile=profile
	buf,_ := json.Marshal(l)
	sendData(buf,client)
	buf = make([]byte,1024*1024)
	ok :=false
	msg:=""
	for !ok {
		count, err := client.Read(buf)
		if err != nil {
			break
		}
		com.set(string(buf[:count]))
		msg,ok = com.get()
	}
	ret := &TCPServer.JsonUpdateRes{}
	err =json.Unmarshal([]byte(msg),ret)
	if err!=nil{
		return false
	}
	return ret.Result=="success"
}