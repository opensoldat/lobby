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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

// types
type Configuration struct {
	Port int
}

type Servers struct {
	//Servers map[uint64]Server `json:"Servers"`
	Servers []Server `json:"Servers"`
}

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

// globals
// var servers = make(map[uint64]Server)
// var servers = Servers{Servers: make(map[uint64]Server)}
var servers = Servers{Servers: []Server{}}

var templateServers = Servers{Servers: []Server{}}
var configuration Configuration = Configuration{}

// functions
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// workaround for missing players endpoint
	w.Header().Set("Content-Type", "application/json")
	setupCorsResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	fmt.Fprintf(w, "{\"Players\":[]}")

	//fmt.Fprintf(w, "<h1 style='font-family:monospace'>slobby v0.0.0 alpha1</h1>")
}

func serversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupCorsResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	/*
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
	*/

	output, err := json.Marshal(servers)
	if err == nil {
		fmt.Fprintf(w, string(output))
	} else {
		fmt.Fprintf(w, "{\"error\": \"503: internal server error\"}")
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func hasFlag(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func handleSettings() {
	var configPath string
	var port int
	var writeServerTemplate bool

	flag.StringVar(&configPath, "config", ".", "config folder")
	flag.IntVar(&port, "port", 80, "lobby port")
	flag.BoolVar(&writeServerTemplate, "write-servers-example", false, "writes servers-example.json template file")
	flag.Parse()

	if writeServerTemplate {
		generateTemplateServersConfig()
		os.Exit(0)
	}

	readConfigFile(configPath, hasFlag("config"))

	if port <= math.MaxUint16 && port >= 0 && hasFlag("port") {
		configuration.Port = port
	}

	readServersFile(configPath)

}

func readConfigFile(path string, isCustom bool) {
	if isCustom {
		println("Using custom config path:", path+"/config.json\n")
	}
	content, err := ioutil.ReadFile(path + "/config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &configuration)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	if configuration.Port > math.MaxInt16 || configuration.Port < 0 {
		configuration.Port = 80
	}
}

func generateTemplateServersConfig() {
	//var templateServers =  Servers{Servers: []Server{}}
	templateServers.Servers = append(templateServers.Servers, Server{
		AC:             false,
		Advanced:       false,
		Dedicated:      true,
		Private:        false,
		Realistic:      false,
		Survival:       false,
		WM:             false,
		BonusFreq:      240,
		ConnectionType: 3,
		MaxPlayers:     24,
		NumBots:        0,
		NumPlayers:     0,
		Port:           23073,
		Respawn:        60,
		Country:        "DE",
		CurrentMap:     "ctf_Ash",
		Version:        "1.8.0",
		GameStyle:      "CTF",
		IP:             "127.0.0.1",
		Info:           "Server Info",
		Name:           "OpenSoldat Test Server",
		OS:             "linux",
		Players:        nil,
	})

	file, err := json.MarshalIndent(templateServers, "", " ")
	if err != nil {
		log.Fatal("Error during Marshal(): ", err)
	}
	err = ioutil.WriteFile("servers-example.json", file, 0644)
	if err != nil {
		log.Fatal("Error when writing file: ", err)
	}
}

func readServersFile(path string) {
	content, err := ioutil.ReadFile(path + "/servers.json")
	if err != nil {
		log.Print("Error when opening file: ", err)
		return

		//log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &servers)
	if err != nil {
		log.Print("Error during Unmarshal(): ", err)
		return
		//log.Fatal("Error during Unmarshal(): ", err)
	}
	//fmt.Println("Debug:", servers.Servers[0])
}

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

func registerHandlers() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/v0/servers", serversHandler)
	http.HandleFunc("/v0/register", registerHandler)
}

func main() {
	handleSettings()

	registerHandlers()

	//println("Port: %d\n", configuration.Port)
	fmt.Printf("Server started: http://localhost:%d/v0/servers\n", configuration.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(configuration.Port), nil)
	if err != nil {
		panic(err)
	}
}

// vim: ts=2:sts=2:sw=2:noet
