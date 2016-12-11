### 今日头条图片爬虫

根据指定标签爬取图片，以文章名为目录存储。

### RUN

```
$ git clone git@github.com/qichengzx/toutiaoSpider.git
$ cd toutiaoSpider
$ //main.go后添加需要爬取的标签名
$ go run main.go 街拍 摄影
```

### TODO 

> 并发爬取
> 以 标签名/文章名/文件名 结构存储
> 错误处理

### 已知问题

某些情况下会出现 unexpected EOF 错误导致退出