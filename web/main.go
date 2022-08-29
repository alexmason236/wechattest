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
	intChan := make(chan int)
	go func() {
		initRouter()
		log.Fatal(http.ListenAndServe(":80", nil))
	}()

	time.Sleep(10 * time.Second)
	service := getService()
	_ = service.ApiStartPushTicket()
	<-intChan
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

	ticketHandler := func(w http.ResponseWriter, req *http.Request) {
		wechatService := getService()
		resp, err := wechatService.ServeHTTP(req)
		if err != nil {
			fmt.Println("微信第三方开放平台component_verify_ticket获取失败:", err.Error())
			return
		}
		// 将ticket缓存,并在服务重启时取用.
		err = wechatService.SetTicket(resp.ComponentVerifyTicket)

		if err != nil {
			fmt.Println("微信第三方开放平台component_verify_ticket设置失败:", err.Error())
		}
		h := &HtmlWriter{rw: w}
		h.html(200, "success")
	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/wxcallback", ticketHandler)
}

type HtmlWriter struct {
	rw http.ResponseWriter
}

func (hw *HtmlWriter) html(code int, html string) {
	hw.rw.Header().Set("Content-Type", "text/html")
	hw.rw.WriteHeader(code)
	hw.rw.Write([]byte(html))
}
