package main 

import (
    "net/http"
    "fmt"
    "strings"
    "net"
    "html/template"
    "io"
    "os"
    "time"
    "crypto/md5"
    "strconv"
)

// Default Request Handler
func createHandler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "<h1>Hello %s!</h1>", r.URL.Path[1:])
    if r.Method == "GET" {
    	//fmt.Fprintf(w, "Error Method")
    	hasher := md5.New()
        io.WriteString(hasher, strconv.FormatInt(time.Now().Unix(), 10))
        token := fmt.Sprintf("%x", hasher.Sum(nil))
    	t, _ := template.ParseFiles("/Users/rachel/Documents/goworkspace/musicsync/src/server/create.html")
    	t.Execute(w,token)
    } else if r.Method == "POST" {
    	r.ParseForm()
    	//ip := strings.TrimSpace(r.PostFormValue("clientip"))
    	ip,_,_ := net.SplitHostPort(r.RemoteAddr)
    	//fmt.Println(ip)
    	groupName := strings.TrimSpace(r.PostFormValue("groupname"))
    	fmt.Println(ip)
    	fmt.Println(groupName)
    	createGroup(ip, groupName)
    } else {
    	fmt.Fprintf(w, "Error Method")
    }
    
}

/**
 * Update local list
 * Reply client requested group members
 */
func joinHandler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "<h1>%s!</h1>", r.URL.Path[1:])
    
    if r.Method == "GET" {
    	hasher := md5.New()
        io.WriteString(hasher, strconv.FormatInt(time.Now().Unix(), 10))
        token := fmt.Sprintf("%x", hasher.Sum(nil))
    	t, _ := template.ParseFiles("/Users/rachel/Documents/goworkspace/musicsync/src/server/join.html")
    	t.Execute(w,token)
    } else if r.Method == "POST"{
    	r.ParseForm()
    	ip := strings.TrimSpace(r.PostFormValue("clientip"))
    	groupid := strings.TrimSpace(r.PostFormValue("groupid"))
    	groupMember := joinGroup(ip, groupid)
    	fmt.Println(groupMember)
    	//send back server ips to client
    	for i := range groupMember {
    		w.Write([]byte(groupMember[i]))
    		w.Write([]byte("\n"))
    	} 	
    	
    }
}

func leaveHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "<h1>%s!</h1>", r.URL.Path[1:])
    r.ParseForm()
    if r.Method == "GET" {
    	fmt.Fprintf(w, "Error Method")
    } else {
    	ip := strings.TrimSpace(r.PostFormValue("clientip"))
    	groupid := strings.TrimSpace(r.PostFormValue("groupid"))
    	leaveGroup(ip, groupid)
    } 
}


func test(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("ok"))
	addrs, err := net.InterfaceAddrs()
	if err != nil {
   		panic(err)
	}   
	for i, addr := range addrs {
    	fmt.Fprintf(w, "%d %v\n", i, addr)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        hasher := md5.New()
        io.WriteString(hasher, strconv.FormatInt(time.Now().Unix(), 10))
        token := fmt.Sprintf("%x", hasher.Sum(nil))
        t, _ := template.ParseFiles("/Users/rachel/Documents/goworkspace/musicsync/src/server/upload.gtpl")
        t.Execute(w, token)
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
        	http.Redirect(w, r, "http://127.0.0.1:8080/upload", 301)
            fmt.Println(err)
            return
        }
        defer file.Close()
        fmt.Fprintf(w, "%v", handler.Header)
        fmt.Println(handler.Filename)
        f, err := os.OpenFile("/Users/rachel/Documents/goworkspace/musicsync/src/server/test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
    }
}

func main() {
	initGroups()
	
	http.HandleFunc("/upload", upload)
    http.HandleFunc("/create", createHandler)
    http.HandleFunc("/join", joinHandler)
    http.HandleFunc("/leave", leaveHandler)
    http.ListenAndServe(":8080", nil)
    
}

