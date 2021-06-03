package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

func SSHClient(c *gin.Context) {
	// 建立SSH客户端连接
	client, err := ssh.Dial("tcp", "10.24.2.4:22", &ssh.ClientConfig{
		User:            "admin",
		Auth:            []ssh.AuthMethod{ssh.Password("admin@123")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		ResponseNotFoundCode(c, err.Error())
		return
	}
	// 建立新会话
	session, Err := client.NewSession()
	if Err != nil {
		ResponseNotFoundCode(c, Err.Error())
		return
	}
	defer session.Close()
	// 输入命令执行
	result, RErr := session.Output("sh spanning-tree vlan 25 detail | xml")
	if RErr != nil {
		ResponseNotFoundCode(c, RErr.Error())
		return
	}
	fmt.Println(string(result))
}

func Terminal(c *gin.Context) {
	// 建立SSH客户端连接
	client, err := ssh.Dial("tcp", "192.168.100.120:22", &ssh.ClientConfig{
		User:            "asialink",
		Auth:            []ssh.AuthMethod{ssh.Password("Admin@123")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		ResponseNotFoundCode(c, err.Error())
		return
	}
	// 建立新会话
	session, Err := client.NewSession()
	if Err != nil {
		ResponseNotFoundCode(c, Err.Error())
		return
	}
	defer session.Close()
	session.Stdout = os.Stdout // 会话输出关联到系统标准输出设备
	session.Stderr = os.Stderr // 会话错误输出关联到系统标准错误输出设备
	session.Stdin = os.Stdin   // 会话输入关联到系统标准输入设备
	modes := ssh.TerminalModes{
		ssh.ECHO:          0, // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("linux", 32, 160, modes); err != nil {
		log.Fatalf("request pty error: %s", err.Error())
	}
	if err = session.Shell(); err != nil {
		log.Fatalf("start shell error: %s", err.Error())
	}
	if err = session.Wait(); err != nil {
		log.Fatalf("return error: %s", err.Error())
	}
}
