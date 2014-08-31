package vagrant

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateVagrantDockerBox() {
	vagrantCommand := exec.Command("/usr/bin/vagrant", "up", "docker")
	fmt.Println("===================================")
	fmt.Println("Launching Vagrant Docker Ubuntu Box")
	fmt.Println("===================================")
	vagrantCommand.Stderr = os.Stderr
	vagrantCommand.Stdout = os.Stdout
	err := vagrantCommand.Run()
	if err != nil {
		fmt.Printf("error occured %s\n", err)
	}
}
