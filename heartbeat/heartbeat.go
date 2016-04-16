package heartbeat

/*
 * Original thought: 
 */

import (
	"fmt"
	"os"
	"net"
	"sync"
)

const (
	HEARTBEAT_FREQUENCY = 1
	DEAD_DETECT = 3
)

type HeartBeat struct {
	host string
	track_server []string
	track_server_addr []*net.UDPAddr
	listenSock *net.UDPConn
	timeStamps map[string]time.time
	deadChannel chan string
	lock *sync.Mutex
}

func checkErr(err error){
	if err != nil {
		fmt.Println("<Error> ", err)
		os.Exit(0)
	}
}

func (this *HeartBeat) newInstance(host string, connect_servers []string){
	this.host = host
	// Set up listen socket
	addr, err := net.ResolveUDPAddr("udp", this.host)
	checkErr(err)
	this.listenSock, err := net.ListenUDP("udp", addr)
	checkErr(err)

	// Initiallize the arguments
	this.lock = new(sync.Mutex)
	this.deadChannel = make(chan string)
	this.updateAliveList(connect_servers)
	go this.recvAliveMsg()
	go this.sendAliveMsg()
}

func (this *HeartBeat) updateAliveList(connect_servers []string){
	this.lock.Lock()
	this.track_server_addr = make([]*net.UDPAddr, len(connect_servers))
	for idx, server := range connect_servers{
		addr, err := net.ResolveUDPAddr("udp", server)
		checkErr(err)
		this.track_server_addr[idx] = addr
	}
	this.track_server = connect_servers
	this.timeStamps = make(map[string]time.time)
	this.lock.Unlock()
}

func (this *HeartBeat) recvAliveMsg(){
	for{
		buffer := make([]byte, 64)
		numBytes, _, err := this.listenSock.ReadFromUDP(buffer)
		checkErr(err)
		recvServerName := string(buffer[:numBytes])
		// Update the timer
		this.lock.Lock()
		this.timeStamps[recvServerName] = time.Now()
		this.lock.Unlock()
	}
}

func (this *HeartBeat) sendAliveMsg(){
	var ticker time.Ticker = time.NewTicker(time.Second * HEARTBEAT_FREQUENCY)
	for _ := range ticker.C {
		this.lock.Lock()
		// Send message to every other servers
		for _, addr := range this.to {
			_, _ := this.listenSock.WriteToUDP([]byte(this.host), addr)
		}
		
		// To check whether the track servers are still alive
		for server, latestTime := range this.timeStamps {
			if time.Now().After(latestTime.Add(time.Second * DEAD_DETECT)) {
				fmt.Println("Found a dead server", server)
				delete(this.timeStamps, server)
				this.dead <- server
			}
		}
		this.lock.Unlock()
	}
}
