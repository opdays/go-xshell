package xshell

import (
	"golang.org/x/crypto/ssh"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
	"unicode/utf8"
	"bufio"
	"encoding/json"
	"log"
)

type Terminal struct {
	Columns uint32 `json:"cols"`
	Rows    uint32 `json:"rows"`
}
type SSH struct {
	Ip       string
	Port     int
	Username string
	Password string
	Client   *ssh.Client
	Session  *ssh.Session
	channel  ssh.Channel
}

func (this *SSH) Connect() (error) {
	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = this.Username
	config.Auth = []ssh.AuthMethod{ssh.Password(this.Password)}
	config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	config.Timeout = time.Second * 10
	addr := fmt.Sprintf("%s:%d", this.Ip, this.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	} else {
		this.Client = client
		return nil
	}
}

func (this *SSH) Exec(cmd string) (stdout string, exitCode int, err error) {
	//var buf bytes.Buffer
	session, err := this.Client.NewSession()
	if err != nil {
	}

	defer this.Client.Close()
	defer session.Close()
	//session.Stdout = &buf
	//session.Stderr = &buf
	//err = session.Run(cmd)
	this.Session = session
	//https://github.com/golang/go/issues/14251
	bs, err := session.CombinedOutput(cmd)
	if err != nil {
		if v, ok := err.(*ssh.ExitError); ok {
			exitCode = v.Waitmsg.ExitStatus()
		}
		log.Println(err)
	}
	stdout = string(bs)
	return stdout, exitCode, err
}

func (this *SSH) ExecNotClose(cmd string) (stdout string, exitCode int, err error) {
	//var buf bytes.Buffer
	session, err := this.Client.NewSession()
	if err != nil {
	}

	//session.Stdout = &buf
	//session.Stderr = &buf
	//err = session.Run(cmd)
	this.Session = session
	//https://github.com/golang/go/issues/14251
	bs, err := session.CombinedOutput(cmd)
	if err != nil {
		if v, ok := err.(*ssh.ExitError); ok {
			exitCode = v.Waitmsg.ExitStatus()
		}
		log.Println(err)
	}
	stdout = string(bs)
	return stdout, exitCode, err
}
func (this *SSH)Close()  {
	if this.Client != nil{
		this.Client.Close()
	}
	if this.Session != nil{
		this.Session.Close()
	}

}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

func (this *SSH) RequestTerminal(terminal Terminal) (*SSH) {

	session, err := this.Client.NewSession()
	this.Session = session
	if err != nil {
		log.Println(err)
	}

	channel, incomingRequests, err := this.Client.OpenChannel("session", nil)
	this.channel = channel
	if err != nil {
		return nil
	}
	go func() {
		for req := range incomingRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  terminal.Columns,
		Rows:     terminal.Rows,
		Width:    uint32(terminal.Columns * 8),
		Height:   uint32(terminal.Columns * 8),
		Modelist: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	return this
}


func (this *SSH) Handle(ws *websocket.Conn) {
	defer func() {
		//有可能panic 导致程序退出 在这里捕获下
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	//modes := ssh.TerminalModes{
	//	ssh.ECHO:          1, // enable echoing
	//	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	//	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	//}
	//err = this.Session.RequestPty("xterm-256color", 80, 24, modes)
	//go func() {
	//	//channel 读出来的写入ws
	//	ws.WriteMessage(websocket.TextMessage, []byte("hello "))
	//	for{
	//		br := bufio.NewReader(channel)
	//		r,size,_:=br.ReadRune()
	//		if size > 0{
	//			global.GLog.Debug(string(r))
	//			p := make([]byte, utf8.RuneLen(r))
	//			utf8.EncodeRune(p, r)
	//			ws.WriteMessage(websocket.TextMessage, []byte(p))
	//		}
	//
	//	}
	//
	//}()
	//go func() {
	//	for {
	//		t,message,_:=ws.ReadMessage()
	//		global.GLog.Debug("message",message)
	//		channel.Write(message)
	//		time.Sleep(1*time.Second)
	//		var buf []byte
	//		channel.Read(buf)
	//		global.GLog.Debug("buf",buf)
	//		ws.WriteMessage(t,buf)
	//	}
	//}()

	go func() {
		//第一个协程获取ws用户输入 写入shell =this.channel
		defer func() {
			//有可能panic 导致程序退出 在这里捕获下
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		for {
			m, p, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if m == websocket.TextMessage {
				resize := &Terminal{}
				err = json.Unmarshal(p, resize)
				if err == nil {
					log.Println("resize", resize)
					this.Session.Close()
					//判断重置request terminal大小
					this.RequestTerminal(*resize).Handle(ws)
				} else {
					//global.GLog.Debug("resize",err)
					if _, err := this.channel.Write(p); nil != err {
						return
					}
				}

			}
		}
	}()
	//模仿https://github.com/shibingli/webconsole.git
	go func() {
		//第2个协程
		// 1  获取shell =this.channel输出
		// 2  将shell =this.channel输出 写入ws
		// time.Microsecond * 100 是切换间隔时间
		defer this.Client.Close()
		defer this.Session.Close()
		defer func() {
			//有可能panic 导致程序退出 在这里捕获下
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		br := bufio.NewReader(this.channel)
		buf := []byte{}
		t := time.NewTimer(time.Microsecond * 100)
		defer t.Stop()
		r := make(chan rune)

		go func() {
			defer this.Client.Close()
			defer this.Session.Close()
			defer func() {
				//有可能panic 导致程序退出 在这里捕获下
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			for {
				//if err:=ws.PingHandler()("ping");err!=nil{
				//	global.GLog.Error("ping error",err)
				//	return
				//}
				//if err:=ws.WriteMessage(websocket.PingMessage,[]byte("ping"));err!=nil{
				//	//心跳
				//	global.GLog.Error("心跳",err)
				//	ws.Close()
				//}
				x, size, err := br.ReadRune()
				if err != nil {
					//客户端exit  ctrl+d会导致EOF错误
					log.Println(err)
					ws.WriteMessage(1, []byte("\033[31m已经关闭连接!\033[0m"))
					ws.Close()
					//wsclose 之后会到229 会触发defer this.Client.Close()
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		for {
			select {
			case <-t.C:
				if len(buf) != 0 {
					err := ws.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						//ws客户端tcp连接指针关闭会导致use of closed network connection
						log.Println(err)
						return
					}
				}
				t.Reset(time.Microsecond * 100)
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					log.Println("d", d)
					buf = append(buf, []byte("@")...)
				}
			}
		}

	}()

	//if err != nil {
	//	global.GLog.Error(err)
	//}
	//ok,err := this.Session.SendRequest("shell", true, nil)
	//if !ok ||err != nil {
	//	global.GLog.Error(err)
	//}

}
