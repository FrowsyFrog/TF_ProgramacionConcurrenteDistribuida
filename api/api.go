package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	hostAddress string
	apiPort     = ":9000"
	nodePorts   = []string{":8000", ":8001", ":8002"}
	currentNode = 0
	mutex       sync.Mutex
)

func predecirHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get Client Array
	log.Println("Llamada a endpoint /predecir")
	aValue := req.FormValue("arrayValue")
	arr := readArray(aValue)
	fmt.Println("Arreglo recibido:", arr)

	// Get Predictions from Server
	nodeAddr := hostAddress + getNodePort()
	fmt.Println("---------------------------------------------------")
	fmt.Println(nodeAddr)
	predictionsArr := makePredictionRequest(nodeAddr, arr)
	fmt.Println("Predicciones: ", predictionsArr)

	// Send message to Client
	res.Header().Set("Content-Type", "application/json")
	bytesArrayMsg, _ := json.Marshal(predictionsArr)
	io.WriteString(res, string(bytesArrayMsg))
}

func makePredictionRequest(nodeAddr string, arr []float64) []float64 {
	conn, _ := net.Dial("tcp", nodeAddr)
	defer conn.Close()

	bytesMsg, _ := json.Marshal(arr)
	fmt.Fprintln(conn, string(bytesMsg))

	pred := getPredictions(conn)

	return pred
}

func getPredictions(conn net.Conn) []float64 {
	msg, _ := bufio.NewReader(conn).ReadString('\n')
	msg = strings.TrimSpace(msg)
	var y []float64
	json.Unmarshal([]byte(msg), &y)
	return y
}

func getNodePort() string {
	mutex.Lock()
	defer mutex.Unlock()
	auxCurrentNode := currentNode
	currentNode = (currentNode + 1) % len(nodePorts)
	return nodePorts[auxCurrentNode]
}

func readArray(strArr string) []float64 {
	strArr = strings.TrimSpace(strArr)

	arr := []float64{}
	for _, strNum := range strings.Split(strArr, " ") {
		num, _ := strconv.ParseFloat(strNum, 64)
		arr = append(arr, num)
	}
	return arr
}

func discoverIP() string {
	interfaces, _ := net.Interfaces()
	fmt.Println("Tamaño de Interfaces:", len(interfaces))
	for _, iface := range interfaces {
		fmt.Println("Nombre de interfaz:", iface.Name)
		if strings.HasPrefix(iface.Name, "Ethernet") {
			fmt.Println("¡Tiene prefijo 'Ethernet'!")
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				switch t := addr.(type) {
				case *net.IPNet:
					if t.IP.To4() != nil {
						return t.IP.To4().String()
					}
				}
			}
		}
	}
	return "127.0.0.1"
}

func requestHandler() {
	http.HandleFunc("/predecir", predecirHandler)
	log.Fatal(http.ListenAndServe(apiPort, nil))
}

func main() {
	hostAddress = discoverIP()
	requestHandler()
}
