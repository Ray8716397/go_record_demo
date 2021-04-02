package main

import (
	"./config"
	"./logger"
	"fmt"
	"github.com/gorilla/websocket"
	wave "github.com/zenwerk/go-wave"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func checkFileIsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	defer c.Close()
	if err != nil {
		logger.Error.Println("Error upgrade request: " + err.Error())
		return
	}

	var uid string
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params) == 0 {
		logger.Error.Println("no userid")
		return
	} else {
		uid = params["uid"][0]
		logger.Info.Println(uid + " ws connected")
	}

	udirPath := path.Join(config.G_wav_dir, uid)
	if !checkFileIsExist(udirPath) {
		err := os.MkdirAll(udirPath, os.ModePerm)
		if err != nil {
			logger.Error.Println("udir mkdir failed :" + err.Error())
			return
		}
	}
	filename := path.Join(udirPath, time.Now().Format("2006-01-02_15:04:05")+".wav")
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		logger.Error.Println("Cannot create file: " + err.Error())
		return
	}
	param := wave.WriterParam{
		Out:           f,
		Channel:       1,
		SampleRate:    16000,
		BitsPerSample: 16,
	}

	w1, err2 := wave.NewWriter(param)
	defer w1.Close()

	if err2 != nil {
		panic(err2)
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.Error.Println("websocket read: " + err.Error())
			break
		}

		_, err = w1.WriteSample8(message)
		if err != nil {
			logger.Error.Println("websocket recv: " + err.Error())
			break
		}

		//log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			logger.Error.Println("websocket write: " + err.Error())
			break
		}
	}

	absFp, err := filepath.Abs(filename)
	dbFp := "/home/ray/Workspace/project/asr/src/dictation-kit/test.db"
	if checkFileIsExist(dbFp) { //如果文件存在
		f, err = os.OpenFile(dbFp, os.O_WRONLY, os.ModePerm) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err = os.Create(dbFp) //创建文件
		fmt.Println("文件不存在")
	}
	if err != nil {
		logger.Error.Println("open file error: " + err.Error())
		return
	}
	n3, err := f.WriteString(absFp) //写入文件(字节数组)
	fmt.Printf("写入 %d 个字节n", n3)
	f.Sync()
	if err != nil {
		logger.Error.Println("write file error: " + err.Error())
		return
	}
	exec.Command("/bin/bash", "-c", `cd /home/ray/Workspace/project/asr/src/dictation-kit && ./run-linux-dnn.sh`).Run()
}

func main() {
	logger.Info.Println("Service start.")
	log.SetFlags(0)

	if !checkFileIsExist(config.G_wav_dir) {
		err := os.MkdirAll(config.G_wav_dir, os.ModePerm)
		if err != nil {
			logger.Error.Println("wav mkdir failed :" + err.Error())
			return
		}
	}

	http.HandleFunc("/ws", echo)

	logger.Info.Println("Service running on " + config.G_host)
	log.Fatal(http.ListenAndServe(config.G_host, nil))
}
