package main

import (
	"os"

	"br.com.mlabs/api"
	"br.com.mlabs/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02/01/2006.15:04:05",
	})
	logrus.SetOutput(os.Stdout)
	storage.Connect()
	api.Start()
}
