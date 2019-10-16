package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/kardianos/service"
)

var logger service.Logger

func main() {
	port := flag.Int("p", 20517, "Port Number to listen on.")
	svcFlag := flag.String("service", "", "Service action.  Valid actions are: 'start', 'stop', 'restart', 'install' and 'uninstall'")
	flag.Parse()

	// Create a new server
	s := &Server{
		PortNo: *port,
	}

	// Create the service
	svcConfig := &service.Config{
		Name:        "PiMonitor",
		DisplayName: "PiMonitor",
		Description: "Monitors Raspberry pi status.",
	}
	v, err := service.New(s, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Set up the logger
	errs := make(chan error, 5)
	logger, err = v.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	// Start the service
	if *svcFlag != "" {
		// Service control request
		if err := service.Control(v, *svcFlag); err != nil {
			e := err.Error()
			if strings.Contains(e, "Unknown action") {
				fmt.Println(*svcFlag, "is an invalid action.")
				fmt.Println("Valid actions are", service.ControlAction)
			} else {
				fmt.Println(err.Error())
			}
		}
	} else {
		// Start the service in debug if we are running in a terminal
		s.VerboseLogging = service.Interactive()
		if err := v.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
