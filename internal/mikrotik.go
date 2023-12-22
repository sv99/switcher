package internal

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

func MikrotikPing(conf *Config) (bool, error) {
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", conf.MikrotikAddr, conf.SshClientConfig)
	if err != nil {
		return false, err
	}
	defer client.Close()
	return true, nil
}

// MikrotikRunScript connect to the router by ssh and execute command.
func MikrotikRunScript(conf *Config, command string) (string, error) {
	res := ""
	client, err := ssh.Dial("tcp", conf.MikrotikAddr, conf.SshClientConfig)
	if err != nil {
		return res, err
	}
	defer func(client *ssh.Client) {
		_ = client.Close()
	}(client)

	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return res, err
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(command)
	if err != nil {
		return "", err
	}
	return stdoutBuf.String(), nil
}

// MikrotikGetVersion get router version
func MikrotikGetVersion(conf *Config) (string, error) {
	return MikrotikRunScript(conf, "/system script run get_version")
}

// MikrotikGetProvider get current provider
func MikrotikGetProvider(conf *Config) (string, error) {
	return MikrotikRunScript(conf, "/system script run get_provider")
}

// MikrotikSwitchProvider circle switch provider
func MikrotikSwitchProvider(conf *Config) (string, error) {
	return MikrotikRunScript(conf, "/system script run switch_provider")
}
