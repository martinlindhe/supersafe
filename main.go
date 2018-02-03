package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/feixiao/httpprof"
	"github.com/julienschmidt/httprouter"
)

func init() {
	// prefix file:line
	log.SetFlags(log.Lshortfile)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "supersafe!\n")
}

func run(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	comReader, _, err := r.FormFile("com")
	if err != nil {
		log.Println("run com post err:", err)
		return
	}
	com, err := ioutil.ReadAll(comReader)
	if err != nil {
		log.Println("run com read err:", err)
		return
	}

	if runtime.GOOS == "windows" {
		file, err := ioutil.TempFile(os.TempDir(), "prober")
		if err != nil {
			log.Println("error creating temp file", err)
			return
		}
		proberCom := file.Name() + ".com"
		defer os.Remove(proberCom)

		err = ioutil.WriteFile(proberCom, com, 0644)
		if err != nil {
			log.Println("error writing com to disk", err)
			return
		}

		t := time.Now()
		fmt.Printf("%s: exec %s (%d bytes)\n", t.Format("15:04:05"), proberCom, len(com))
		out, err := exec.Command("cmd", "/C", proberCom).Output()
		if err != nil {
			log.Println("error running command", err)
			return
		}
		fmt.Fprintf(w, string(out))
	} else {
		fmt.Printf("not windows: skipping execution")
	}
}

var port = 28111

func main() {
	runtime.SetBlockProfileRate(1)

	router := httprouter.New()
	router = httpprof.WrapRouter(router) // Register pprof handlers
	router.GET("/", index)
	router.POST("/run", run)

	ip, err := externalIP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("supersafe started, accepting connections on http://%s:%d/\n", ip, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
