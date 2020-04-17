package TCPServer

import (
	"fmt"
	"strconv"
	"strings"
)

type completer struct {
	msg string
}

func (c *completer)set(msg string){
	c.msg+=msg
}

func (c *completer)get() (msg string, ok bool){
	index :=strings.Index(c.msg,"{")
	if index==-1{
		return "",false
	}
	length,err :=strconv.Atoi(c.msg[:index])
	if err!=nil{
		fmt.Println(err.Error())
		return "",false
	}
	if len(c.msg)-index<length{
		return "",false
	}
	msg=c.msg[index:index+length]
	ok=true
	c.msg=c.msg[index+length:]
	return
}