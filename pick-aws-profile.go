package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCommand(command string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func main() {
	// Runs the aws list-profiles command to collect all existing profiles
	listCmd := "aws"
	listCmdArgs := []string{"configure", "list-profiles"}
	listStdout, _, err := runCommand(listCmd, listCmdArgs...)

	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// Formats the list of profiles
	profiles := strings.Split(strings.TrimSpace(listStdout), "\n")
	nProfiles := len(profiles)
	var pretty_profiles []string

	for idx, val := range profiles {
		pretty_profile := fmt.Sprintf("[%d] %s", idx, val)
		pretty_profiles = append(pretty_profiles, pretty_profile)
	}

	fmt.Println(fmt.Sprintf("Profiles found: [%d]", nProfiles))
	fmt.Println(strings.Join(pretty_profiles, "\n"))

	// Asks the user to pick a profile from the printed ones by index
	fmt.Print("Pick a profile by index: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	idx, err := strconv.Atoi(text)

	if err != nil {
		log.Fatalf("error during input parsing: %v\n", err)
	}

	if idx < 0 || idx >= nProfiles {
		log.Fatalf("Invalid index: %d\n", idx)
	}

	fmt.Printf("Copying profile: %s\n", profiles[idx])

	// Copy the choice to the clipboard (macos pbcopy, macos zsh, not really portable)
	shell := "zsh"
	cpCommand := fmt.Sprintf("echo \"export AWS_PROFILE=%s\" | pbcopy", profiles[idx])
	cpStdout, _, err := runCommand(shell, "-c", cpCommand)

	if err != nil {
		log.Fatalf("error during copy: %d\n", err)
	}

	fmt.Println("You can now paste and run the export command")
	fmt.Println(cpStdout)

	os.Exit(0)
}
