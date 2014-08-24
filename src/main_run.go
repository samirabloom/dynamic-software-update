package main

import (
	proxy "proxy"
//	proxy "content_counting_proxy"
//	proxy "zmq_proxy"
)

func main() {
	proxy.CLI()
}

//import (
//	"fmt"
//	"os/exec"
//	"time"
//)
//
//func exe_cmd(name string, arg ...string) {
//	out, err := exec.Command(name, arg...).CombinedOutput()
//	time.Sleep(3000 * time.Millisecond)
//	if err != nil {
//		fmt.Printf("error occured %s\n", err)
//	}
//	fmt.Printf("%s", out)
//}
//
//func main() {
//	exe_cmd("docker", "pull", "jamesdbloom/couchbase")
//	exe_cmd("docker", "stop", "couch_one")
//	exe_cmd("docker", "rm", "couch_one")
//	exe_cmd("docker", "ps", "-a")
//	exe_cmd("docker", "run", "-d", "--name", "couch_one", "-p", "11210:11210", "-p", "8091:8091", "-p", "8092:8092", "-e", "CLUSTER_INIT_USER=Administrator", "-e", "CLUSTER_INIT_PASSWORD=password", "-e", "SAMPLE_BUCKETS=\\\"beer-sample\\\"", "jamesdbloom/couchbase")
//	exe_cmd("docker", "ps", "-a")
//}
//
