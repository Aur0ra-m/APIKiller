package logger

import (
	log "github.com/sirupsen/logrus"
)

func Errorln(args ...interface{}) {
	//fmt.Println(args)

	log.Errorln(args)
}

func Infoln(args ...interface{}) {
	//fmt.Println(args)

	log.Infoln(args)
}

func Debugln(args ...interface{}) {
	//fmt.Println(args)

	log.Debugln(args)
}
