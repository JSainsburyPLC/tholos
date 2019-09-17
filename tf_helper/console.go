package tf_helper

import (
	"log"
	"os/exec"
	"strings"
)

func (c *Config) Console() {

	//finding terraform as may be different per client

	cmd := exec.Command("which", "terraform")

	cmd_path, err := cmd.Output()

	if err != nil {
		log.Printf("Error finding terraform executable: %s, err")
	}

	exec_args := []string{"terraform", "console"}

	if !ExecOsCmd(strings.TrimRight(string(cmd_path), "\n"), exec_args) {
		log.Fatal("[ERROR] Failed to run Terraform Console. Aborting.")
	}

}
