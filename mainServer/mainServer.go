package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"encoding/csv"
	"strconv"
	"os"
	"io"
)

var groups []Group
var group_id int //assign group id
//var groupServerMap map[string][]string //groupId <-> serverList

func requestHandler(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	msg := &Message{}
	dec.Decode(msg)
	//fmt.Printf("Received : %+v", msg.Data);
	switch msg.Kind {
		case "create_group": createGroup(msg)
		case "join_group": 
			serverList := joinGroup(msg)
			for i := range serverList {
    			conn.Write([]byte(serverList[i]))
    			conn.Write([]byte(" "))
    		}
			conn.Write([]byte("\n"))
		case "remove_server": removeServer(msg)
		//case "group_list": groupList(msg)
	}
	
}

func createGroup(msg *Message) {
	var newGroup Group
	newGroup.id = strconv.Itoa(group_id)
	group_id++
	newGroup.name = msg.Data
	newGroup.serverList = make(map[string]bool)
	newGroup.addServer(msg.Src)
	groups = append(groups, newGroup)
	fmt.Println(groups)
}

func joinGroup(msg *Message) []string {
	var serverSlice []string
	for i:=0; i<len(groups); i++ {
		if groups[i].id == msg.Data {
			groups[i].addServer(msg.Src)
			for k := range groups[i].serverList {
				serverSlice = append(serverSlice, k)
			}	
		}
	}
	fmt.Println(groups)
	return serverSlice
}

func removeServer(msg *Message) {
	for i:=0; i<len(groups); i++ {
		if groups[i].id == msg.Data {
			groups[i].delServer(msg.Src)
			fmt.Println("delete: ", msg.Src)
		}
	}
	fmt.Println(groups)
}

func readSubServerConfig(){
	file, err:= os.Open("/Users/rachel/Documents/goworkspace/musicsync/src/server/initGroups.csv")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()
	
	reader := csv.NewReader(file)	
	for {
		record, err := reader.Read()
		if err == io.EOF {
	    	break
		} else if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		var newGroup Group
		newGroup.id = record[0]
		newGroup.name = record[1]
		newGroup.serverList = make(map[string]bool)
		newGroup.addServer(record[2])
		groups = append(groups, newGroup)
		group_id++
    }
	fmt.Println(groups)
	
}


func initCServer(){
	fmt.Println("hello, I am music sync centre server") 
	group_id = 0
	//groups = make([]string)
	//groupServerMap = make(map[string][]string)
	readSubServerConfig()
}

func main() {
	fmt.Println("Launching server...")
	initCServer()
	socket, err := net.Listen("tcp", ":5000")
  	if err != nil { fmt.Println("tcp listen error") }
 
  	for {
    	conn, err := socket.Accept()
    	if err != nil { 
    		fmt.Println("connection error") 
    	}
    	go requestHandler(conn)
    	//go send_data(conn, channel)
  	}
}