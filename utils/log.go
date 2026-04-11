package utils

import "fmt"

func Log(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

func LogErr(err error) {
	fmt.Printf("ERROR: " + err.Error() + "\n")
}

func LogErrM(msg string, args ...interface{}) {
	fmt.Printf("ERROR: "+msg+"\n", args...)
}
