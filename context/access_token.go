package context

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/util"
)

const (
	//AccessTokenURL 获取access_token的接口
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"
)

//ResAccessToken struct
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

//SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) SetAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

//GetAccessToken 获取access_token
func (ctx *Context) GetAccessToken() (accessToken string, err error) {
	ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", ctx.AppID)
	val := ctx.Cache.Get(accessTokenCacheKey)
	if val != nil {
		accessToken = val.(string)
		return
	}

	//从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = ctx.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (ctx *Context) GetAccessTokenFromServer() (resAccessToken ResAccessToken, err error) {
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", AccessTokenURL, ctx.AppID, ctx.AppSecret)
	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	log.Println("GetAccessTokenFromServer body:", string(body))
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrMsg != "" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", ctx.AppID)

	expires := resAccessToken.ExpiresIn - wechat.TokenOrTicketRefreshBufferPeriod
	if !(wechat.MinimumCacheLife > 0 && wechat.MaximumCacheLife > 0 && wechat.MinimumCacheLife > wechat.MaximumCacheLife) {
		if wechat.MinimumCacheLife > 0 && expires < wechat.MinimumCacheLife {
			expires = wechat.MinimumCacheLife
		}
		if wechat.MaximumCacheLife > 0 && expires > wechat.MaximumCacheLife {
			expire = wechat.MinimumCacheLife
		}
	}
	err = ctx.Cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	return
}
