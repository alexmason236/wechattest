package main

import (
	"fmt"
	"github.com/l306287405/wechat3rd"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var srv *wechat3rd.Server

func main() {

	go func() {
		initRouter()
		log.Fatal(http.ListenAndServe(":80", nil))
	}()

	time.Sleep(10 * time.Second)
	service := getService()
	_ = service.ApiStartPushTicket()
}

func getService() *wechat3rd.Server {
	if srv != nil {
		return srv
	}
	service, err := wechat3rd.NewService(wechat3rd.Config{
		AppID:     "wx0775b18bb5d55acc",                          //第三方平台appid
		AppSecret: "fc38474968b3a00f3af02a2b4ff4e818",            //第三方平台app_secret
		AESKey:    "JcAKpeGTnPGSVCPPhYLgbFCWXENgIeDfeogXZbooLzo", //第三方平台消息加解密Key
		Token:     "wds_token",                                   //消息校验Token
	}, nil, nil, nil)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	srv = service
	return srv
}

func initRouter() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		service := getService()
		resp := service.PreAuthCode()
		if !resp.Success() {
			fmt.Println("获取授权链接失败:", resp.ErrMsg)
			return
		}

		r := url.Values{}

		// 必选参数
		r.Add("component_appid", os.Getenv("WX_OPEN_APP_ID"))
		r.Add("pre_auth_code", resp.PreAuthCode)
		r.Add("redirect_uri", "你的回调url")
		r.Add("auth_type", string(wechat3rd.PREAUTH_AUTH_TYPE_MINIAPP))

		// 网页方式授权：授权注册页面扫码授权
		authUrl := "https://mp.weixin.qq.com/cgi-bin/componentloginpage?"
		authUrl += r.Encode()
		println(authUrl)
	}
	http.HandleFunc("/wxcallback", helloHandler)
}
