
## 生成api代码
`goctl api go --api demo.api --dir .  --style go_zero`

## 安装goctl-swagger
`https://github.com/zeromicro/goctl-swagger`

## 生成swagger文档
`goctl api plugin -plugin goctl-swagger="swagger -filename demo.json" -api demo.api -dir ./etc`
