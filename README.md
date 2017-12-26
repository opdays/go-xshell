## go 实现的web xshell

[github地址](https://github.com/yy584089958/go-xshell)

### 依赖

- [gin](https://github.com/gin-gonic/gin)
- [websocket](https://github.com/gorilla/websocket)
- [xterm.js](https://github.com/xtermjs/xterm.js)
- [crypto](https://github.com/golang/crypto)


### 快速使用

下载我编译好的二进制

- [linux64](https://github.com/yy584089958/go-xshell/files/1586776/xshell_v0.0.1_linux-amd64.tar.gz)
- [windows64](https://github.com/yy584089958/go-xshell/files/1586785/xshell_v0.0.1_windows-amd64.zip)

下载之后解压

windows用户打开cmd `./go-xshell.exe` 启动

linux用户打开terminal `./go-xshell` 启动

### 编译安装
```bash

git clone http://github.com/yy584089958/go-xshell

cd go-xshell

govendor sync

go build

./go-xshell

```
打开浏览器 http://127.0.0.1:18080/console/ip=x.x.x.x



登录截图


![登录](https://github.com/yy584089958/go-xshell/raw/master/screenshot/console-1.png)

登录成功截图


![登录成功](https://github.com/yy584089958/go-xshell/raw/master/screenshot/console-sucs.png)


vim截图



