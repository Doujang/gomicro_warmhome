package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
)

func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_img_addr + "/" + url

	return domain_url
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
