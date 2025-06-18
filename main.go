package main

import (
	"errors"
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

func checkUnitFile(filePath string) (result *bool) {
	var file, fileError = os.OpenFile(filePath, 0, os.ModePerm)
	defer file.Close()
	if fileError == nil {
		var unitInfo, unitError = unit.DeserializeOptions(file)
		if unitError == nil {
			var memoryEnabled = false
			result = &memoryEnabled
			for _, unitValue := range unitInfo {
				if unitValue.Section == "Service" && unitValue.Name == "MemoryKSM" {
					memoryEnabled = unitValue.Value == "true"
				}
			}
		}
	}
	return
}

func insertMergeMemory(lines []string) (outputLines []string) {
	for _, line := range lines {
		outputLines = append(outputLines, line)
		if strings.TrimSpace(line) == "[Service]" {
			outputLines = append(outputLines, "MemoryKSM=true")
		}
	}
	return
}

func enableMergeMemory(filePath string) error {
	var bytes, fileError = os.ReadFile(filePath)
	var permissions = assertResultError(os.Stat(filePath))
	var backupBytes = bytes
	if fileError != nil {
		return fileError
	}
	var text = string(bytes)
	var lines = strings.Split(text, "\n")
	lines = insertMergeMemory(lines)
	text = strings.Join(lines, "\n")
	bytes = []byte(text)
	var writeError = os.WriteFile(filePath, bytes, permissions.Mode())
	if writeError != nil {
		return writeError
	}

	var status = checkUnitFile(filePath)
	if status == nil || !*status {
		var writeError = os.WriteFile(filePath, backupBytes, permissions.Mode())
		if writeError != nil {
			return writeError
		}
		return errors.New("sanity check failed, restoring backup")
	}
	return nil
}

func main() {
	log.Println("STARTING")
	var allServices = getAllServices()
	for _, serviceName := range allServices {
		var filePath = getServiceFilePath(serviceName)
		if len(filePath) == 0 {
			continue
		}
		var memoryEnabled = checkUnitFile(filePath)
		var status = "x"
		if memoryEnabled != nil {
			status = " "
			if *memoryEnabled {
				status = "Y"
			} else {
				enableMergeMemory(filePath)
			}
		}
		fmt.Printf("[%v] %v\n", status, filePath)
		if memoryEnabled != nil && !*memoryEnabled {
			var error = enableMergeMemory(filePath)
			if error != nil {
				log.Println(error)
			}
		}
		break
	}
}
