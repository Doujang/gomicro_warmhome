package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"gomicro_warmhome/homeweb/models"
	"gomicro_warmhome/homeweb/utils"

	"github.com/micro/go-log"

	example "gomicro_warmhome/GetUserInfo/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetUserInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取用户信息 url：api/v1.0/user")

	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//连接缓存根据sessionId查询用户信息
	bm, err := utils.GetRedisConnector()
	if err != nil {
		beego.Info("连接缓存失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	userInfo_redis := bm.Get(req.SessionId)
	//如果缓存中无数据
	if userInfo_redis == nil {
		beego.Info("缓存中无数据，session过期")
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	userInfo_string, _ := redis.String(userInfo_redis, nil)
	user := models.User{}
	err = json.Unmarshal([]byte(userInfo_string), &user)
	rsp.UserId = user.Uid
	rsp.Email = user.Email
	rsp.Name = user.Name
	rsp.RealName = user.Real_name
	rsp.IdCard = user.Id_card
	rsp.AvatarUrl = user.Avatar_url
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
