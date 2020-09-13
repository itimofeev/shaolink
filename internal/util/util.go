package util

import "log"

func CheckErr(err error, msg ...string) {
	if err != nil {
		log.Panic(err, msg)
	}
}

func CheckOK(ok bool, msg ...string) {
	if !ok {
		log.Panic(msg)
	}
}
