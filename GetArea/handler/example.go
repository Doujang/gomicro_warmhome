package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"gomicro_warmhome/homeweb/models"
	"gomicro_warmhome/homeweb/utils"
	"time"

	"github.com/micro/go-log"

	example "gomicro_warmhome/GetArea/proto/example"
)

type Example struct{}

//获取地区数据
func (e *Example) GetArea(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取地区请求客户端 url:api/v1.0/areas")

	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//连接redis
	bm, err := utils.GetRedisConnector()
	if err != nil {
		beego.Info("连接缓存失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//1.获取缓存数据
	areas_info_value := bm.Get("areas_info")
	//a.缓存中有数据
	if areas_info_value != nil {
		//存放解码后的json数据
		areas_info := []map[string]interface{}{}
		//解码
		err = json.Unmarshal(areas_info_value.([]byte), &areas_info)
		//进行赋值
		for key, value := range areas_info {
			beego.Info(key, value)
			area := example.Response_Address{Aid: int32(value["aid"].(float64)), Aname: value["aname"].(string)}
			rsp.Data = append(rsp.Data, &area)
		}
		return nil
	}
	//b.缓存中无数据，我们需要从mysql中读取并且加载到redis中
	o := orm.NewOrm()
	//用来接收数据
	var areas []models.Area
	//查询area表
	qs := o.QueryTable("area")
	//查询全部区域数据
	num, err := qs.All(&areas)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//无数据
	if num == 0 {
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//写入缓存
	area_json, _ := json.Marshal(areas)
	err = bm.Put("areas_info", area_json, time.Second*3600)
	if err != nil {
		beego.Info("区域数据写入缓存失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	//返回区域数据
	for _, value := range areas {
		area := example.Response_Address{Aid: int32(value.Id), Aname: value.Name}
		rsp.Data = append(rsp.Data, &area)
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
