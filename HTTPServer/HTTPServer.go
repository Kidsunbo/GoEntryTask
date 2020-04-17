/*
Mainly hold a server to handle the HTTP request
 */
package HTTPServer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Run(){
	http.HandleFunc("/",index)
	http.HandleFunc("/login/",login)
	http.HandleFunc("/update/",update)
	log.Fatal(http.ListenAndServe(":8080",nil))
}

func index(writer http.ResponseWriter, request *http.Request) {
	t,err :=template.ParseFiles("HTTPServer/template/index.html")
	if err!=nil{
		_,_ =fmt.Fprintf(writer,err.Error())
		return
	}
	_ =t.Execute(writer,nil)
}

func update(writer http.ResponseWriter, request *http.Request) {
	profile:=request.PostFormValue("profile_photo")
	username:=request.PostFormValue("username")
	nickname:=request.PostFormValue("nickname")
	postUpdate(username,nickname,profile)
	http.Redirect(writer,request, "/", http.StatusFound)
}

func login(writer http.ResponseWriter, request *http.Request) {
	username:=request.PostFormValue("username")
	password :=request.PostFormValue("password")

	if ok,reason,ui :=tryToLogin(username,password);!ok{
		p := struct{
			Success bool
			Reason string
		}{
			false,reason,
		}
		t,err :=template.ParseFiles("HTTPServer/template/login.html")
		if err!=nil{
			_,_ = fmt.Fprintf(writer,err.Error())
			return
		}
		_=t.Execute(writer,p)
	}else{
		p := struct{
			Success bool
			Username string
			Nickname string
			Profile template.URL
		}{
			true,ui.Username,ui.Nickname,template.URL(ui.Profile),
		}
		t,err :=template.ParseFiles("HTTPServer/template/login.html")
		if err!=nil{
			_,_ = fmt.Fprintf(writer,err.Error())
			return
		}
		_=t.Execute(writer,p)
	}
}

