package main

import (
	"bufio"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/spf13/viper"
	"github.com/suntt2019/EduOJJudger/base"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func readConfig() {
	// TODO: set default
	log.Debug("Reading config.")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.WithField("error", err).Fatal("Could not read config.")
	}
	base.ScriptPath = viper.GetString("path.scripts")
	if base.ScriptPath[len(base.ScriptPath)-1] == '/' {
		base.ScriptPath = base.ScriptPath[:len(base.ScriptPath)-1]
	}
	base.RunPath = viper.GetString("path.runs")
	if base.RunPath[len(base.RunPath)-1] == '/' {
		base.RunPath = base.RunPath[:len(base.RunPath)-1]
	}
}

func initConsoleLogger() {
	log.Debug("Initializing console logger.")
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)

	log.SetFormatter(&log.TextFormatter{
		PadLevelText:    true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05 MST",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			function = f.Function + "()]"
			file = f.File
			if strings.HasPrefix(file, prefixPath) {
				file = file[len(prefixPath):]
			}
			file = fmt.Sprintf(" [.%s:%d", file, f.Line)
			return
		},
	})
	log.SetReportCaller(true)
}

func initFileLogger() {
	log.Debug("Initializing file logger.")
	filePath := viper.GetString("log.log_path")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.WithField("error", err).Error("Failed to open log file")
	}
	defer file.Close()
	// TODO: fix log file writing

	log.AddHook(&writer.Hook{
		Writer: bufio.NewWriter(file),
		LogLevels: []log.Level{
			log.ErrorLevel,
			log.FatalLevel,
			log.PanicLevel,
		},
	})
}

func initHttpClient() {
	log.Debug("Initializing http client.")
	base.HC = resty.New()
	base.HC.
		SetHeader("Authorization", viper.GetString("auth.token")).
		SetHeader("Judger-Name", viper.GetString("auth.name")).
		SetHostURL(fmt.Sprintf("http://%s:%d/judger", viper.GetString("backend.host"), viper.GetInt("backend.port")))
}
