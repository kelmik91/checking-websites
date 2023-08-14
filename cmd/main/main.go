package main

import (
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"yandex/internal/db"
	"yandex/internal/logger"
	"yandex/internal/site"
)

func main() {
	createLockFileOrDie()
	defer func() {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		err = os.Remove(exPath + "/mainCheck.lock")
		if err != nil {
			return
		}
	}()
	start := time.Now()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	wg := sync.WaitGroup{}
	hosts := db.GetHosts()

	for i := range hosts {
		wg.Add(1)
		go site.CheckSite(hosts[i], &wg)
	}

	wg.Wait()

	duration := time.Since(start)
	log.Println(duration)
	endWork(duration)
}

func endWork(duration time.Duration) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	file, err := os.OpenFile(exPath+"/logGoWork.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		//logger.WriteWork(err.Error())
		fmt.Println(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	loc, _ := time.LoadLocation("Europe/Moscow")
	_, err = writer.WriteString(time.Now().In(loc).Format(time.RFC822) + " Время выполнения " + duration.String() + "\n")
	if err != nil {
		//logger.WriteWork(err.Error())
		fmt.Println(err)
		return
	}
}

func createLockFileOrDie() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	_, errFile := os.Stat(exPath + "/mainCheck.lock")
	if os.IsNotExist(errFile) {
		file, err := os.Create(exPath + "/mainCheck.lock")
		if err != nil {
			logger.WriteWork(err.Error())
			log.Println(err)
		}
		defer file.Close()
	} else {
		panic("Long run")
	}
}
