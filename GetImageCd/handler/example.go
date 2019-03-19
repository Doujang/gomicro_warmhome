package handler

import (
	"context"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"gomicro_warmhome/homeweb/utils"
	"image/color"
	"time"

	"github.com/micro/go-log"

	example "gomicro_warmhome/GetImageCd/proto/example"
)

type Example struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) GetImageCd(ctx context.Context, req *example.Request, rsp *example.Response) error {
	beego.Info("获取验证码图片请求客户端 url:api/v1.0/imagecode/:uuid")
	//创建句柄
	cap := captcha.New()
	//通过句柄调用字体文件
	if err := cap.SetFont("comic.ttf"); err != nil {
		beego.Info("没有字体文件")
		panic(err.Error())
	}
	//设置图片大小
	cap.SetSize(91, 41)
	//设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	//SetFrontColor(colors ...color.Color)  这两个颜色设置的函数属于不定参函数
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 153, 0, 255})
	//生成图片 返回图片和 字符串(图片内容的文本形式)
	img, str := cap.Create(4, captcha.NUM)
	b := *img      //解引用
	c := *(b.RGBA) //解引用
	//默认返回成功
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	//图片信息
	rsp.Pix = []byte(c.Pix)
	rsp.Stride = int64(c.Stride)
	rsp.Max = &example.Response_Point{X: int64(c.Rect.Max.X), Y: int64(c.Rect.Max.Y)}
	rsp.Min = &example.Response_Point{X: int64(c.Rect.Min.X), Y: int64(c.Rect.Min.Y)}

	//将uuid与验证码存入redis
	bm, err := utils.GetRedisConnector()
	if err != nil {
		beego.Info("GetImageCd() cache.NewCache err", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	bm.Put(req.Uuid, str, time.Second*300)
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
