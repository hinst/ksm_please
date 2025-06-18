package main

import (
	"log"
	"os/exec"
	"strings"
)

func getAllServices() []string {
	var allUnitsText = string(assertResultError(exec.Command("systemctl", "list-units", "--type=service", "--all").Output()))
	for index, unitText := range strings.Split(allUnitsText, "\n") {
		if 0 == index {
			continue
		}
		log.Println(unitText)
	}
	return nil
}

func main() {
	log.Println("STARTING")
	getAllServices()
}
