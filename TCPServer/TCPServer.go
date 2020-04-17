/*
I don't know why I need this because it's just a web app, but it is here.
 */
package TCPServer

import (
	"awesomeProject/MySQLPart"
	"awesomeProject/RedisPart"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

type ServerTcp struct {
	port uint16
}

func NewServer(port uint16) *ServerTcp{
	return &ServerTcp{port}
}

func handleConnection(con net.Conn){
	// Do some stuff
	buf :=make([]byte,1024)
	com :=completer{}
	for{
		l,err:=con.Read(buf)
		if err!=nil{
			log.Println(err.Error(),con.RemoteAddr().String(),"close")
			break
		}
		msg:=string(buf[:l])
		com.set(msg)
		data,ok := com.get()
		for ok {
			process(data, con)
			data,ok =com.get()
		}
	}

}


func process(msg string,con net.Conn){
	t:= JsonType{}
	err :=json.Unmarshal([]byte(msg),&t)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	switch t.Type {
	case "login":
		handleLogin(msg,con)
	case "update":
		handleUpdate(msg,con)
	default:
		// Only means that I don't forget default
	}
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

func handleLogin(msg string,con net.Conn){
	login :=&JsonLogin{}
	err :=json.Unmarshal([]byte(msg),login)
	if err!=nil{
		log.Println(err.Error())
		return
	}
	username,password:=login.Username,login.Password
	res :=&JsonLoginRes{}
	res.Type="login"
	if ok,reason :=RedisPart.Authenticate(username,password);ok{
		ui := RedisPart.GetUserInfo(username)
		if ui==nil{
			res.Result="fail"
			res.Reason="No user info found"
			buf,_:=json.Marshal(res)
			sendData(buf,con)
		}else {
			res.Result = "success"
			res.UserInfo = *ui
			buf,_:=json.Marshal(res)
			sendData(buf,con)
		}
	} else{
		res.Result="fail"
		res.Reason=reason
		buf,_:=json.Marshal(res)
		sendData(buf,con)
	}
}

func handleUpdate(msg string,con net.Conn){
	update := new(JsonUpdate)
	err:=json.Unmarshal([]byte(msg),update)
	if err!=nil{
		log.Println(err.Error())
		return
	}
	res :=&JsonLoginRes{}
	res.Type="update"
	if RedisPart.Exist(update.Username) {
		err =MySQLPart.WriteToMysql(update.Username, update.Nickname, update.Profile, MySQLPart.UPDATE)
	}else{
		err =MySQLPart.WriteToMysql(update.Username,update.Nickname,update.Profile,MySQLPart.INSERT)
	}
	if err!=nil{
		res.Result="fail"
		res.Reason=err.Error()
	}else{
		res.Result = "success"
	}
	buf,err := json.Marshal(*res)
	if err!=nil{
		log.Println(err)
		return
	}
	sendData(buf,con)
}


func (s *ServerTcp)Run(){
	go RedisPart.SyncUserInfoWithMySql()

	listener,err :=net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	defer listener.Close()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	for {
		con, err := listener.Accept()
		if err!=nil{
			fmt.Println(err)
			return
		}
		go handleConnection(con)
	}
}