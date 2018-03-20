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
	mutex = &sync.Mutex{}
)

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll("_tempTsFiles")
	if err != nil {
		//check if the folder is really deleted
	}

	os.Mkdir("_tempTsFiles", 1)
	if err != nil {
		fmt.Println(err)
	}

	destination := filepath.Dir(ex) + "\\_tempTsFiles\\"

	if len(os.Args) < 2 {
		log.Fatal("need the url of the playlist")
	}

	var url = os.Args[1]

	//var toDl="chunklist_w135572285_b3300956.m3u8"
	var toDl = getChunkList(url)
	done := make(chan bool)
	pDone := &done
	var nbrGoRoutine = 0
	var nbrGoRoutineDone = 0
	var fileName = "1"

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

			fmt.Println("on telecharge   " + url + line)
			go downloadFile(destination, fileName+".ts", pDone, url+line)
			fileName = binaryorder(fileName)
			nbrGoRoutine++
		}
		if line == "#EXT-X-ENDLIST" {
			//trigger
			fmt.Println("trouvé")
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("il y a " + strconv.Itoa(nbrGoRoutine) + " en cours")
	for i := 0; i < nbrGoRoutine; i++ {
		<-done
		nbrGoRoutineDone = debug(nbrGoRoutineDone)
	}

	t := time.Now()
	elapsed := t.Sub(startTime)

	fmt.Println(elapsed)

	c := exec.Command("cmd", "/C copy /b _tempTsFiles\\*.ts final.ts")
	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
	fmt.Println("Result: " + out.String())

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

		//return err
	}
	out.Close()
	respBody.Close()
	*done <- true

}

func binaryorder(s string) string {

	//taille max des nom du fichier => len(s)> x
	if len(s) > 20 {
		sInt := s[0:1]

		var newInt, _ = strconv.Atoi(sInt)
		newInt++
		return strconv.Itoa(newInt)

	}

	return s + s[0:1]

}

func debug(i int) int {

	i++
	fmt.Println("routine terminé " + strconv.Itoa(i))
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
		if strings.HasSuffix(line, "RESOLUTION=1280x720") {
			nextIsChunkList = true

		}

	}
	return "nothing"
}
