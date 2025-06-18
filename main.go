package main

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func getAllServices() (serviceNames []string) {
	var matcher = regexp.MustCompile(`\s[\S-]+`)
	var allUnitsText = string(assertResultError(exec.Command("systemctl", "list-units", "--type=service", "--all").Output()))
	for index, unitText := range strings.Split(allUnitsText, "\n") {
		if len(strings.TrimSpace(unitText)) == 0 {
			break
		}
		if index == 0 {
			continue
		}
		var serviceName = matcher.FindString(unitText)
		serviceName = strings.TrimSpace(serviceName)
		serviceNames = append(serviceNames, serviceName)
	}
	return
}

func getServiceStatus(serviceName string) {
	var output, commandError = exec.Command("systemctl", "status", serviceName).Output()
	var text = string(output)
	if commandError == nil {
		log.Println(text)
	} else {
		log.Println("Cannot read " + serviceName)
		log.Println("error: " + text)
	}
}

func main() {
	log.Println("STARTING")
	var allServices = getAllServices()
	for _, serviceName := range allServices {
		getServiceStatus(serviceName)
	}
}
