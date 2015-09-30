package main

import (
	"flag"
	"fmt"
	"github.com/stevenbrichards/easyssh"
)

var user = flag.String("user", "", "username")
var pswd = flag.String("pswd", "", "password")
var host = flag.String("host", "", "host server or IP address")
var cmd = flag.String("cmd", "", "command to run")
var keypath = flag.String("keypath", "", "keypath")

func main() {
	flag.Parse()

	sshSession := easyssh.Config{
		User:     *user,
		Password: *pswd,
		KeyPath:  *keypath,
		Host:     *host,
		Port:     "22",
	}
	err := sshSession.Connect()
	if err != nil {
		panic(err)
	}
	result, err := sshSession.Command(*cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	fmt.Scanln()
}
