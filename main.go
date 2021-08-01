package ACGImage

import (
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	//r18有一定保护，一般不会发出图片
	RANDOM_API_URL = "&loli=true&r18=true"
	BLOCK_REQUEST  = false
	msgof          = make(map[int64]int64)
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
				if RANDOM_API_URL == "" {
					go Classify(ctx, RANDOM_API_URL, false)
				}
			}
		})
	// 直接随机图片，无r18保护，后果自负。如果出r18图可尽快通过发送"太涩了"撤回
	zero.OnFullMatch("直接随机", zero.AdminPermission).SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID > 0 {
				if BLOCK_REQUEST {
					ctx.Send("请稍后再试哦")
				} else if RANDOM_API_URL != "" {
					BLOCK_REQUEST = true
					setLastMsg(ctx.Event.GroupID, ctx.Send(message.Image(RANDOM_API_URL)))
					BLOCK_REQUEST = false
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
				ctx.Send("少女祈祷中......")
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
