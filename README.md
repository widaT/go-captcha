# go captcha

受[AJ-Captcha](https://gitee.com/anji-plus/captcha)启发，用golang实现的滑动验证码。


## 运行

```bash
$ go mod tidy
$ go run example/puzzle_captcha/main.go

```

## 前端

[前端代码仓库](https://gitee.com/wida/gocatcha-ui)

运行前端

```bash
$ git clone git@gitee.com:wida/gocatcha-ui.git
$ cd gocatcha-ui
$ npm i
$ npm run dev
```

## 运行截图

![](./doc/1.gif)


## 感谢

本程序算法受[AJ-Captcha](https://gitee.com/anji-plus/captcha)启发，本程序的前端UI则来着[AJ-Captcha](https://gitee.com/anji-plus/captcha)的vue代码，感谢`AJ-Captcha`的开源。