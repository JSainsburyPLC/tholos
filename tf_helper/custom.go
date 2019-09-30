package tf_helper

import (
	"log"
	"os/exec"
	"strings"
)

func (c *Config) Custom(customArgs []string) {

	//finding terraform as may be different per client

	cmd := exec.Command("which", "terraform")

	cmd_path, err := cmd.Output()

	if err != nil {
		log.Printf("Error finding terraform executable: %s", err)
	}

	c.Setup_remote_state()

	exec_args := []string{"terraform"}

	if len(customArgs) > 0 {
		exec_args = append(exec_args, customArgs...)
	}

	log.Println("[INFO] Running Terraform custom command.")

	if !ExecOsCmd(strings.TrimRight(string(cmd_path), "\n"), exec_args) {
		log.Fatal("[ERROR] Failed to run Terraform custom command. Aborting.")
	}

}
