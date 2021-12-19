package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	server string
	port int
	name string
	conn net.Conn
	flag int

}

func NewClient(server string,port int)*Client{

	client:=&Client{server: server,port: port,flag: 999}

	conn,err:=net.Dial("tcp",fmt.Sprintf("%s:%d",server,port))

	if err!=nil{

		fmt.Println("链接失败")
		return nil
	}
	client.conn=conn
	return client

}

func (client *Client)UpdateName()bool{
	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&client.name)

	sendmsg :="rename|"+client.name +"\n"
	_,err:=client.conn.Write([]byte(sendmsg))
	if err!=nil{
		fmt.Println("conn.Write err",err )
		return  false
	}

	return true
}


func (client *Client)DealResponse(){
	io.Copy(os.Stdout,client.conn)



}


func (client Client)SelectUser(){
	sendmsg:="who\n"
	_,err:=client.conn.Write([]byte(sendmsg))
	if err!=nil{
		fmt.Println("conn write is error",err)
		return
	}



}
func (client Client)privateChat(){

	var remoteName string
	var chatMsg string
	client.SelectUser()
	fmt.Println(">>>>>请输入聊天对象[用户名],exit退出:")
	fmt.Scanln(&remoteName)

	for remoteName !="exit"{
		fmt.Println(">>>>>请输入消息内容,exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg!="exit"{

			if len(chatMsg)!=0{
				sendmsg:="to|"+remoteName+"|"+chatMsg+"\n"
				_,err:=client.conn.Write([]byte(sendmsg))
				if err!=nil{
					fmt.Println("conn write is error",err)
					break
				}
			}
			chatMsg=""
			fmt.Println(">>>>>请输入消息内容,exit退出")
			fmt.Scanln(&chatMsg)

		}

		client.SelectUser()
		fmt.Println(">>>>>请输入聊天对象[用户名],exit退出:")
		fmt.Scanln(&remoteName)
	}


}

func (client Client)PublicChat(){
	var  chatmsg string

	fmt.Println(">>>>>请输入聊天内容，exit退出")

	fmt.Scanln(&chatmsg)

	for chatmsg!="exit"{

		if len(chatmsg)!=0{
			sendMsg :=chatmsg +"\n"
			_,err:=client.conn.Write([]byte(sendMsg))
			if err!=nil{
				fmt.Println("conn write is error",err)
				break
			}

		}

		chatmsg=""

		fmt.Println(">>>>>请输入聊天内容，exit退出")

		fmt.Scanln(&chatmsg)
	}




}

func (this *Client)Run(){

	for this.flag!=0{
		for this.Menu()!=true{

		}

		switch  this.flag {
		case 1:
			fmt.Println("公聊模式请选择")
			this.PublicChat()
			break
		case 2:
			fmt.Println("私聊模式请选择")
			this.privateChat()
			break
		case 3:
			fmt.Println("更新用户名...")
			this.UpdateName()
			break
		}
	}



}

func (this *Client)Menu()bool{
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag>=0 && flag<=3{
		this.flag=flag
		return true
	}else{
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}

}
func main() {

	client:=NewClient("127.0.0.1",8890)

	if client==nil{
		fmt.Printf("链接失败")
	}

	go client.DealResponse()
	fmt.Println("链接成功")

	client.Run()





}
