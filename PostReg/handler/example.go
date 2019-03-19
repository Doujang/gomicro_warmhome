package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
	"github.com/micro/go-log"
	"gomicro_warmhome/homeweb/models"
	"gomicro_warmhome/homeweb/utils"
	"time"

	example "gomicro_warmhome/PostReg/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostReg(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("注册请求  /api/v1.0/users")

	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	bm, err := utils.GetRedisConnector()
	if err != nil {
		beego.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//查询相关数据
	code_redis := bm.Get(req.Email)
	if code_redis == nil {
		beego.Info("缓存数据为空")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//拿到邮件验证码
	code, _ := redis.String(code_redis, nil)
	//如果邮件验证码错误
	if req.EmailCode != code {
		beego.Info("邮件验证码错误")
		rsp.Errno = utils.RECODE_SMSERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	user := models.User{}
	user.Uid = uuid.New().String()
	user.Name = req.Email
	pwd_hash := utils.Sha256Encode(req.Password)
	user.Password_hash = pwd_hash
	user.Email = req.Email
	beego.Info(user.Uid)
	o := orm.NewOrm()
	_, err = o.Insert(&user)
	if err != nil {
		beego.Info("数据插入失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//返回的sessionid
	sessionId := utils.Sha256Encode(user.Password_hash)
	rsp.SessionID = sessionId
	user.Password_hash = ""
	userInfo, _ := json.Marshal(user)
	bm.Put(sessionId, userInfo, time.Second*3600)
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
