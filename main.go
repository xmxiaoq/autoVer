package main

import (
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"os"
	"github.com/ungerik/go-dry"
	_"time"
	"github.com/rjeczalik/notify"
	"os/signal"
	"fmt"
)

var FS = afero.Afero{Fs: afero.NewOsFs()}

func init() {
	_, _ = FS.Exists(os.Args[0] + ".log")
	_ = cast.ToString(100)
	_ = dry.SyncMap{}
}

type VerInfo struct {
	code_url   string `json:"code_url"`
	update_url string `json:"update_url"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	//sugar.Infow("Failed to fetch URL.",
	//	// Structured context as loosely-typed key-value pairs.
	//	"url", "url",
	//	"attempt", 3,
	//	"backoff", time.Second,
	//)
	//sugar.Infof("Failed to fetch URL: %s", "url")

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening on events within current working directory.
	// Dispatch each create and remove events separately to c.
	rootPath := `D:\fanfan\mj_h5\bin-release\native`
	jsonPath := `D:\fanfan\mj_h5\bin-release\native\version.json`
	if err := notify.Watch(rootPath, c, notify.Create, notify.Remove); err != nil {
		sugar.Fatal(err)
	}
	defer notify.Stop(c)

	go func() {
		for evt := range c {
			if evt.Event() == notify.Create {
				pth := evt.Path()
				st, err := os.Stat(pth)
				if err == nil {
					if st.IsDir() {
						str := `{
"code_url": "http://10.0.0.35/%s/game_code_%s.zip",
"update_url": "http://10.0.0.35/%s"
}`
						dirName := st.Name()
						//var info VerInfo
						//info.code_url = fmt.Sprintf("http://10.0.0.35/%s/game_code_%s.zip", dirName, dirName)
						//info.update_url = fmt.Sprintf("http://10.0.0.35/%s", dirName)
						//dry.FileSetJSONIndent(jsonPath,info, "    ")

						jsonStr := fmt.Sprintf(str, dirName, dirName, dirName)
						sugar.Info(jsonStr)
						dry.FileSetString(jsonPath, jsonStr)
					}
				}
			}
		}
	}()

	//ei := <-c
	//sugar.Info("Got event:", ei)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	sugar.Info("exit")
}
