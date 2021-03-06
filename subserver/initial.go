package main

import (
	"fmt"
    "io"
    "os"
    "encoding/csv"
)

func readMusicConfig(){
	file, err:= os.Open("./initMusic.csv")
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
		
		if hasGroups[record[0]] {
			var newGroup GroupMusic
			newGroup.name = record[0]
			newGroup.musicList = []string{record[1]}
			localGroups = append(localGroups, newGroup)
		}
		
    }
	fmt.Println("Local group: ", localGroups)	
}

func readGroupConfig(){
	hasGroups = make(map[string]bool)
	
	file, err:= os.Open("./initGroups.csv")
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
		newGroup.name = record[0]
		newGroup.serverList = make(map[string]bool)
		newGroup.addServer(record[1])
		groups = append(groups, newGroup)
		
		if(record[1] == myAddr){
			hasGroups[record[0]] = true
		}
    }
	fmt.Println("Groups: ", groups)	
	fmt.Println("Has Groups: ", hasGroups)
}

func readServerConfig(){
	file, err:= os.Open("./initServers.csv")
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
		// Server struct {ip, comm_port, http_port, heart_beat_port}
		servers = append(servers, Server{record[1], record[2], record[3], record[4]}) 
		// servers = append(servers, record[1])
    }
	fmt.Println("Servers: ", servers)	
}
