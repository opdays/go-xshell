package console

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/yy584089958/go-xshell/xshell"
	"github.com/yy584089958/go-xshell/util"
	"github.com/gorilla/websocket"
	"time"
	"encoding/base64"
	"strings"
	"math"
	"github.com/pkg/errors"
	"log"
)

func ConsoleHtml(c *gin.Context) {
	c.HTML(200, "console.html", nil);
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ConsoleWs(c *gin.Context) {
	key := c.Query("key")
	s,err:=checkConsoleKey(key)
	if err!=nil{
		return
	}
	cols := c.DefaultQuery("cols", "150")
	rows := c.DefaultQuery("rows", "35")
	col, _ := strconv.Atoi(cols)
	row, _ := strconv.Atoi(rows)
	log.Println(cols)
	log.Println(rows)
	terminal := xshell.Terminal{
		Columns: uint32(col),
		Rows:    uint32(row),
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}

	//for {
	//    t, msg, err := conn.ReadMessage()
	////client,err=client.Connect()
	////out,err:=client.Exec(string(msg))
	//    if err != nil {
	//        break
	//    }
	//global.GLog.Debug(msg)
	//    conn.WriteMessage(t, []byte("hello world"))
	//}
	err = s.Connect()
	if err != nil {
		log.Println("connect failure", err)
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
		return
	}
	s.RequestTerminal(terminal)
	s.Handle(conn)
}

func SimpleEncrypt(text string) (result string){
	resultBs, err := util.AesEncrypt([]byte(text), []byte(util.AesKey))
	if err != nil {
		log.Println("encrypt",err)
		return
	}
	result = base64.StdEncoding.EncodeToString(resultBs)
	return
}

func checkConsoleKey(key string) (client *xshell.SSH,err error){
	//ip::user::pass[x]timestamp
	result := SimpleDecrypt(key)
	log.Println(result)
	sliceStr := strings.Split(result,"[x]")
	now := time.Now().Unix()
	timeStamp := sliceStr[1]
	timeStampInt,_ :=  strconv.Atoi(timeStamp)
	if (math.Abs(float64(now - int64(timeStampInt))) < 10){
		//10秒过期
		sliceConfig:=strings.Split(sliceStr[0],"::")
		client = &xshell.SSH{
			Ip:sliceConfig[0],
			Username:sliceConfig[1],
			Password:sliceConfig[2],
			Port:22,
		}
		log.Println(sliceConfig)
		return client,nil
	}
	return nil,errors.New("超时")
}
func SimpleDecrypt(text string) (result string) {
	resultBs, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Println("base64.StdEncoding.DecodeString",err)
		return
	}
	resultBs, err = util.AesDecrypt([]byte(resultBs), []byte(util.AesKey))
	if err != nil {
		log.Println("util.AesDecrypt",err)
		return
	}
	result = string(resultBs)
	return
}


func AesEncrypt(c *gin.Context) {
	//加密接口 带上时间戳加密
	text := c.PostForm("text")
	timeStamp := time.Now().Unix()
	//int64到string
	//string:=strconv.FormatInt(int64,10)
	timeStampString :=strconv.FormatInt(timeStamp,10)
	text = text + "[x]" + timeStampString
	resultStr := SimpleEncrypt(text)
	c.String(200, resultStr)
}

func AesDecrypt(c *gin.Context) {
	text := c.PostForm("text")
	resultStr := SimpleDecrypt(text)
	c.String(200,resultStr)
}