package ACGImage

import (
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	LOLI_PROXY_URL = "http://saki.fumiama.top:62002/dice?class=0&loli=true&r18=true"
	//r18有一定保护，一般不会发出图片
	RANDOM_API_URL = "&loli=true&r18=true"
	msgof          = make(map[int64]int64)
	lastvisit      = time.Now().Unix()
)

func init() { // 插件主体
	zero.OnRegex(`^设置随机图片网址(.*)$`, zero.SuperUserPermission).SetBlock(true).SetPriority(20).
		Handle(func(ctx *zero.Ctx) {
			url := ctx.State["regex_matched"].([]string)[1]
			if !strings.HasPrefix(url, "http") {
				ctx.Send("URL非法!")
			} else {
				RANDOM_API_URL = url
			}
		})
	// 有保护的随机图片
	zero.OnFullMatch("随机图片").SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID > 0 {
				if time.Now().Unix()-lastvisit > 5 {
					go Classify(ctx, RANDOM_API_URL, false)
					lastvisit = time.Now().Unix()
				} else {
					ctx.Send("你太快啦!")
				}
			}
		})
	// 直接随机图片，无r18保护，后果自负。如果出r18图可尽快通过发送"太涩了"撤回
	zero.OnFullMatch("直接随机", zero.AdminPermission).SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID > 0 {
				if time.Now().Unix()-lastvisit > 5 {
					ctx.Send("请稍后再试哦")
				} else if RANDOM_API_URL != "" {
					var url string
					if RANDOM_API_URL[0] == '&' {
						url = LOLI_PROXY_URL
					} else {
						url = RANDOM_API_URL
					}
					setLastMsg(ctx.Event.GroupID, ctx.Send(message.Image(url).Add("cache", "0")))
					lastvisit = time.Now().Unix()
				}
			}
		})
	// 撤回最后的直接随机图片
	zero.OnFullMatch("太涩了").SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			go cancel(ctx)
		})
	// 上传一张图进行评价
	zero.OnFullMatch("评价图片", MustHasPicture()).SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID > 0 {
				ctx.Send("少女祈祷中...")
				for _, pic := range ctx.State["image_url"].([]string) {
					//fmt.Println(pic)
					go Classify(ctx, pic, true)
				}
			}
		})
}

func setLastMsg(id int64, msg int64) {
	msgof[id] = msg
}

func cancel(ctx *zero.Ctx) {
	msg, ok := msgof[ctx.Event.GroupID]
	if ok {
		ctx.DeleteMessage(msg)
		delete(msgof, ctx.Event.GroupID)
	}
}
