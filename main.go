// soldat lobby - KISS reimplementation
//
// TODO: add mutex
// TODO: add register server support
// TODO: add more filters
// TODO: add all fields
// TODO: add old api

package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

type Server struct {
	//AC             bool     `json:"AC"`
	//Advanced       bool     `json:"Advanced"`
	//Dedicated      bool     `json:"Dedicated"`
	//Private        bool     `json:"Private"`
	//Realistic      bool     `json:"Realistic"`
	//Survival       bool     `json:"Survival"`
	//WM             bool     `json:"WM"`
	//BonusFreq      uint8    `json:"BonusFreq"`
	//ConnectionType uint8    `json:"ConnectionType"`
	//MaxPlayers     uint8    `json:"MaxPlayers"`
	//NumBots        uint8    `json:"NumBots"`
	//NumPlayers     uint8    `json:"NumPlayers"`
	Port uint16 `json:"Port"`
	//Respawn        uint16   `json:"Respawn"`
	//Country        string   `json:"Country"`
	//CurrentMap     string   `json:"CurrentMap"`
	Version string `json:"Version"`
	//GameStyle      string   `json:"GameStyle"`
	IP   string `json:"IP"`
	Info string `json:"Info"`
	//Name           string   `json:"Name"`
	OS string `json:"OS"`
	//Players        []string `json:"-"`
}

var servers = make(map[uint64]Server)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1 style='font-family:monospace'>slobby v0.0.0 alpha1</h1>")
}

func serversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setupCorsResponse(&w, r)

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
	var ip uint32 = 0
	//var port uint16 = 0
	var port uint64 = 0
	var net_ip = net.ParseIP("0.0.0.1")
	if net_ip != nil {
		ipp := uint64(ip2int(net_ip))
		fmt.Println("ipp", ipp)
		ip = ip2int(net_ip)
		port = uint64(80 << 32)
	}

	id := generate_id("27.0.0.1", 80)
	fmt.Println("data", ip, port, id)

	servers[id] = Server{
		IP:      "127.0.0.1",
		Port:    80,
		Info:    "hello",
		OS:      "windows",
		Version: "1.7.1",
	}

	// example adding new entries
	servers[generate_id("127.0.0.1", 81)] = Server{
		IP:      "227.0.0.1",
		Port:    81,
		Info:    "hi",
		OS:      "mac",
		Version: "1.7.1",
	}

	servers[generate_id("227.0.0.1", 82)] = Server{
		IP:      "327.0.0.1",
		Port:    82,
		Info:    "hi",
		OS:      "linux",
		Version: "1.7.0",
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

func generate_id(ip_str string, port uint16) uint64 {
	const ERROR_RESULT = 0

	var result uint64 = ERROR_RESULT
	var net_ip = net.ParseIP(ip_str)
	if net_ip == nil {
		result = uint64(ip2int(net_ip))
		result += uint64(port << 32)
	}

	return result
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

/* vim: se ts=2:sts=2:sw=2:noet */
