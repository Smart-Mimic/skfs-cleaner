package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var routeDevices []RouteDevice

func fetchRouteDevices(routeID string) error {
	cmd := exec.Command("helium-config-service-cli", "route", "skfs", "list", "--route-id", routeID)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	err = json.Unmarshal(out.Bytes(), &routeDevices)
	if err != nil {
		return fmt.Errorf("failed to unmarshal route devices: %v", err)
	}

	return nil
}

func runRemoveActions(routeID string, config Config) error {
	today := time.Now().Format("20060102")
	removedDir := filepath.Join("updates", today, "removed")
	dryRun := config.DryRun

	files, err := filepath.Glob(filepath.Join(removedDir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list files in %s: %v", removedDir, err)
	}

	if len(files) == 0 {
		fmt.Println("No files to process in the removed directory.")
		return nil
	}

	// Regular expression to extract the "x" value from filenames like "part_x_y.json"
	re := regexp.MustCompile(`part_(\d+)_\d+\.json`)

	// Sort the files by the extracted "x" value
	sort.Slice(files, func(i, j int) bool {
		x1 := extractXValue(files[i], re)
		x2 := extractXValue(files[j], re)
		return x1 < x2
	})

	for _, file := range files {
		err := runUpdateCommand(routeID, file, dryRun)
		if err != nil {
			log.Printf("Error running update for file %s: %v", file, err)
		} else {
			log.Printf("Successfully processed %s\n", file)
		}
	}

	return nil
}

func runAddActions(routeID string, config Config) error {
	today := time.Now().Format("20060102")
	addedDir := filepath.Join("updates", today, "added")
	dryRun := config.DryRun

	files, err := filepath.Glob(filepath.Join(addedDir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list files in %s: %v", addedDir, err)
	}

	if len(files) == 0 {
		fmt.Println("No files to process in the added directory.")
		return nil
	}

	re := regexp.MustCompile(`part_(\d+)_\d+\.json`)

	sort.Slice(files, func(i, j int) bool {
		x1 := extractXValue(files[i], re)
		x2 := extractXValue(files[j], re)
		return x1 < x2
	})

	for _, file := range files {
		err := runUpdateCommand(routeID, file, dryRun)
		if err != nil {
			log.Printf("Error running update for file %s: %v", file, err)
		} else {
			log.Printf("Successfully processed %s\n", file)
		}
	}

	return nil
}

func extractXValue(file string, re *regexp.Regexp) int {
	match := re.FindStringSubmatch(filepath.Base(file))
	if len(match) > 1 {
		x, err := strconv.Atoi(match[1])
		if err == nil {
			return x
		}
	}
	return 0
}

func runUpdateCommand(routeID, file string, dryRun bool) error {
	cmdArgs := []string{
		"route", "skfs", "update",
		"--route-id", routeID,
		"--update-file", file,
	}

	if !dryRun {
		cmdArgs = append(cmdArgs, "--commit")
	}

	cmd := exec.Command("helium-config-service-cli", cmdArgs...)
	output, err := cmd.CombinedOutput()
	log.Printf("Processing file: %s\n", file)
	log.Printf("output: %v", string(output))
	time.Sleep(50 * time.Millisecond)
	if err != nil {
		return fmt.Errorf("command failed with error: %v, output: %s", err, string(output))
	}

	log.Printf("Command output: %s\n", string(output))
	return nil
}
