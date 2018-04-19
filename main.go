// Author: zheng-ji.info

package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	pConfig    ProxyConfig
	pLog       *logrus.Logger
	configFile = flag.String("c", "etc/conf.yaml", "配置文件，默认etc/conf.yaml")
)

func onExitSignal() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGUSR1, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
L:
	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGUSR1:
			log.Fatal("Reopen log file")
		case syscall.SIGHUP:
			log.Fatal("Reload config file")
		case syscall.SIGTERM, syscall.SIGINT:
			log.Fatal("Catch SIGTERM singal, exit.")
			break L
		}
	}
}

func main() {

	flag.Parse()

	if parseConfigFile(*configFile) != nil {
		return
	}

	// need to reload some config
	initLogger()
	initBackendSvrs(pConfig.Backend)
	initAllowedIPs(pConfig.AllowedIPs)

	go onExitSignal()
	fmt.Println("Start Proxy bind ", pConfig.Bind)

	// init status service
	initStats()

	// init proxy service
	initProxy()
}
