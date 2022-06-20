package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Chandler-WQ/go_exercise_demo/pkg/log"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

var sourcePath string
var targetPath string

// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o txt2xlsx.exe ./excel
//go build -o txt2xlsx ./excel
//./txt2xlsx -source=./excel target=data.xlsx
func main() {
	flag.StringVar(&sourcePath, "source", "./", "the source dir path")
	flag.StringVar(&targetPath, "target", "data.xlsx", "the target file path")

	helpFlag := flag.Bool("help", false, "Show usage info and exit.")
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	f := excelize.NewFile()

	// Create a new sheet.

	var files []string
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "txt") {
			files = append(files, path)
			return nil
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	var i = 0
	for _, file := range files {
		sheetName := "sheet" + cast.ToString(i)
		index := f.NewSheet(sheetName)
		wg := &sync.WaitGroup{}
		c := make(chan string, 100)
		wg.Add(1)
		go read(file, c)

		i++
		var j = 0
		for {
			str, open := <-c
			if !open {
				wg.Done()
				break
			}
			nums, err := conv(str)
			if err != nil {
				log.Errorf("conv %s", err)
				continue
			}
			for i := 0; i < len(nums); i++ {
				row := 'A' + i
				posCell := string(rune(row)) + cast.ToString(j)
				f.SetCellValue(sheetName, posCell, nums[i])
			}
			j++
		}
		f.SetActiveSheet(index)
		wg.Wait()
	}

	if err := f.SaveAs(targetPath); err != nil {
		log.Errorf("SaveAs %s", err)
	}
}

func conv(str string) ([]float64, error) {
	vals := strings.Split(str, " ")
	if len(vals) == 0 {
		return nil, errors.New("strs is null")
	}

	nums := make([]float64, 0, len(vals))
	for _, val := range vals {
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		num, err := cast.ToFloat64E(val)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func read(filePath string, c chan string) error {
	defer close(c)
	f, err := os.OpenFile(filePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Errorf("open file err %v", err)
		return err
	}
	defer f.Close()
	//每次读10MB,保证可以读取一个文件
	r := bufio.NewReaderSize(f, 1024*1024*10)
	for {
		lineStr, err := ReadLine(r)
		if err != nil {
			if err == io.EOF {
				log.Infof("readLine end,it is eof:%v", err)
				break
			}
			log.Errorf("readLine err :%v", err)
			return err
		}
		log.Infof("read %s", lineStr)
		if lineStr != "" {
			c <- lineStr
		}
	}
	return nil
}

func ReadLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}
