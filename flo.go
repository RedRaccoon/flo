package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	//mutex for creating *.ts files
	mutex            = &sync.Mutex{}
	nbrGoRoutine     = 0
	nbrGoRoutineDone = 0
	fileName         = "1"
)

const (
	//max size of a ts file name
	tsFileMaxSize = 20
	nameTempFile  = "_tempTsFiles"
)

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(nameTempFile)
	if err != nil {
		//check if the folder is deleted to be shure having an empty directory

	}

	//create the directory for the *.ts files, 1 = full right
	os.Mkdir(nameTempFile, 1)
	if err != nil {
		fmt.Println(err)
	}

	destination := filepath.Dir(ex) + "\\" + nameTempFile + "\\"

	if len(os.Args) < 2 {
		log.Fatal("need the url of the playlist")
	}

	var url = os.Args[1]

	var toDl = getChunkList(url)

	//channel for the files to download
	done := make(chan bool)
	pDone := &done

	startTime := time.Now()

	resp, err := http.Get(strings.Replace(url, "playlist.m3u8", toDl, -1))
	if err != nil {
		fmt.Println(("error get the playlist => "), err)
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		line := s.Text()
		fmt.Println(line)

		if strings.HasSuffix(line, ".ts") {
			go downloadFile(destination, fileName+".ts", pDone, url+line)
			fileName = binaryorder(fileName)
			nbrGoRoutine++
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("There is " + strconv.Itoa(nbrGoRoutine) + " goroutines downloading")
	for i := 0; i < nbrGoRoutine; i++ {
		<-done
		nbrGoRoutineDone = debug(nbrGoRoutineDone)
	}

	t := time.Now()
	elapsed := t.Sub(startTime)

	fmt.Println(elapsed)

	c := exec.Command("cmd", "/C copy /b _tempTsFiles\\*.ts final.ts")

	//out & stdeer help give an error more usefull for debuging
	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
	fmt.Println("File available !  " + out.String())

}

func downloadFile(filepath string, suffix string, done *chan bool, url string) error {

	// Create the file
	mutex.Lock()
	out, err := os.Create(filepath + suffix)
	mutex.Unlock()
	if err != nil {
		fmt.Println("error create file +> ", err)

		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(("error get the data => "), err)
		return err
	}

	go writeOnDisk(out, resp.Body, err, done)

	return nil
}

func writeOnDisk(out *os.File, respBody io.ReadCloser, err error, done *chan bool) {

	_, err = io.Copy(out, respBody)
	if err != nil {
		fmt.Println("error write the body file => ", err)

	}
	out.Close()
	respBody.Close()
	*done <- true

}

//by default windows make a binary order with the files
//the order of 1 - 2 - 11 - 21 is 1 - 11 - 2 - 21
func binaryorder(s string) string {

	if len(s) > tsFileMaxSize {
		sInt := s[0:1]

		var newInt, _ = strconv.Atoi(sInt)
		newInt++
		return strconv.Itoa(newInt)

	}

	return s + s[0:1]

}

//get the const for flo_test
func getTsFileMaxSize() int {
	return tsFileMaxSize
}

func debug(i int) int {

	i++
	fmt.Println("routine finished " + strconv.Itoa(i))
	return (i)

}

func getChunkList(url string) string {
	resp, err := http.Get(url)
	nextIsChunkList := false
	if err != nil {
		fmt.Println(("error get the playlist => "), err)
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		line := s.Text()
		if nextIsChunkList {
			return (line)
		}
		//we obviously choose the best resolution ...
		if strings.HasSuffix(line, "RESOLUTION=1280x720") {
			nextIsChunkList = true

		}

	}
	return "nothing"
}
