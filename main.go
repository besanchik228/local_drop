package main

import (
	"flag"
	"fmt"
	qr "github.com/skip2/go-qrcode"
	"net"
	"net/http"
	"path/filepath"
)

var path *string

func server_get(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("/storage/projects/local_drop", *path)
	http.ServeFile(w, r, filePath)
}

func main() {
	path = flag.String("p", ".", "Путь к объекту")
	flag.Parse()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Ошибка")
	}

	var addr string

	for _, i := range addrs {
		if ip, ok := i.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			v4 := ip.IP.To4()
			if v4 != nil && v4[0] == 192 && v4[1] == 168 {
				addr = ip.IP.String()
				break
			}
		}
	}

	http.HandleFunc("/", server_get)
	In, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println("Нету или ошибка порта")
	}
	defer In.Close()
	port := In.Addr().(*net.TCPAddr).Port
	code, err := qr.New(fmt.Sprintf("http://%s:%d", addr, port), qr.Medium)
	fmt.Printf("http://%s:%d\n", addr, port)
	if err != nil {
		fmt.Println("Ошибка при созданиии qr")
	}
	fmt.Println(code.ToSmallString(false))

	http.Serve(In, nil)
}
