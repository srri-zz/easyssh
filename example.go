package main

import (
	"flag"
	"fmt"
	"github.com/stevenbrichards/easyssh"
)

//Usage: either pass pswd or keypath to id_rsa
//Example: easyssh -user root -pswd yourPswd -host 192.168.56.101 -cmd uptime
//Example: easyssh -user root -keypath id_rsa -host 192.168.56.101 -cmd uptime

var user = flag.String("user", "", "username")
var pswd = flag.String("pswd", "", "password")
var host = flag.String("host", "", "host server or IP address")
var cmd = flag.String("cmd", "", "command to run")
var keypath = flag.String("keypath", "", "keypath to id_rsa")

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
}
