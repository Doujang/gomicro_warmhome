package handler

import (
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro/client"
	"gomicro_warmhome/homeweb/models"
	"gomicro_warmhome/homeweb/utils"
	"image"
	"image/png"
	"net/http"
	"reflect"
	"time"

	example "github.com/micro/examples/template/srv/proto/example"
	DELETESESSION "gomicro_warmhome/DeleteSession/proto/example"
	GETAREA "gomicro_warmhome/GetArea/proto/example"
	GETEMAILCD "gomicro_warmhome/GetEmailcd/proto/example"
	GETIMAGECD "gomicro_warmhome/GetImageCd/proto/example"
	GETSESSION "gomicro_warmhome/GetSession/proto/example"
	GETUSERINFO "gomicro_warmhome/GetUserInfo/proto/example"
	POSTLOGIN "gomicro_warmhome/PostLogin/proto/example"
	POSTREG "gomicro_warmhome/PostReg/proto/example"
)

func ExampleCall(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	exampleClient := example.NewExampleService("go.micro.srv.template", client.DefaultClient)
	rsp, err := exampleClient.Call(context.TODO(), &example.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

//获取地区信息
func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("GetArea url:api/v1.0/areas")

	//创建新grpc返回句柄
	server := grpc.NewService()
	//服务初始化
	server.Init()

	//创建获取地区的服务并返回句柄
	exampleClient := GETAREA.NewExampleService("go.micro.srv.GetArea", server.Client())

	//调用函数并且获得返回数据
	rsp, err := exampleClient.GetArea(context.TODO(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//接收数据
	//准备接收切片
	area_list := []models.Area{}
	//循环接收数据
	for _, value := range rsp.Data {
		tmp := models.Area{Id: int(value.Aid), Name: value.Aname}
		area_list = append(area_list, tmp)
	}

	// 返回给前端的map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   area_list,
	}

	//会传数据的时候三直接发送过去的并没有设置数据格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImageCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("获取验证码图片请求客户端 url:api/v1.0/imagecode/:uuid")

	//创建服务
	server := grpc.NewService()
	//服务初始化
	server.Init()

	//获取前端传送过来的uuid
	exampleClient := GETIMAGECD.NewExampleService("go.micro.srv.GetImageCd", server.Client())
	rsp, err := exampleClient.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid: ps.ByName("uuid"),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//处理从服务端传送过来的图片信息
	var img image.RGBA
	img.Stride = int(rsp.Stride)
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)
	img.Pix = []uint8(rsp.Pix)

	var image captcha.Image
	image.RGBA = &img

	//将图片发送给前端
	png.Encode(w, image)

}

func GetEmailCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("获取邮箱验证码请求客户端 url:api/v1.0/emailcode/:email")
	//创建服务并初始化
	server := grpc.NewService()
	server.Init()

	//获取前端发送过来的邮箱号
	email := ps.ByName("email")
	beego.Info(email)

	beego.Info(r.URL.Query())
	//获取url携带的图片验证码和uuid
	text := r.URL.Query()["text"][0]
	id := r.URL.Query()["id"][0]

	//调用服务
	exampleClient := GETEMAILCD.NewExampleService("go.micro.srv.GetEmailcd", server.Client())
	rsp, err := exampleClient.GetEmailCd(context.TODO(), &GETEMAILCD.Request{
		Email: email,
		Uuid:  id,
		Text:  text,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置返回格式
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostReg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("注册请求  /api/v1.0/users")

	//解析前端发送过来的json数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for key, value := range request {
		beego.Info(key, value, reflect.TypeOf(value))
	}

	if request["email"] == "" || request["password"] == "" || request["email_code"] == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误请重新输入",
		}
		w.Header().Set("Content-type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		beego.Info("有数据为空")
		return
	}

	//创建服务并初始化
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := POSTREG.NewExampleService("go.micro.srv.PostReg", server.Client())
	rsp, err := exampleClient.PostReg(context.TODO(), &POSTREG.Request{
		Email:     request["email"].(string),
		Password:  request["password"].(string),
		EmailCode: request["email_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//读取cookie
	cookie, err := r.Cookie("userlogin")
	//如果读取失败或者cookie中的value不存在则创建cookie
	if err != nil || "" == cookie.Value {
		cookie := http.Cookie{Name: "userlogin", Value: rsp.SessionID, Path: "/", MaxAge: 600}
		http.SetCookie(w, &cookie)
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取Session url：api/v1.0/session")

	//创建服务并初始化
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETSESSION.NewExampleService("go.micro.srv.GetSession", server.Client())

	//获取cookie
	userlogin, err := r.Cookie("userlogin")
	//未登录或登录超时
	if err != nil || "" == userlogin.Value {
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	//如果cookie有值就发送到服务端
	rsp, err := exampleClient.GetSession(context.TODO(), &GETSESSION.Request{
		SessionId: userlogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := make(map[string]string)
	data["name"] = rsp.Data
	//创建返回数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("登录请求  /api/v1.0/sessions")

	//解析前端发送过来的json数据
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for key, value := range request {
		beego.Info(key, value, reflect.TypeOf(value))
	}

	if request["email"] == "" || request["password"] == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误请重新输入",
		}
		w.Header().Set("Content-type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		beego.Info("有数据为空")
		return
	}

	//创建服务并初始化
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := POSTLOGIN.NewExampleService("go.micro.srv.PostLogin", server.Client())
	rsp, err := exampleClient.PostLogin(context.TODO(), &POSTLOGIN.Request{
		Email:    request["email"].(string),
		Password: request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//读取cookie
	cookie, err := r.Cookie("userlogin")
	//如果读取失败或者cookie中的value不存在则创建cookie
	if err != nil || "" == cookie.Value {
		cookie := http.Cookie{Name: "userlogin", Value: rsp.SessionId, Path: "/", MaxAge: 600}
		http.SetCookie(w, &cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func DeleteSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	beego.Info("登出请求  /api/v1.0/session")

	//创建服务并初始化
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := DELETESESSION.NewExampleService("go.micro.srv.DeleteSession", server.Client())
	//获取cookie
	userlogin, err := r.Cookie("userlogin")
	//Cookie为空
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}

	rsp, err := exampleClient.DeleteSession(context.TODO(), &DELETESESSION.Request{
		SessionId: userlogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//读取cookie
	cookie, err := r.Cookie("userlogin")
	//如果读取失败或者cookie中的value不存在则创建cookie
	if err != nil || "" == cookie.Value {
		return
	} else {
		cookie := http.Cookie{Name: "userlogin", Path: "/", MaxAge: 600}
		http.SetCookie(w, &cookie)
	}

	// we want to augment the response
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	w.Header().Set("Content-type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取用户信息 url：api/v1.0/user")

	//初始化服务
	server := grpc.NewService()
	server.Init()

	// call the backend service
	exampleClient := GETUSERINFO.NewExampleService("go.micro.srv.GetUserInfo", server.Client())

	userlogin, err := r.Cookie("userlogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		// encode and write the response as json
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			beego.Info(err)
			return
		}
		return
	}
	//成功就将信息发送给前端
	rsp, err := exampleClient.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		SessionId: userlogin.Value,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}
	//准备数据
	data := make(map[string]interface{})
	//将信息发送给前端
	data["user_id"] = rsp.UserId
	data["name"] = rsp.Name
	data["email"] = rsp.Email
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IdCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}
	//设置格式
	w.Header().Set("Content-Type", "application/json")
	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}
	return
}

func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取首页轮播 url：api/v1.0/houses/index")

	//创建返回数据map
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}
	w.Header().Set("Content-Type", "application/json")

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
