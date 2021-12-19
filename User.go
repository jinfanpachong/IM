package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C chan string  //管道
	conn net.Conn  //conn
	server *Server
}

func NewUser(conn net.Conn,server *Server)*User{
	userAddr:=conn.RemoteAddr().String()
	user:=&User{userAddr,userAddr,make(chan string),conn,server}
    go user.ListerMessage()
	return user

}

// ListerMessage 监听当前User channel 的方法，一旦有消息，就直接发送给客户端
func (this *User)ListerMessage(){

	for{
		msg:= <-this.C

		this.conn.Write([]byte(msg+"\n"))


	}

}

//用户上线了
func (this *User)Online(){

	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name]=this
	this.server.mapLock.Unlock()


	this.server.BroadCast(this,"已上线")
}

//用户下线
func (this *User)Offline(){
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()

	this.server.BroadCast(this,"下线")
}


func (this User)UserWrite(msg string){
	this.conn.Write([]byte(msg))

}

//用户处理消息的业务
func (this *User)DoMessage(msg string){

	if msg=="who"{
		this.server.mapLock.Lock()
		for _,u:=range this.server.OnlineMap{
			newmsg:= "["+u.Addr+"]"+u.Name+":"+"在线"
			this.UserWrite(newmsg)
		}
		this.server.mapLock.Unlock()

	}else if len(msg)>7 && msg[:7]=="rename|"{
		newName := strings.Split(msg,"|")[1]

		//判断name是否存在
		_,ok:=this.server.OnlineMap[newName]
		if ok{
			return
		}else{
		    this.server.mapLock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName]=this
			this.server.mapLock.Unlock()
			this.Name=newName
			this.UserWrite("您已经更新用户名:"+this.Name+"\n")
		}

	}else if len(msg)>4 && msg[:3]=="to|"{
		toname:=strings.Split(msg,"|")[1]
		if toname==""{
			this.UserWrite("消息格式不正确\n")
			return
		}

		tomsg:=strings.Split(msg,"|")[2]
		if v,ok:=this.server.OnlineMap[toname];ok{
			v.UserWrite(tomsg)
		}

	} else{
		this.server.BroadCast(this,msg)
	}

}

