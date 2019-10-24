package tf_helper

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func ExecCmd(cmdName string, args []string) bool {

	success := true

	log.Printf("[INFO] Executing command: %s %s", cmdName, strings.Join(args, " "))

	cmd := exec.Command(cmdName, args...)

	cmdReader, err := cmd.StdoutPipe()

	var wg sync.WaitGroup

	if err != nil {
		success = false
		log.Printf("Error creating StdoutPipe for Command: %s, Error: %s\n", cmdName, err.Error())
	}

	cmdErrorReader, err := cmd.StderrPipe()
	if err != nil {
		success = false
		log.Printf("Error creating StderrPipe for Command: %s, Error: %s\n", cmdName, err.Error())
	}

	wg.Add(1)

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		wg.Done()
	}()

	wg.Add(1)

	errorScanner := bufio.NewScanner(cmdErrorReader)
	go func() {
		for errorScanner.Scan() {
			fmt.Println(errorScanner.Text())
		}

		wg.Done()
	}()

	err = cmd.Start()
	if err != nil {
		success = false
		log.Printf("Error starting Command: %s, Error: %s\n", cmdName, err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		success = false
		log.Printf("Error waiting for Command: %s, Error: %s\n", cmdName, err.Error())
	}

	wg.Wait()

	return success

}

func ExecOsCmd(cmdPath string, args []string) bool {

	success := true

	log.Printf("[INFO] Executing command: %s", strings.Join(args, " "))

	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	proc, err := os.StartProcess(cmdPath, args, &pa)
	if err != nil {
		success = false
		log.Printf("Error starting process, Path: %s, Args: %s, Error: %s", cmdPath, args, err)
	}

	state, err := proc.Wait()
	if err != nil {
		success = false
		log.Printf("Error creating process wait: %s", err)
	}

	log.Printf("<< Custom command exited with code :%s\n", state.String())

	return success
}
