package main

import (
	"./config"
	"./logger"
	"github.com/gorilla/websocket"
	wave "github.com/zenwerk/go-wave"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func checkFileIsExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	defer c.Close()
	if err != nil {
		logger.Error.Println("Error upgrade request: " + err.Error())
		return
	}

	filename := "./test.wav"
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

}

func main() {
	logger.Info.Println("Service start.")
	log.SetFlags(0)

	if !checkFileIsExist(config.G_wav_dir) {
		err := os.MkdirAll(config.G_wav_dir, 0666)
		if err != nil {
			logger.Error.Println("wav mkdir failed :" + err.Error())
			return
		}
	}

	http.HandleFunc("/ws", echo)

	logger.Info.Println("Service running on " + config.G_host)
	log.Fatal(http.ListenAndServe(config.G_host, nil))
}
