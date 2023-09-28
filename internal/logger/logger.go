package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func WriteWork(message string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
		//TODO заменить панику на логирование
	}
	exPath := filepath.Dir(ex)
	file, err := os.OpenFile(exPath+"/logGoWork.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	loc, _ := time.LoadLocation("Europe/Moscow")
	_, err = writer.WriteString(time.Now().In(loc).Format(time.RFC822) + " " + message + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func WriteWorkTelegram(message string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
		//TODO заменить панику на логирование
	}
	exPath := filepath.Dir(ex)
	file, err := os.OpenFile(exPath+"/logGoWorkTelegram.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	loc, _ := time.LoadLocation("Europe/Moscow")
	_, err = writer.WriteString(time.Now().In(loc).Format(time.RFC822) + " " + message + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}
