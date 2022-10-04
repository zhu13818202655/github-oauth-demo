package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type GitHubConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
//
var conf = GitHubConfig{
	ClientId: "da0e42c8ecdec813xxxx",
	ClientSecret: "522099cc4d945eaf63ea655347e663a0acc2xxxx",
	RedirectUrl: "http://localhost:8080/authorization",
}

func main() {

	engine := gin.Default()
	engine.LoadHTMLGlob("html/*")

	engine.GET("/", func(c *gin.Context) {
		c.HTML(200, "githubLogin.tmpl", conf)
	})
	//1. 认证服务器在发行 access_token 之前要先通过用户的同意, 这里会自动调用回调函数
	engine.GET("authorization", func(c *gin.Context) {
		code, _ := c.GetQuery("code")
		if code != "" {
			token, err := getToken(code)
			if err != nil {
				panic(err)
			}
			info, err := getUserInfo(token.AccessToken)
			if err != nil {
				panic(err)
			}
			fmt.Println(info)
			c.String(200, info)
		} else {
			c.String(500, "nil")
		}
	})

	engine.Run(":8080")

}

//2. 认证服务器负责生成并且发行 access_token 给第三方应用程序
func getToken(code string) (*TokenResponse, error) {
	var token TokenResponse
	client := &http.Client{}
	request, err := http.NewRequest("", fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", conf.ClientId, conf.ClientSecret, code), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func getUserInfo(token string) (string, error) {
	url := "https://api.github.com/user"
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "token " + token) //3. 拿到token凭证请求资源服务器
	request.Header.Add("Accept", "application/json")
	resp1, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp1.Body.Close()
	all, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		panic(err)
	}
	return string(all), nil
}
