// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/svc"
	"switcher"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, debug, start, or stop.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

func main() {
	// This is the name you will use for the NET START command
	const svcName = "switcher"
	// This is the name that will appear in the Services control panel
	const svcDesc = "Provider Switcher Service"

	// detect service mode
	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}

	service := switcher.NewService(svcName, svcDesc, isIntSess)

	// check service mode
	if !isIntSess {
		service.Run()
		return
	}

	// check params in the interactive mode
	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		service.Run()
		return
	case "install":
		err = service.Install()
	case "remove":
		err = service.Remove()
	case "start":
		err = service.Start()
	case "stop":
		err = service.Control(svc.Stop, svc.Stopped)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
	return
}
