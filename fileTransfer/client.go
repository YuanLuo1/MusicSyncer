package main

import (
	"net"
	"os"
	"fmt"
	"log"
 	"bufio"
	"strings"
	"io"
)

func checkFileExist(fileName string) bool{
	if _, err := os.Stat(fileName); err == nil{
		return true
	}
	return false
}

func dealAction(text string, connection net.Conn){
	actions := strings.Fields(text)
	if len(actions) == 0{
		return
	}
	switch actions[0]{
	case "upload":
		fmt.Println("file upload .....")
		if len(actions)!= 2 {
			fmt.Println("<Error Format> Usage: upload fileName")
			return
		}

		sendFile(actions[1], connection)
		return
	case "get":
		if len(actions) != 	2{
			fmt.Println("<Error Format> Usage: upload fileName")
			return
		}
		
		requestFile(actions[1], connection)
		return

	case "list":
		if len(actions) != 	1{
			fmt.Println("<Error Format> Usage: upload fileName")
			return
		}
		listGroupFiles(connection)
		return
	case "exit":
		connection.Close()
		os.Exit(0)
	}
	connection.Close()
	fmt.Println("<Error Action> No action matched")
}

// Trying to get all file list within the group
func listGroupFiles(conn net.Conn){
	// send get list request
	conn.Write([]byte("list\n"))

	// Recv
	fmt.Println("======== List You can download from ===========")
	msg, _ := bufio.NewReader(conn).ReadString('\n')
	files := strings.Split(msg, "; ")
	fmt.Println(files[:len(files)-1])
	fmt.Println("=========== End of List ===========")
}

// Trying to get file from server
func requestFile(fileName string, conn net.Conn){

	// Dircetory -- where file saved
	directory := "./test2/"

	// send action
	conn.Write([]byte("get\n"))
	// send request file name
	conn.Write([]byte(fileName + "\n"))
	// fmt.Fprintf(conn, fileName)
	
	msg, _ := bufio.NewReader(conn).ReadString('\n')
	// if server doesn't have that file || client isn't in the group
	if strings.Compare(msg, "success\n") != 0{
		fmt.Println("<ERROR> ", msg)
		return
	}

	var receivedBytes int64
	// reader := bufio.NewReader(conn)
	f, err := os.Create(directory + fileName)
	defer f.Close()
	if err != nil {
	    fmt.Println("Error creating file")
	}
	receivedBytes, err = io.Copy(f, conn)
	conn.Close()
	if err != nil {
	    panic("Transmission error")
	}

	fmt.Printf("Finished transferring file. Received: %d \n", receivedBytes)
}

// upload file to server or group
func sendFile(fileName string, connection net.Conn){
	if(!checkFileExist(fileName)){
		fmt.Println("<Error> File Not Exist")
		return
	}
	// Send action
	connection.Write([]byte("upload\n"))
	filePath := strings.Split(fileName, "/")
	// Send file name
	connection.Write([]byte(filePath[len(filePath)-1]+"\n"))
	msg, _ := bufio.NewReader(connection).ReadString('\n')
	// if already exists
	if strings.Compare(msg, "success\n") != 0{
		fmt.Println("msg: ", msg)
		fmt.Println("File already exists in server")
		return
	}

	var n int64
	file, err := os.Open(strings.TrimSpace(fileName))
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close() // make sure to close the file even if we panic.
	n, err = io.Copy(connection, file)
	connection.Close()
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println(n, "bytes sent")
}

func main(){
	var ip string = "127.0.0.1"
	var port string = "9999"
	reader := bufio.NewReader(os.Stdin)
	for {
		connection, err  := net.Dial("tcp", ip + ":" + port)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Unable to connect server")
		}
		fmt.Println("Connected to server ....")
		text, _ := reader.ReadString('\n')
		dealAction(text, connection)
		// fmt.Fprintf(connection, text + "\n")
	}
}
