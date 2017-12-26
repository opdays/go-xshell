package xshell

//import (
//	"github.com/pkg/sftp"
//	"os"
//	"opsgit.feidee.org/yangyang/gops/global"
//	"io"
//)
//
//func (this *SSH) Sftp(localFile string, remoteFile string) (err error) {
//	sftpClient, err := sftp.NewClient(this.Client)
//	defer func() {
//		if errI:=recover();errI!=nil{
//			if v,ok:=errI.(error);ok{
//				err = v
//				return
//			}
//
//		}
//	}()
//	if err != nil {
//		global.GLog.Debug("NewClient", err)
//		return
//	}
//	_, err = sftpClient.Lstat(remoteFile)
//	//https://godoc.org/github.com/pkg/sftp
//	if err != nil {
//		//判断远程文件是否存在 不存在进入体内
//	}else{
//		global.GLog.Warning(remoteFile," is already exist")
//	}
//	remoteFiler, err := sftpClient.Create(remoteFile)
//	defer remoteFiler.Close()
//	if err != nil {
//		global.GLog.Debug("sftpClient.Create", remoteFile, err)
//		return
//	}
//	localFiler, err := os.Open(localFile)
//	if err != nil {
//		global.GLog.Debug("os.Open", localFiler, err)
//		return
//	}
//	io.Copy(remoteFiler, localFiler)
//	return
//}
