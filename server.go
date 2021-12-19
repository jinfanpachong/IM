package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip string
	Port int
	OnlineMap map[string]*User
	mapLock sync.RWMutex
	Message chan  string
}


func NewServer(ip string,port int)*Server{
	server:=&Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message:make(chan  string),
	}

	return  server
}
func (this *Server)listenMessager(){
	for {
		msg:=<-this.Message

		this.mapLock.Lock()
		for _,v:=range this.OnlineMap{
			v.C<-msg
		}
		this.mapLock.Unlock()
	}

}


func (this *Server)BroadCast(user *User,msg string){
	sendMsg :="[" + user.Addr +"]" +user.Name +":" +msg

	this.Message <-sendMsg

}


// Handler 启动服务器的接口
func (this *Server)Handler(conn net.Conn){
	user:=NewUser(conn,this)
	user.Online()

	//接收当前客户端发送的消息并广播
	//
	isLive:=make(chan bool)
	go func (){

		buf:=make([]byte,4096)
		for {
			n,err:=conn.Read(buf)
			if n==0{
				user.Offline()

				return
			}

			if err !=nil && err != io.EOF{
				fmt.Println("Conn Read err:",err )
				return
			}
			//去除用户的消息\n
			msg :=string(buf[:n-1])


			user.DoMessage(msg)
            //是否活跃
			isLive<-true
		}
	}()



	//堵塞 不让子hanfler结束
	for{
		select {
		 case <-isLive:

		 case <-time.After(300*time.Second):
		 	user.Offline() //用户下线
		 	close(user.C)
		 	conn.Close()
			 return  //或者runtime.Goexit()

		}
	}



	fmt.Println("链接建立成功")
}



func (this *Server)Start(){
	listener,error:= net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))

	if error !=nil{
		fmt.Println("net.listener err",error)
		return
	}

	defer listener.Close()
    go this.listenMessager()

	for {
		conn,err:=listener.Accept()

		if err !=nil{
			fmt.Println("net.accept err",err)
			continue
		}


		//d hander

		go this.Handler(conn)


	}




}



