package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"gomicro_warmhome/homeweb/models"
	"gomicro_warmhome/homeweb/utils"
	"math/rand"
	"strconv"
	"time"

	"github.com/micro/go-log"

	example "gomicro_warmhome/GetEmailcd/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetEmailCd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取邮箱验证码请求客户端 url:api/v1.0/emailcode/:email")
	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//验证邮箱是否存在
	o := orm.NewOrm()
	user := models.User{Email: req.Email}
	err := o.Read(&user)
	if err == nil {
		beego.Info("用户已经存在")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errmsg)
		return nil
	}
	//连接redis
	bm, err := utils.GetRedisConnector()
	if err != nil {
		beego.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	value := bm.Get(req.Uuid)
	if value == nil {
		beego.Info("缓存查询失败", value)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	value_str, _ := redis.String(value, nil)
	//校验验证码
	if req.Text != value_str {
		beego.Info("图片验证码错误")
		rsp.Errno = utils.RECODE_SMSERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code_number := r.Intn(9999) + 1001
	beego.Info(code_number)
	code := strconv.Itoa(code_number)
	//发送邮箱验证码
	err = utils.SendEmail(req.Email, code)
	if err != nil {
		beego.Info("邮件发送失败")
		rsp.Errno = utils.RECODE_SERVERERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	err = bm.Put(req.Email, code, time.Second*300)
	if err != nil {
		beego.Info("缓存异常")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Example) Stream(ctx context.Context, req *example.StreamingRequest, stream example.Example_StreamStream) error {
	log.Logf("Received Example.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&example.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Example) PingPong(ctx context.Context, stream example.Example_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&example.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
