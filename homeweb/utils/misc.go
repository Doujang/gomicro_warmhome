package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"path"
	"time"
)

func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_img_addr + "/" + url

	return domain_url
}

func Md5String(s string) string {
	//创建1个md5对象
	h := md5.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

func Sha256Encode(value string) string {
	encoder := sha256.New()
	encoder.Write([]byte(value))
	hash := encoder.Sum(nil)
	result := hex.EncodeToString(hash)
	return string(result)
}

func SendEmail(emailTo string, code string) error {
	//异常捕获
	defer func() {
		if err := recover(); err != nil {
			beego.Info("邮件发送失败")
		} else {
			beego.Info("邮件发送成功")
		}
	}()

	config := `{"username":"` + G_email_user + `","password":"` + G_email_passwd + `","host":"smtp.163.com","port":25}`
	beego.Info(config)
	temail := utils.NewEMail(config)
	//指定邮件基本信息
	//收件人
	temail.To = []string{emailTo}
	//发件人
	temail.From = "温暖小家"
	//邮件主题
	temail.Subject = "温暖小家注册验证码"
	//邮件内容
	temail.HTML = "欢迎注册温暖小家，您的验证码为:" + code
	err := temail.Send()
	return err
}

func connect() (*sftp.Client, error) {
	var (
		sftpClient *sftp.Client
		err        error
	)
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	pemBytes, err := ioutil.ReadFile("D:/Projects/mygo/src/gomicro_warmhome/homeweb/conf/id_rsa")
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 30 * time.Second,
	}
	conn, err := ssh.Dial("tcp", G_ssh_addr+":22", config)
	if err != nil {
		log.Fatalf("dial failed:%v", err)
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(conn); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

//上传二进制文件到fdfs中的操作
func UploadByBuffer(filebuffer []byte, fileExt string) (string, error) {
	var (
		err        error
		sftpClient *sftp.Client
	)

	sftpClient, err = connect()
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	// 用来测试的本地文件路径 和 远程机器上的文件夹
	var remoteDir = "/home/vsftpd/sher/"
	var fileName = Md5String(time.Now().String())
	var remoteFileName = fileName + "." + fileExt
	fmt.Println(remoteFileName)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()
	dstFile.Write(filebuffer)
	//返回远程文件名称
	return remoteFileName, nil
}
