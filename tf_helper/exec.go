package tf_helper

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func ExecCmd(cmdName string, args []string) bool {

	success := true

	log.Printf("[INFO] Executing command: %s %s", cmdName, strings.Join(args, " "))

	cmd := exec.Command(cmdName, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		success = false
		log.Printf("Error starting Command: %s, Error: %s\n", cmdName, err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		success = false
		log.Printf("Error waiting for Command: %s, Error: %s\n", cmdName, err.Error())
	}

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
