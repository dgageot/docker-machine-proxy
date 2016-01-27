package main

import (
	"net/http"

	"flag"

	"path/filepath"

	mcn "github.com/dgageot/docker-machine-proxy/machine"
	"github.com/dgageot/docker-machine-proxy/proxy"
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	machineName := flag.String("machine", "default", "Docker machine name")
	addr := flag.String("addr", "localhost:2375", "Address of the proxy")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		return nil
	}

	api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())
	defer api.Close()

	machine, err := api.Load(*machineName)
	if err != nil {
		return err
	}

	url, err := machine.URL()
	if err != nil {
		return err
	}

	return http.ListenAndServe(*addr, &proxy.DockerMachineProxy{
		Machine: &mcn.DockerMachine{
			Url:      url,
			CertPath: filepath.Join(mcndirs.GetBaseDir(), "machines", *machineName),
		},
	})
}
