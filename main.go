// soldat lobby - KISS reimplementation
//
// TODO: add mutex
// TODO: add register server support
// TODO: add more filters
// TODO: add all fields
// TODO: add old api

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Server struct {
	AC             bool     `json:"AC"`
	Advanced       bool     `json:"Advanced"`
	Dedicated      bool     `json:"Dedicated"`
	Private        bool     `json:"Private"`
	Realistic      bool     `json:"Realistic"`
	Survival       bool     `json:"Survival"`
	WM             bool     `json:"WM"`
	BonusFreq      uint8    `json:"BonusFreq"`
	ConnectionType uint8    `json:"ConnectionType"`
	MaxPlayers     uint8    `json:"MaxPlayers"`
	NumBots        uint8    `json:"NumBots"`
	NumPlayers     uint8    `json:"NumPlayers"`
	Port           uint16   `json:"Port"`
	Respawn        uint16   `json:"Respawn"`
	Country        string   `json:"Country"`
	CurrentMap     string   `json:"CurrentMap"`
	Version        string   `json:"Version"`
	GameStyle      string   `json:"GameStyle"`
	IP             string   `json:"IP"`
	Info           string   `json:"Info"`
	Name           string   `json:"Name"`
	OS             string   `json:"OS"`
	Players        []string `json:"-"`
}

var servers = make(map[uint64]Server)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1 style='font-family:monospace'>slobby v0.0.0 alpha1</h1>")
}

func serversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupCorsResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	fmt.Fprintf(w, "{\"Servers\":[")

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		return
	}
	queryParams, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}

	param_version := queryParams["version"]
	param_os := queryParams["os"]

	// TODO: lock here
	// TODO: defer unlock
	i := 0
	for _, server := range servers {
		if (param_version == nil || param_version[0] == server.Version) &&
			(param_os == nil || param_os[0] == server.OS) {
			output, err := json.Marshal(server)
			if err == nil {
				if i != 0 {
					fmt.Fprintf(w, ","+string(output))
				} else {
					fmt.Fprintf(w, string(output))
					i++
				}
			}
		}
	}

	fmt.Fprintf(w, "]}")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func main() {
	servers[0] = Server{
		AC:             false,
		Advanced:       false,
		Dedicated:      true,
		Private:        true,
		Realistic:      false,
		Survival:       false,
		WM:             false,
		BonusFreq:      240,
		ConnectionType: 3,
		MaxPlayers:     16,
		NumBots:        0,
		NumPlayers:     0,
		Port:           80,
		Respawn:        60,
		Country:        "DE",
		Version:        "1.8.0",
		GameStyle:      "CTF",
		IP:             "127.0.0.1",
		Info:           "info",
		Name:           "OpenSoldatServer 1!",
		OS:             "windows",
		Players:        nil,
	}

	servers[1] = Server{
		AC:             false,
		Advanced:       false,
		Dedicated:      true,
		Private:        true,
		Realistic:      false,
		Survival:       false,
		WM:             false,
		BonusFreq:      240,
		ConnectionType: 3,
		MaxPlayers:     16,
		NumBots:        0,
		NumPlayers:     0,
		Port:           81,
		Respawn:        60,
		Country:        "DE",
		Version:        "1.7.1",
		GameStyle:      "DM",
		IP:             "227.0.0.1",
		Info:           "info",
		Name:           "OpenSoldatServer 2!",
		OS:             "mac",
		Players:        nil,
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/v0/servers", serversHandler)
	http.HandleFunc("/v0/register", registerHandler)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

// vim: ts=2:sts=2:sw=2:noet
