package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"github.com/micro/go-log"
	"gomicro_warmhome/homeweb/models"
	"gomicro_warmhome/homeweb/utils"
	"time"

	example "gomicro_warmhome/PostUserAuth/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostUserAuth(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info(" 实名认证 Postuserauth  api/v1.0/user/auth ")

	//创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*从session中获取我们的user_id*/
	//连接redis数据库
	bm, err := utils.GetRedisConnector()

	userInfo_redis := bm.Get(req.SessionId)
	userInfo_string, _ := redis.String(userInfo_redis, nil)
	userOld := models.User{}
	json.Unmarshal([]byte(userInfo_string), &userOld)

	//创建user对象
	user := models.User{Uid: userOld.Uid, Real_name: req.RealName, Id_card: req.IdCard}
	/*更新user表中的 姓名和 身份号*/
	o := orm.NewOrm()
	//更新表
	_, err = o.Update(&user, "real_name", "id_card")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//更新缓存
	userOld.Real_name = req.RealName
	userOld.Id_card = req.IdCard
	userInfo, _ := json.Marshal(userOld)
	bm.Put(req.SessionId, userInfo, time.Second*600)
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
