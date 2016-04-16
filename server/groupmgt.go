package main

import (
	"fmt"
	"strconv"
	"encoding/csv"
	"os"
	"io"
)

type Group struct{
	id string
	name string
	members map[string]bool
}

func (g Group) addMember(ip string){
	g.members[ip] = true
}

func (g Group) removeMember(ip string){
	delete(g.members, ip)
}

var groups []Group
var group_num int

func joinGroup(clientIp string, groupId string) []string{
	var memberSlice []string
	for i:=0; i<len(groups); i++ {
		if groups[i].id == groupId {
			groups[i].addMember(clientIp)
			for k := range groups[i].members {
				memberSlice = append(memberSlice, k)
			}	
		}
	}
	return memberSlice
}

func createGroup(creatorIp string, groupName string){
	var newGroup Group
	newGroup.id = strconv.Itoa(group_num + 1)
	group_num++
	newGroup.name = groupName
	newGroup.members = make(map[string]bool)
	newGroup.addMember(creatorIp)
	groups = append(groups, newGroup)
	fmt.Println(groups)
}

func leaveGroup(clientIp string, groupId string){
	for i:=0; i<len(groups); i++ {
		if groups[i].id == groupId {
			groups[i].removeMember(clientIp)
			fmt.Println("delete: ", clientIp)
		}
	}
	/*****************debug**************/
	for i := range groups {
		fmt.Printf("  %v\n", groups[i])
	}
	/*****************debug**************/
}

func initGroups(){
	fmt.Println("hello, I am music sync server") 
	readCSV()
}

func readCSV(){
	file, err:= os.Open("/Users/rachel/Documents/goworkspace/musicsync/src/server/initGroups.csv")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	group_num = 0
	
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
		newGroup.members = make(map[string]bool)
		newGroup.addMember(record[2])
		groups = append(groups, newGroup)
		group_num++
    }
}


