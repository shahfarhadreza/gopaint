package reza

import "log"

func logError(msg string) {
	log.Printf("error: %s\n", msg)
}

func logInfo(msg string) {
	//log.Println(msg)
}

func init() {
	logInfo("initialize reza library!")
	loadApiFunctions()
}
