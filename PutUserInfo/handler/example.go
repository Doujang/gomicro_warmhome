package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"gomicro_warmhome/homeweb/models"
	"gomicro_warmhome/homeweb/utils"
	"time"

	"github.com/micro/go-log"

	example "gomicro_warmhome/PutUserInfo/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PutUserInfo(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("修改用户名 url：api/v1.0/user/name")

	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//连接缓存
	bm, err := utils.GetRedisConnector()
	if err != nil {
		beego.Info("缓存连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//从缓存中拿到用户数据
	userInfo_redis := bm.Get(req.SessionId)
	userInfo_string, _ := redis.String(userInfo_redis, nil)
	userOld := models.User{}
	json.Unmarshal([]byte(userInfo_string), &userOld)
	user := models.User{Uid: userOld.Uid, Name: req.Username}
	//更新数据库
	o := orm.NewOrm()
	_, err = o.Update(&user, "name")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//更新缓存
	userOld.Name = req.Username
	userInfo, _ := json.Marshal(userOld)
	bm.Put(req.SessionId, userInfo, time.Second*600)
	//返回数据
	rsp.Username = user.Name
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
