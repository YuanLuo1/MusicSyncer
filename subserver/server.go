package main

import (
	"net/http"
    "fmt"
    "strings"
    "html/template"
    "io"
    "time"
    "crypto/md5"
    "strconv"
    "os"
)

type Server struct {
    ip string
    comm_port string
    http_port string
    heartbeat_port string
}

var (
    groups []Group
    // servers []string
    // heartbeatPort []string
    localGroups []GroupMusic
    hasGroups map[string]bool
    dir string
    // myIp string
    // myPort string
    // myAddr string
    servers []Server
    myServer Server
    heartBeatTracker = new(HeartBeat)
)

// var groups []Group 
// var servers []string 
// var localGroups []GroupMusic
// var hasGroups map[string]bool

// var dir string
// var myIp string
// var myPort string
// var myAddr string

func createHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(dir)
    if r.Method == "GET" {
    	t, _ := template.ParseFiles("UI/create.html")

    	t.Execute(w,nil)
    } else if r.Method == "POST" {
    	r.ParseForm()
    	//ip,_,_ := net.SplitHostPort(r.RemoteAddr)
    	groupName := strings.TrimSpace(r.PostFormValue("groupname"))
    	if !isGroupNameExist(groupName) {  		
    		multicastServers(groupName, "create_group")
    		var newGroup Group
			newGroup.name = groupName
			newGroup.serverList = make(map[string]bool)
			newGroup.addServer(myServer.ip + ":" + myServer.comm_port)
			groups = append(groups, newGroup)
			hasGroups[groupName] = true
			fmt.Println(groups)
			fmt.Println(hasGroups)
			http.Redirect(w, r, "http://127.0.0.1:8282/upload", 301)
    	} else {
    		w.Write([]byte("group name exist, please try another"))
    	}
    } else {
    	fmt.Fprintf(w, "Error Method")
    }
    
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
    	hasher := md5.New()
        io.WriteString(hasher, strconv.FormatInt(time.Now().Unix(), 10))
        token := fmt.Sprintf("%x", hasher.Sum(nil))
    	t, _ := template.ParseFiles("UI/join.html")
    	t.Execute(w,token)
    } else if r.Method == "POST"{
    	r.ParseForm()
    	groupName := strings.TrimSpace(r.PostFormValue("groupname"))
    	if isGroupHere(groupName) {
    		w.Write([]byte("you are in the group " + groupName))
    	} else {
    		multicastServers(groupName, "join_group")
    		hasGroups[groupName] = true
    		w.Write([]byte("join successful"))    		
    	}
    	//groupMember := joinGroup(ip, groupid)
    	//comMainServer(groupId, "join_group")
    	//fmt.Println(groupMember)
    	//send back server ips to client
    	//for i := range groupMember {
    		//w.Write([]byte(groupMember[i]))
    		//w.Write([]byte("\n"))
    	//} 	
    	
    }
}

func addfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        hasher := md5.New()
        io.WriteString(hasher, strconv.FormatInt(time.Now().Unix(), 10))
        token := fmt.Sprintf("%x", hasher.Sum(nil))
        t, _ := template.ParseFiles("UI/upload.html")
        t.Execute(w, token)
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
        	http.Redirect(w, r, "localhost:8282/upload.html", 301)
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


func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
    	t, _ := template.ParseFiles("UI/index.html")
    	t.Execute(w, nil)
    } else {
    	fmt.Fprintf(w, "Error Method")
    }
}

func startHTTP() {
	fmt.Println("HTTP server start")
  	
  	http.Handle("/css/", http.FileServer(http.Dir("UI")))
    http.Handle("/js/", http.FileServer(http.Dir("UI")))
    http.Handle("/images/", http.FileServer(http.Dir("UI")))
    http.Handle("/fonts/", http.FileServer(http.Dir("UI")))
    
	http.HandleFunc("/index.html", homeHandler)
	http.HandleFunc("/create.html", createHandler)
    http.HandleFunc("/join.html", joinHandler)
    http.HandleFunc("/upload.html", addfileHandler)

    //http.HandleFunc("/leave", leaveHandler)
    // var httpPort string
    // switch myPort {
    // 	case "9292": httpPort = ":8282"
    // 	case "9293": httpPort = ":8283"
    // 	case "9294": httpPort = ":8284"
    // 	case "9295": httpPort = ":8285"
    // }
    http.ListenAndServe(myServer.http_port, nil)
   
}


/* Update the servers if recv a dead server */
func getDeadServer(){
    fmt.Println("Get dead servers from heartbeat manager's deadchannel")
    deadServerChannel := heartBeatTracker.GetDeadChannel()
    // TODO: Do something for the dead servers
    for {

    }
}

func InitialHeartBeat(){
    fmt.Println("Initialize heartbeat")
    // argument : (myIP, other servers)
    heartBeatTracker.newInstance(myIP+":"+myPort, )

}

func main() {
	readServerConfig() 
	
	//select server's configuration
    fmt.Print("Enter a number(0-3) set up this server: ")
    var i int
    fmt.Scan(&i)
    myServer = servers[i]
    // myAddr = servers[i]
    // localInfo := strings.Split(servers[i],":")
    // myIp = localInfo[0]
    // myPort = localInfo[1]
    // fmt.Println(myIp, myPort)

	readGroupConfig()
	readMusicConfig()
	
	heartbeat()
	go listeningMsg(myServer.ip, myServer.comm_port)
    startHTTP()
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