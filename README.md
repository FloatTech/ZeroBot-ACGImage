> 本插件已合并入[ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin)，本仓库不再迭代/维护
# ~~ZeroBot-ACGImage~~
~~基于[Zerobot-ACGImage-Classify](https://github.com/FloatTech/Zerobot-ACGImage-Classify)编写的二次元图片随机、鉴别插件。~~

# 使用说明

```go
import _ "github.com/FloatTech/ZeroBot-ACGImage
```
> 本插件检测如下口令

### 评价图片
发送一张图片并使用AI进行评分。

### 设置随机图片网址[url]
设置随机图片/直接随机调用的API网址。

### 随机图片
从随机图片API随机一张二次元图片，并回复其评价。默认从`loliconapi`随机且开启r18模式，但是r18图片不会发出而是以提示给出。若想获取图片，可访问sayuri.fumiama.top:8080/img?path=提示的汉字。

本口令有5秒CD。

### 直接随机
仅管理可用，从设置的随机图片网址获取一张图片，不做r18检测直接发送，后果自负。如果出r18图可尽快通过发送"太涩了"撤回。

本口令有5秒CD。

### 太涩了
手动撤回最近发的一张图，无法撤回更多。
