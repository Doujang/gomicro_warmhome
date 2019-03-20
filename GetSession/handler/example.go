package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"gomicro_warmhome/homeweb/models"
	"gomicro_warmhome/homeweb/utils"

	"github.com/micro/go-log"

	example "gomicro_warmhome/GetSession/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetSession(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取Session url：api/v1.0/session")

	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//连接缓存
	bm, err := utils.GetRedisConnector()
	if err != nil {
		beego.Info("获取缓存连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//从缓存中拿到用户信息
	beego.Info(req.SessionId)
	userInfo_redis := bm.Get(req.SessionId)
	userInfo_string, _ := redis.String(userInfo_redis, nil)
	beego.Info(userInfo_string)
	userInfo := []byte(userInfo_string)
	user := models.User{}
	err = json.Unmarshal(userInfo, &user)
	if err != nil {
		beego.Info("Json解析异常")
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.Data = user.Name

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
