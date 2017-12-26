# go 实现的web xshell

依赖
- [gin](https://github.com/gin-gonic/gin)
- [websocket](https://github.com/gorilla/websocket)
- [xterm.js](https://github.com/xtermjs/xterm.js)
- [crypto](https://github.com/golang/crypto)

安装
```bash

git clone http://github.com/yy584089958/go-xshell

cd go-xshell

govendor sync

go build

./xshell

```
打开浏览器 http://127.0.0.1:18080/console/ip=x.x.x.x



登录截图


![登录](https://github.com/yy584089958/go-xshell/raw/master/screenshot/console-1.png)

登录成功截图


![登录成功](https://github.com/yy584089958/go-xshell/raw/master/screenshot/console-sucs.png)


vim截图


![vim](https://github.com/yy584089958/go-xshell/raw/master/screenshot/console-vim.png)