# easyssh
SSH wrapper for Go

## Features:  
* Supports PEM key authentication (password protected, and non-password protected)
* Supports password based authentication
* Batch and single command processing

## Docs:
http://godoc.org/github.com/stevenbrichards/easyssh

## Usage:  

`go get github.com/stevenbrichards/easyssh`

Single command:

```
sshSession := easyssh.Config{
	User:    username,
	KeyPath: pathToFile,
	Host:    host,
	Port:    "22",
}
err := sshSession.Connect()
if err != nil {
	panic(err)
}
result, err := sshSession.Command("whoami")
if err != nil {
	panic(err)
}
fmt.Println(result)
```

Batch commands

```
sshSession := easyssh.Config{
        User:    username,
        KeyPath: pathToFile,
        Host:    host,
        Port:    "22",
}
err := sshSession.Connect()
if err != nil {
        panic(err)
}
results, err := sshSession.BatchCommands([]string{"whoami","uptime"},"\n")
if err != nil {
        panic(err)
}
for _, result := range results {
	fmt.Println(result)
}
```


