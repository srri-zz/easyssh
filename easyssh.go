/*
	Wrapper for SSH in Go
	Copyright Steven Richards 2015 - <sbrichards@mit.edu>
*/

package easyssh

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/howeyc/gopass"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

type Config struct {
	User     string
	Password string
	KeyPath  string
	Host     string
	Port     string
	Client   *ssh.Client
}

func getKeyFile(keyPath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	p, rest := pem.Decode(buf)
	if len(rest) > 0 {
		return nil, fmt.Errorf("Failed to decode key")
	}
	pBlock := pem.Block{
		Bytes:   buf,
		Type:    p.Type,
		Headers: p.Headers,
	}
	if x509.IsEncryptedPEMBlock(&pBlock) {
		fmt.Println("Detected a password protected key file")
		decodedKeyBytes, err := tryDecrypt(p, 0)
		if err != nil {
			return nil, err
		}
		// Encode decrypted bytes into decoded PEM structure as if it were in a file
		pBlock = pem.Block{
			Bytes:   decodedKeyBytes,
			Type:    p.Type,
			Headers: p.Headers,
		}
		// Extract SSH key from decrypted PEM in-memory file
		rawkey, err := ssh.ParsePrivateKey(pem.EncodeToMemory(&pBlock))
		if err != nil {
			return nil, err
		}
		return rawkey.(ssh.Signer), nil
	}
	// Non-encrypted key
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func tryDecrypt(p *pem.Block, count int) ([]byte, error) {
	if count >= 3 {
		return nil, fmt.Errorf("Too many tries")
	} else if count >= 1 {
		fmt.Println("Wrong password, try again")
	}
	fmt.Println("Enter the key's password: ")
	password := gopass.GetPasswd()
	// Decrypt using user's key
	decodedKeyBytes, err := x509.DecryptPEMBlock(p, password)
	if err != nil {
		return tryDecrypt(p, count+1)
	}
	return decodedKeyBytes, nil
}

// Connect takes no arguments and attempts to open an SSH connection.
// It uses the easyssh Config and the parameters it contains.
func (e *Config) Connect() error {
	var sshConfig = &ssh.ClientConfig{}
	var err error
	if e.KeyPath != "" {
		keyFile, err := getKeyFile(e.KeyPath)
		if err != nil {
			return fmt.Errorf("Unable to use specified key file")
		}
		sshConfig = &ssh.ClientConfig{
			User: e.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(keyFile),
			},
		}
	} else if e.Password != "" {
		sshConfig = &ssh.ClientConfig{
			User: e.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(e.Password),
			},
		}
	}
	e.Client, err = ssh.Dial("tcp", e.Host+":"+e.Port, sshConfig)
	return err
}

// Command executes a single string command.
// To do so it opens a new SSH session, executes the command, and returns stdout as a string.
// If it fails for any reasion it returns an empty string and an error.
func (e *Config) Command(command string) (string, error) {
	session, err := e.Client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Failed to create session: %s\n" + err.Error())
	}
	var stdout bytes.Buffer
	session.Stdout = &stdout
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("Failed to execute command: %s\n" + err.Error())
	}
	return stdout.String(), nil
}

// BatchCommands takes a string slice of commands and a string delimiter.
// The delimiter is appended to the result of each command, it can be an empty string as well.
// It returns whether there were any errors or not and outputs if there are any errors as it runs.
func (e *Config) BatchCommands(commands []string, delimiter string) ([]string, error) {
	results := make([]string, len(commands))
	errors := false
	for _, command := range commands {
		session, err := e.Client.NewSession()
		if err != nil {
			fmt.Printf("Failed to create session: %s\n" + err.Error())
			return nil, err
		}
		var stdout bytes.Buffer
		session.Stdout = &stdout
		err = session.Run(command)
		if err != nil {
			fmt.Printf("Failed to execute %s : %s\n", command, err.Error())
			errors = true
			continue
		}
		output := strings.TrimSpace(stdout.String())
		results = append(results, output)+delimiter)
	}
	if errors {
		return results, fmt.Errorf("There were some errors while executing the commands")
	}
	return results, nil
}
