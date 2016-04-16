package main

import(
	"encoding/gob"
    "fmt"
    "net"
    "bufio"
    "os"
)

type Message struct{
	Dst string
	Src string
	Kind string
	Data string	
}

func sendOneMsg(dest string, src string, kind string, data string) {
	//fmt.Println("start client");
	msg := &Message{dest, src, kind, data}
    conn, err := net.Dial("tcp", msg.Dst)
    if err != nil {
        fmt.Println("Connection error: ", err)
    }
    encoder := gob.NewEncoder(conn)
    encoder.Encode(msg)
   
   /*connbuf := bufio.NewReader(conn)
	for{
    	str, err := connbuf.ReadString('\n')
    	if len(str)>0{
        	fmt.Println(str)
    	}
    	if err!= nil {
        	break
    	}
	}*/
    conn.Close()
    fmt.Println("done");
}

func requestHandler(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	msg := &Message{}
	dec.Decode(msg)
	
	fmt.Printf("Received:\n %+v\n", msg);
	/*switch msg.Kind {
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
	}*/	
}

func listeningMsg() {
	fmt.Println("listening messages...")
	socket, err := net.Listen("tcp", ":5000")
  	if err != nil { 
  		fmt.Println("tcp listen error") 
  	} 
  	for {
    	conn, err := socket.Accept()
    	if err != nil { 
    		fmt.Println("connection error") 
    	}
    	go requestHandler(conn)
  	}
}

func main() {
	fmt.Println("hello, press enter sending message")
	go listeningMsg()	
	for{
		//fmt.Print("1 for send message:")
		reader := bufio.NewReader(os.Stdin)
		data, _ := reader.ReadString('\n')
		if data != "" {
			fmt.Print("send\n")
			sendOneMsg("localhost:5000","Myself","test", "this is a test")
		}
	}
}