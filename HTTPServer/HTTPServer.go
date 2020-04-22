/*
Mainly hold a server to handle the HTTP request
 */
package HTTPServer

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
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

func IsExist(fileAddr string)(bool){
	_,err := os.Stat(fileAddr)
	if err!=nil{
		if os.IsExist(err){
			return true
		}
		return false
	}
	return true
}

func update(writer http.ResponseWriter, request *http.Request) {
	username:=request.PostFormValue("username")
	nickname:=request.PostFormValue("nickname")
	img,handle,_:=request.FormFile("profile")
	if img!=nil && handle!=nil {
		defer img.Close()
		if !IsExist("./static/") {
			err := os.Mkdir("./static/", os.ModePerm)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
		imgType := [...]string{"jpg", "jpeg", "png"}
		for _, imgT := range imgType {
			_ = os.Remove(fmt.Sprintf("./static/%s.%s", username, imgT))
		}
		temp := strings.Split(handle.Filename, ".")
		extension := strings.ToLower(temp[len(temp)-1])
		saveImg(img, fmt.Sprintf("%s.%s", username, extension))
		postUpdate(username, nickname, fmt.Sprintf("%s.%s", username, extension))
	}else{
		profile:=request.PostFormValue("last_profile")
		postUpdate(username, nickname, profile)
	}
	http.Redirect(writer,request, "/", http.StatusFound)
}

func saveImg(img multipart.File,name string){
	file,err :=os.Create(fmt.Sprintf("./static/%s",name))
	if err!=nil{
		log.Println(err.Error())
		return
	}
	_,err =io.Copy(file,img)
	if err!=nil{
		log.Println(err.Error())
		return
	}
	_= file.Close()
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
		p := struct {
			Success     bool
			Username    string
			Nickname    string
			ProfileName string
			Profile     template.URL
		}{}
		imgFile,err := os.Open(fmt.Sprintf("./static/%s",ui.Profile))
		if err==nil {
			buf := make([]byte, 1024*1024)
			count, err := imgFile.Read(buf)
			if err != nil {
				log.Println(err)
			}
			extension := strings.Split(ui.Profile, ".")[len(strings.Split(ui.Profile, "."))-1]
			base64Str := base64.StdEncoding.EncodeToString(buf[:count])
			base64Str = fmt.Sprintf("data:image/%s;base64,%s", extension, base64Str)
			p = struct {
				Success     bool
				Username    string
				Nickname    string
				ProfileName string
				Profile     template.URL
			}{
				true, ui.Username, ui.Nickname, ui.Profile, template.URL(base64Str),
			}
		}else{
			p = struct {
				Success     bool
				Username    string
				Nickname    string
				ProfileName string
				Profile     template.URL
			}{
				true, ui.Username, ui.Nickname, ui.Profile, template.URL(""),
			}
		}
		t,err :=template.ParseFiles("HTTPServer/template/login.html")
		if err!=nil{
			_,_ = fmt.Fprintf(writer,err.Error())
			return
		}
		_=t.Execute(writer,p)
	}
}

