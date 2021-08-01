package ACGImage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/tidwall/gjson"
)

var (
	BOTPATH, _     = os.Getwd() // 当前bot运行目录
	DATAPATH       = BOTPATH + "/data/acgimage/"
	CACHE_IMG_FILE = DATAPATH + "cache"
	CACHE_URI      = "file:///" + CACHE_IMG_FILE
	CLASSIFY_HEAD  = "http://saki.fumiama.top:62002/dice?class=9&url="
)

func init() {
	os.RemoveAll(DATAPATH) //清除缓存
	err := os.MkdirAll(DATAPATH, 0755)
	if err != nil {
		panic(err)
	}
}

func Classify(ctx *zero.Ctx, targeturl string, noimg bool) {
	if targeturl[0] != '&' {
		targeturl = url.QueryEscape(targeturl)
	}
	get_url := CLASSIFY_HEAD + targeturl
	if noimg {
		get_url += "&noimg=true"
	}
	resp, err := http.Get(get_url)
	if err != nil {
		ctx.Send(fmt.Sprintf("ERROR: %v", err))
	} else {
		if noimg {
			data, err1 := ioutil.ReadAll(resp.Body)
			if err1 == nil {
				dhash := gjson.GetBytes(data, "img").String()
				class := int(gjson.GetBytes(data, "class").Int())
				replyClass(ctx, dhash, class, noimg)
			} else {
				ctx.Send(fmt.Sprintf("ERROR: %v", err1))
			}
		} else {
			class, err1 := strconv.Atoi(resp.Header.Get("Class"))
			dhash := resp.Header.Get("DHash")
			if err1 != nil {
				ctx.Send(fmt.Sprintf("ERROR: %v", err1))
			}
			defer resp.Body.Close()
			// 写入文件
			data, _ := ioutil.ReadAll(resp.Body)
			f, _ := os.OpenFile(CACHE_IMG_FILE+strconv.FormatInt(lastvisit, 10), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			defer f.Close()
			f.Write(data)
			replyClass(ctx, dhash, class, noimg)
		}
	}
}

func replyClass(ctx *zero.Ctx, dhash string, class int, noimg bool) {
	if class > 5 {
		switch class {
		case 6:
			ctx.Send("[6]影响不好啦!")
		case 7:
			ctx.Send("[7]太涩啦，🐛了!")
		case 8:
			ctx.Send("[8]已经🐛不动啦...")
		}
		if dhash != "" && !noimg {
			b14, err3 := url.QueryUnescape(dhash)
			if err3 == nil {
				ctx.Send("给你点提示哦：" + b14)
			}
		}
	} else {
		var last_message_id int64
		if !noimg {
			last_message_id = ctx.SendChain(message.Image(CACHE_URI + strconv.FormatInt(lastvisit, 10)))
		} else {
			last_message_id = ctx.Event.MessageID
		}
		switch class {
		case 0:
			ctx.SendChain(message.Reply(last_message_id), message.Text("[0]这啥啊"))
		case 1:
			ctx.SendChain(message.Reply(last_message_id), message.Text("[1]普通欸"))
		case 2:
			ctx.SendChain(message.Reply(last_message_id), message.Text("[2]有点可爱"))
		case 3:
			ctx.SendChain(message.Reply(last_message_id), message.Text("[3]不错哦"))
		case 4:
			ctx.SendChain(message.Reply(last_message_id), message.Text("[4]很棒"))
		case 5:
			ctx.SendChain(message.Reply(last_message_id), message.Text("[5]我好啦!"))
		}
	}
}
