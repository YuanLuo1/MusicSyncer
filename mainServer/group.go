package main

type Group struct{
	id string
	name string
	serverList map[string]bool	
}

func (g Group) addServer(ip string) {
	g.serverList[ip] = true
}

func (g Group) delServer(ip string) {
	delete(g.serverList, ip)
}

func (g Group) setId(id string) {
	g.id = id
}

func (g Group) setName(name string) {
	g.name = name
}