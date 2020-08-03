package switcher

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

func mikrotikPing(conf Config) (bool, error) {
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", conf.MikrotikAddr, conf.SshClientConfig)
	if err != nil {
		return false, err
	}
	defer client.Close()
	return true, nil
}

func mikrotikRunScript(conf Config, command string) (string, error) {
	res := ""
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", conf.MikrotikAddr, conf.SshClientConfig)
	if err != nil {
		return res, err
	}
	defer client.Close()
	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return res, err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)
	return stdoutBuf.String(), nil
}

