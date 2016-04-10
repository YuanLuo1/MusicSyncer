package main

import (
	"net"
	"fmt"
 	"bufio"
 	"strings"
	"os"
	"io"
	"log"
	"io/ioutil"
)

func checkFileExist(fileName string) bool{
	fileName = "./test/" + fileName
	if _, err := os.Stat(fileName); err == nil{
		return true
	}
	return false
}

func handleConnection(conn net.Conn){
	fmt.Println("Start handling connection")
	reader := bufio.NewReader(conn)
	request, _ := reader.ReadString('\n')
	request = strings.Trim(request, "\n")
	switch request{
		case "upload":
			fmt.Println("upload file")
			recvUploadFile(conn, reader)
			return
		case "get":
			fmt.Println("user tries to retrieve file")
			sendFile(conn, reader)
			return
		case "list":
			fmt.Println("user tries to retrieve file")
			listFiles(conn, reader)
			return
	}
	fmt.Println("action not valid....")
	conn.Close()
}

func recvUploadFile(conn net.Conn, reader *bufio.Reader){

	// Dirctory
	directory := "./test/"

	fileName, _ := reader.ReadString('\n')
	fileName = strings.Trim(fileName, "\n")
	fmt.Println("Filename: ", fileName)
	// Check if file already exists
	if checkFileExist(fileName){
		fmt.Println("file already exists\n")
		fmt.Fprintf(conn, "File already exists\n")
		return
	}
	fmt.Println("file not exists")
	// send file success
	fmt.Fprintf(conn, "success\n")

	// Wait to read file
	var receivedBytes int64
	// reader := bufio.NewReader(conn)
	f, err := os.Create(directory + fileName)
	defer f.Close()
	if err != nil {
	    fmt.Println("Error creating file")
	}
	receivedBytes, err = io.Copy(f, conn)
	fmt.Println("recvUploadFile succeess!")
	if err != nil {
	    panic("Transmission error")
	}
	fmt.Printf("Finished transferring file. Received: %d \n", receivedBytes)
	conn.Close()
}

func sendFile(conn net.Conn, reader *bufio.Reader){
	directory := "./test/"
	fileName, _ := reader.ReadString('\n')
	fileName = strings.Trim(fileName, "\n")
	fmt.Println("fileName: ", fileName)
	// we don't have that file
	if !checkFileExist(fileName){
		fmt.Fprintf(conn, "No such file\n")
		return
	}
	fmt.Fprintf(conn, "success\n")
	var n int64
	file, err := os.Open(strings.TrimSpace(directory + fileName))
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close()
	n, err = io.Copy(conn, file)
	conn.Close()
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println(n, "bytes sent")
}

func listFiles(conn net.Conn, reader *bufio.Reader){
	directory := "./test"

	var fileList string = ""

	files, _ := ioutil.ReadDir(directory)
    for _, f := range files {
        fmt.Println(f.Name())
        fileList += f.Name() + "; "
    }
    conn.Write([]byte(fileList + "\n"))
    conn.Close()
}

func main() {
	var port string = "9999"

	fmt.Println("Launching Server")
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil{
		fmt.Println("<Error> Can not listen too port!")
		return
	}

	for{
		conn, err := listen.Accept()
		conn.(*net.TCPConn).SetNoDelay(true)
		if err != nil{
			fmt.Println("<Error> Error when connecting to client!")
			continue
		}
		go handleConnection(conn)
	}

}