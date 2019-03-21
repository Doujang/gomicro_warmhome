package utils

import (
	"github.com/astaxie/beego"
	//使用了beego框架的配置文件读取模块
	"github.com/astaxie/beego/config"
)

var (
	G_server_name  string //项目名称
	G_server_addr  string //服务器ip地址
	G_server_port  string //服务器端口
	G_redis_addr   string //redis ip地址
	G_redis_port   string //redis port端口
	G_redis_dbnum  string //redis db 编号
	G_redis_passwd string //redis 密码
	G_mysql_addr   string //mysql ip 地址
	G_mysql_port   string //mysql 端口
	G_mysql_dbname string //mysql db name
	G_mysql_passwd string //mysql db password
	G_img_addr     string //图片服务器地址
	G_ssh_addr     string //ssh远程服务器地址
	G_email_user   string //邮箱账号
	G_email_passwd string //邮箱密码
)

func InitConfig() {
	//从配置文件读取配置信息
	//如果项目迁移需要进行修改
	appconf, err := config.NewConfig("ini", "D:/Projects/mygo/src/gomicro_warmhome/homeweb/conf/app.conf")
	if err != nil {
		beego.Debug(err)
		return
	}
	G_server_name = appconf.String("appname")
	G_server_addr = appconf.String("httpaddr")
	G_server_port = appconf.String("httpport")
	G_redis_addr = appconf.String("redisaddr")
	G_redis_port = appconf.String("redisport")
	G_redis_passwd = appconf.String("redispasswd")
	G_redis_dbnum = appconf.String("redisdbnum")
	G_mysql_addr = appconf.String("mysqladdr")
	G_mysql_port = appconf.String("mysqlport")
	G_mysql_dbname = appconf.String("mysqldbname")
	G_mysql_passwd = appconf.String("mysqlpasswd")
	G_img_addr = appconf.String("imgaddr")
	G_email_user = appconf.String("email_account")
	G_email_passwd = appconf.String("email_passwd")
	G_ssh_addr = appconf.String("sshaddr")
	return
}

func init() {
	InitConfig()
}
