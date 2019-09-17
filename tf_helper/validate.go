package tf_helper

import (
	"log"
)

func (c *Config) Validate() {

	cmd_name := "terraform"

	exec_args := []string{"validate"}

	log.Println("[INFO] Executing Terraform validate.")

	if len(c.TypeTF) > 0 {
		for _, t := range c.TypeTF {
			if t == "json" {
				exec_args = append(exec_args, "-json")
			}
		}
	}

	if !ExecCmd(cmd_name, exec_args) {
		log.Fatal("[ERROR] Failed to execute Terraform validate. Aborting.")
	}
}
