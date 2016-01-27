package main

import (
	"net/http"

	"flag"

	"github.com/dgageot/docker-machine-proxy/machine"
	"github.com/dgageot/docker-machine-proxy/proxy"
)

func main() {
	url := flag.String("url", "tcp:192.168.99.100:2376", "Url of the Docker Machine")
	certPath := flag.String("certPath", "/Users/dgageot/.docker/machine/machines/default", "Location of the certificates")
	addr := flag.String("addr", "127.0.0.1:2376", "Address of the proxy")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	http.ListenAndServe(*addr, &proxy.DockerMachineProxy{
		Machine: &machine.DockerMachine{
			Url:      *url,
			CertPath: *certPath,
		},
	})
}
