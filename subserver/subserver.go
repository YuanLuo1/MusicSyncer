package main

import (
	"encoding/gob"
	"net/http"
    "fmt"
    "strings"
    "net"
    "html/template"
    "io"
    "time"
    "crypto/md5"
    "strconv"
    "bufio"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
    	hasher := md5.New()
        io.WriteString(hasher, strconv.FormatInt(time.Now().Unix(), 10))
        token := fmt.Sprintf("%x", hasher.Sum(nil))
    	t, _ := template.ParseFiles("/Users/rachel/Documents/goworkspace/musicsync/src/subserver/create.html")
    	t.Execute(w,token)
    } else if r.Method == "POST" {
    	r.ParseForm()
    	//ip,_,_ := net.SplitHostPort(r.RemoteAddr)
    	groupName := strings.TrimSpace(r.PostFormValue("groupname"))
    	comMainServer(groupName, "create_group")
    	//createGroup(ip, groupName)
    } else {
    	fmt.Fprintf(w, "Error Method")
    }
    
}

/*func leaveHandler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintln(w, "<h1>%s!</h1>", r.URL.Path[1:])
    r.ParseForm()
    if r.Method == "GET" {
    	fmt.Fprintf(w, "Error Method")
    } else {
    	ip := strings.TrimSpace(r.PostFormValue("clientip"))
    	groupid := strings.TrimSpace(r.PostFormValue("groupid"))
    	comMainGroup(groupid, "remove_server")
    	//leaveGroup(ip, groupid)
    } 
}*/

func joinHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
    	hasher := md5.New()
        io.WriteString(hasher, strconv.FormatInt(time.Now().Unix(), 10))
        token := fmt.Sprintf("%x", hasher.Sum(nil))
    	t, _ := template.ParseFiles("/Users/rachel/Documents/goworkspace/musicsync/src/subserver/join.html")
    	t.Execute(w,token)
    } else if r.Method == "POST"{
    	r.ParseForm()
    	groupId := strings.TrimSpace(r.PostFormValue("groupid"))
    	//groupMember := joinGroup(ip, groupid)
    	comMainServer(groupId, "join_group")
    	/*fmt.Println(groupMember)
    	//send back server ips to client
    	for i := range groupMember {
    		w.Write([]byte(groupMember[i]))
    		w.Write([]byte("\n"))
    	} 	*/
    	
    }
}

func startHTTP() {
	fmt.Println("HTTP server start");
	http.HandleFunc("/create", createHandler)
    http.HandleFunc("/join", joinHandler)
    http.HandleFunc("/leave", leaveHandler)
    http.ListenAndServe(":8080", nil)
   
}

func comMainServer(data string, kind string) {
	//fmt.Println("start client");
	msg := &Message{"192.168.0.107:9090", "localhost:5000", kind, data}
    conn, err := net.Dial("tcp", msg.Dst)
    if err != nil {
        fmt.Println("Connection error: ", err)
    }
    encoder := gob.NewEncoder(conn)
    encoder.Encode(msg)
   
    connbuf := bufio.NewReader(conn)
	for{
    	str, err := connbuf.ReadString('\n')
    	if len(str)>0{
        	fmt.Println(str)
    	}
    	if err!= nil {
        	break
    	}
	}
    conn.Close()
    fmt.Println("done");
}

func main() {
    startHTTP()
}