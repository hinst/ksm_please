package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/coreos/go-systemd/v22/unit"
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

func getServiceFilePath(serviceName string) (filePath string) {
	var filePathMatcher = regexp.MustCompile(`Loaded: loaded \((\S+);`)
	var output, _ = exec.Command("systemctl", "status", serviceName).Output()
	var filePaths = filePathMatcher.FindSubmatch(output)
	if len(filePaths) > 1 {
		filePath = string(filePaths[1])
	}
	return
}

func updateUnit(filePath string) {
	var file, fileError2 = os.OpenFile(filePath, 0, os.ModePerm)
	if fileError2 == nil {
		fmt.Println(filePath)
		var unitInfo, unitError = unit.Deserialize(file)
		if unitError == nil {
			for _, unitValue := range unitInfo {
				if unitValue.Section == "Service" && unitValue.Name == "MemoryKSM" {
					fmt.Println("MemoryKSM")
				}
			}
		}
	}
}

func main() {
	log.Println("STARTING")
	var allServices = getAllServices()
	for _, serviceName := range allServices {
		var filePath = getServiceFilePath(serviceName)
		updateUnit(filePath)
	}
}
