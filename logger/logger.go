package logger

import (
	"encoding/json"
	"fmt"
	"log"
)

func Output(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Println(msg)
}

// Info writes info log.
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Output("[INFO]%s", msg)
}

// Debug writes debug log.
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Output("[DEBUG]%s", msg)
}

// Error writes error log.
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Output("[ERROR]%s", msg)
}

func SetStdFlags() {
	log.SetFlags(log.LstdFlags)
}

func SetDebugFlags() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// JSON marshals the value for logging.
func JSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
