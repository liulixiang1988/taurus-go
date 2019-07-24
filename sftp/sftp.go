package sftp

import (
	"fmt"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type remoteClient struct {
	sshClient  *ssh.Client
	SftpClient *sftp.Client
	SshSession *ssh.Session
}

func (client remoteClient) Close() {
	if client.SshSession != nil {
		if err := client.SshSession.Close(); err != nil {
			fmt.Printf("err occured when close ssh session: %v\n", err)
		} else {
			fmt.Println("close ssh session success")
		}
	}

	if client.SftpClient != nil {
		err := client.SftpClient.Close()
		if err != nil {
			fmt.Printf("error occured when close sftp client: %v\n", err)
		} else {
			fmt.Println("close sftp success")
		}
	}

	if client.sshClient != nil {
		err := client.sshClient.Close()
		if err != nil {
			fmt.Printf("err occured when close ssh client: %v\n", err)
		} else {
			fmt.Println("close ssh success")
		}
	}
}

// Connect 连接sftp
func Connect(user, password, host string, port int, timeout time.Duration) (*remoteClient, error) {
	auth := []ssh.AuthMethod{ssh.Password(password)}
	clientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         timeout * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	session, err := sshClient.NewSession()
	if err != nil {
		sshClient.Close()
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, err
	}

	return &remoteClient{sshClient, sftpClient, session}, nil
}
