package tf_helper

import (
	"fmt"
	"log"
	"strconv"
)

func (c *Config) Graph() {

	cmd_name := "terraform"

	exec_args := []string{"graph"}

	log.Println("[INFO] Executing Terraform graph.")

	if len(c.TypeTF) > 0 {
		for _, t := range c.TypeTF {
			if t == "cycles" {
				exec_args = append(exec_args, "-draw-cycles")
			} else if _, err := strconv.Atoi(t); err == nil {
				exec_args = append(exec_args, fmt.Sprintf("-module-depth=%s", t))
			} else if t == "plan-destroy" || t == "apply" || t == "validate" || t == "refresh" {
				exec_args = append(exec_args, fmt.Sprintf("-type=%s", t))
			}
		}
	}

	if !ExecCmd(cmd_name, exec_args) {
		log.Fatal("[ERROR] Failed to execute Terraform validate. Aborting.")
	} else {
		log.Println("[INFO] Graph creation complete. Use a program like GraphViz to generate an svg file out of the raw data, e.g. through '$ dot -Tsvg output > graph.svg'")
	}
}
