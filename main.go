package main

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func getAllServices() []string {
	var matcher = regexp.MustCompile(`\s[\S-]+`)
	var allUnitsText = string(assertResultError(exec.Command("systemctl", "list-units", "--type=service", "--all").Output()))
	for index, unitText := range strings.Split(allUnitsText, "\n") {
		if len(strings.TrimSpace(unitText)) == 0 {
			break
		}
		if index == 0 {
			continue
		}
		var unitName = matcher.FindString(unitText)
		unitName = strings.TrimSpace(unitName)
		log.Println(unitName)
	}
	return nil
}

func main() {
	log.Println("STARTING")
	getAllServices()
}
