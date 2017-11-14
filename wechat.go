package wechat

import (
	"net/http"
	"sync"

	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/context"
	"github.com/silenceper/wechat/js"
	"github.com/silenceper/wechat/material"
	"github.com/silenceper/wechat/menu"
	"github.com/silenceper/wechat/oauth"
	"github.com/silenceper/wechat/server"
	"github.com/silenceper/wechat/template"
	"github.com/silenceper/wechat/user"
)

// TokenOrTicketRefreshBufferPeriod is the duration before token/ticket
// actually expires that our cache is removed and server request is
// able to be triggered.
var TokenOrTicketRefreshBufferPeriod = 1500

// {Min|Max}imumCacheLife - if set to a positive number, this number will
// override the `expires' field returned by WeChat server and act as the
// actual TTL of entries in our database.
var (
	MinimumCacheLife = 0
	MaximumCacheLife = 0
)

// Wechat struct
type Wechat struct {
	Context *context.Context
}

// Config for user
type Config struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
	Cache          cache.Cache
}

// NewWechat init
func NewWechat(cfg *Config) *Wechat {
	context := new(context.Context)
	copyConfigToContext(cfg, context)
	return &Wechat{context}
}

func copyConfigToContext(cfg *Config, context *context.Context) {
	context.AppID = cfg.AppID
	context.AppSecret = cfg.AppSecret
	context.Token = cfg.Token
	context.EncodingAESKey = cfg.EncodingAESKey
	context.Cache = cfg.Cache
	context.SetAccessTokenLock(new(sync.RWMutex))
	context.SetJsAPITicketLock(new(sync.RWMutex))
}

// GetServer 消息管理
func (wc *Wechat) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	wc.Context.Request = req
	wc.Context.Writer = writer
	return server.NewServer(wc.Context)
}

//GetAccessToken 获取access_token
func (wc *Wechat) GetAccessToken() (string, error) {
	return wc.Context.GetAccessToken()
}

// GetOauth oauth2网页授权
func (wc *Wechat) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(wc.Context)
}

// GetMaterial 素材管理
func (wc *Wechat) GetMaterial() *material.Material {
	return material.NewMaterial(wc.Context)
}

// GetJs js-sdk配置
func (wc *Wechat) GetJs() *js.Js {
	return js.NewJs(wc.Context)
}

// GetMenu 菜单管理接口
func (wc *Wechat) GetMenu() *menu.Menu {
	return menu.NewMenu(wc.Context)
}

// GetUser 用户管理接口
func (wc *Wechat) GetUser() *user.User {
	return user.NewUser(wc.Context)
}

// GetTemplate 模板消息接口
func (wc *Wechat) GetTemplate() *template.Template {
	return template.NewTemplate(wc.Context)
}
