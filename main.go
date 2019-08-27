package main

import (
	"bufio"
	"fmt"
	"log"
	"ops/src/entegor"
	"ops/src/filenum"
	"ops/src/sysutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	err := os.Setenv("LANG", "en_US")
	if err != nil {
		log.Println(err)
	}
	//INIFile := os.Args[1]
	//HostIP12 := os.Args[2]
	//HostIP := os.Args[3]
	INIFile := "./inifile.ini"
	HostIP12 := "011111111111"
	HostIP := "11.111.111.111"
	WorkDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_, fullScriptName := filepath.Split(os.Args[0])
	scriptName := strings.Split(fullScriptName, ".")[0]
	TmpDir := filepath.Join(WorkDir, "temp")
	OutDir := filepath.Join(WorkDir, "out")
	LogDir := filepath.Join(WorkDir, "log")
	Dirs := []string{TmpDir, OutDir, LogDir}
	for _, dir := range Dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}
	}
	LogFileName := scriptName + HostIP12 + ".log"
	OutTmpFileName := scriptName + HostIP12 + ".out"
	OutFileName := "check" + HostIP12 + ".out"
	LogFile := filepath.Join(LogDir, LogFileName)
	OutTmpFile := filepath.Join(TmpDir, OutTmpFileName)
	OutFile := filepath.Join(OutDir, OutFileName)
	fmt.Println(LogFile, OutTmpFile, OutFile)
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hostname)
	now := time.Now().Format(entegor.LongForm)
	ini, err := os.Open(INIFile)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	iniScanner := bufio.NewScanner(ini)
	for iniScanner.Scan() {
		cfgItem := iniScanner.Text()
		walkDir := strings.Split(cfgItem, "|")[0]
		filenum.ThordHold, _ = strconv.Atoi(strings.Split(strings.Split(cfgItem, "#")[0], "|")[1])
		fmt.Println(filenum.ThordHold)
		filepath.Walk(walkDir, filenum.CheckNum)
		fmt.Println(filenum.Files)
		var Data float64
		if filenum.Files == nil {
			Data = float64(0)
		} else {
			Data = float64(1)
		}
		stCode := entegor.GetStCode(Data, cfgItem)
		descMsg := walkDir
		good := entegor.GetGood(cfgItem)
		stCodeString := strconv.Itoa(stCode)
		DataString := strconv.FormatFloat(Data, 'E', 1, 64)
		warnMsg := ""
		var result string
		if stCode == 0 {
			head := entegor.GetHead(cfgItem)
			result = head + "=" + stCodeString + "|" + now + "|" + DataString + "|" + good + "|" + descMsg + "\n"
		} else {
			head := entegor.GetWarningHead(cfgItem)
			result = head + "=" + stCodeString + "|" + now + "|" + "AOMS" + "|" + fullScriptName + "|" + filenum.ErrCode + "|" + hostname + "|" + HostIP + "|" + "" + "|" + "" + "|" + warnMsg + "\n"
		}
		sysutil.AppendToFile(OutFile, result)
		fmt.Println(filenum.Files)
		for _, file := range filenum.Files {
			fmt.Println(file)
			sysutil.AppendToFile(OutTmpFile, file.Name)
		}
		sysutil.AppendToFile(OutFile, result)
		sysutil.WriteToFile(OutTmpFile, result)
		result = entegor.SaveData(stCode, cfgItem, now, Data, walkDir)
		fmt.Println(result)
	}
	//	entegor.SaveData()
}
