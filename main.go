package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	ftp "github.com/liulixiang1988/taurus-go/sftp"
)

var (
	USER = flag.String("user", os.Getenv("USER"), "ssh username")
	HOST = flag.String("host", "localhost", "ssh server hostname")
	PORT = flag.Int("port", 22, "ssh server port")
	PASS = flag.String("pass", os.Getenv("SOCKSIE_SSH_PASSWORD"), "ssh password")
	SIZE = flag.Int("s", 1<<15, "set max packet size")
)

func init() {
	flag.Parse()
}

func main() {
	var remoteClient, err = ftp.Connect(*USER, *PASS, *HOST, *PORT, 30)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	defer remoteClient.Close()

	err = remoteClient.SftpClient.Mkdir("liulx-aaa")
	if err != nil {
		fmt.Println("error occured when creating dir", err)
	}

	var stdOut, stdErr bytes.Buffer
	remoteClient.SshSession.Stdout = &stdOut
	remoteClient.SshSession.Stderr = &stdErr

	err = remoteClient.SshSession.Run("echo $PATH")
	if err != nil {
		fmt.Println("run command error:", err)
	}
	fmt.Println("result: ", stdOut.String())

	fmt.Printf("error: %s\n", stdErr.String())
}
