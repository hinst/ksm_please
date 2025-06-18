package main

import (
	"fmt"
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
	var filePathMatcher = regexp.MustCompile(`Loaded: loaded \((\S+);`)
	var output, _ = exec.Command("systemctl", "status", serviceName).Output()
	var filePaths = filePathMatcher.FindSubmatch(output)
	if len(filePaths) > 1 {
		var filePath = string(filePaths[1])
		fmt.Println(filePath)
	}
}

func main() {
	log.Println("STARTING")
	var allServices = getAllServices()
	for _, serviceName := range allServices {
		getServiceStatus(serviceName)
	}
}
