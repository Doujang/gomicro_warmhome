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
	"path"
	"time"

	example "gomicro_warmhome/PostAvatar/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) PostAvatar(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("上传头像请求  /api/v1.0/avatar")

	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//查看数据是否正常
	beego.Info(len(req.Avatar), req.Filesize)
	//获取文件后缀名
	fileExt := path.Ext(req.Filename)
	//上传数据
	filename, err := utils.UploadByBuffer(req.Avatar, fileExt[1:])
	if err != nil {
		beego.Info("图片上传到图片服务器过程中出错")
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//连接缓存
	bm, err := utils.GetRedisConnector()
	if err != nil {
		beego.Info("缓存连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//从缓存中拿到用户uid,并更新数据库
	userOld := models.User{}

	userInfo_redis := bm.Get(req.SessionId)
	userInfo_string, _ := redis.String(userInfo_redis, nil)
	json.Unmarshal([]byte(userInfo_string), &userOld)
	user := models.User{Uid: userOld.Uid, Avatar_url: filename}
	//更新数据库
	o := orm.NewOrm()
	_, err = o.Update(&user, "avatar_url")
	//更新缓存
	userOld.Avatar_url = filename
	userInfo, _ := json.Marshal(userOld)
	bm.Put(req.SessionId, userInfo, time.Second*600)
	rsp.AvatarUrl = filename
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
