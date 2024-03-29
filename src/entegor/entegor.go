package entegor

import (
	"fmt"
	"log"
	"ops/src/sysutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//GMESSENGER means AOMS
const GMESSENGER = "AOMS"

//LongForm define long time format
const LongForm = "2006-01-02 15:04:05"

//ShortForm define short time format
const ShortForm = "2006-01-02"

// StCode is check result code
var StCode int

//Message store check message
type Message struct {
	OutHead    string
	StCode     int
	CheckTime  string
	CheckData  string
	Threadhold string
	OutDesc    string
	WarnHead   string
	GMESSENGER string
	Script     string
	ErrCode    string
	Hostname   string
	HostIP     string
	WarnDesc   string
}

//OutMessage store check out message
type OutMessage struct {
	OutHead    string
	StCode     int
	CheckTime  string
	CheckData  string
	Threadhold string
	OutDesc    string
}

//WarningMessage store check warning message
type WarningMessage struct {
	WarnHead   string
	StCode     int
	CheckTime  string
	GMESSENGER string
	ScriptName string
	ErrCode    string
	Hostname   string
	HostIP     string
	CheckData  string
	Threadhold string
	WarnDesc   string
}

//SetLang set os locale to en_US.UTF-8
func SetLang() {
	err := os.Setenv("LANG", "en_US.UTF-8")
	if err != nil {
		log.Println(err)
	}
}

//GetHostname reture hostname of os
func GetHostname() (hostname string) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	return hostname
}

//PrepareFile prepare file used to save log ,tempout and out file.
func PrepareFile(HostIP12 string, scriptName string) (LogFile string, OutTmpFile string, OutFile string) {
	WorkDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
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
	LogFile = filepath.Join(LogDir, LogFileName)
	OutTmpFile = filepath.Join(TmpDir, OutTmpFileName)
	OutFile = filepath.Join(OutDir, OutFileName)
	sysutil.WriteToFile(OutTmpFile, "")
	sysutil.WriteToFile(LogFile, "")
	return LogFile, OutTmpFile, OutFile
}

// SaveData save check result
func SaveData(Message *Message, OutTmpFile string, OutFile string) {
	outData := fmt.Sprintf("%v=%v|%v|%v|%v|%v\n", Message.OutHead, Message.StCode, Message.CheckTime, Message.CheckData, Message.Threadhold, Message.OutDesc)
	sysutil.AppendToFile(OutTmpFile, outData)
	sysutil.AppendToFile(OutFile, outData)
	if Message.StCode != 0 {
		warnData := fmt.Sprintf("%v=%v|%v|%v|%v|%v|%v|%v|%v|%v|%v\n", Message.WarnHead, Message.StCode, Message.CheckTime, Message.GMESSENGER, Message.Script, Message.ErrCode, Message.Hostname, Message.HostIP, Message.CheckData, Message.Threadhold, Message.WarnDesc)
		sysutil.AppendToFile(OutTmpFile, warnData)
		sysutil.AppendToFile(OutFile, warnData)
	}
}

//GetHead return head
func GetHead(cfgItem string) string {
	head := strings.Split(cfgItem, "=")[0]
	return head
}

//GetWarningHead return Warning head
func GetWarningHead(cfgItem string) string {
	tempHead := strings.SplitAfter(cfgItem, ",")
	warningHead := tempHead[0] + tempHead[1] + "-1"
	return warningHead
}

// GetStCode return status code
func GetStCode(data float64, cfgItem string) int {
	StCode = 110
	vUpper := float64(999999999999)
	thordHolds := strings.Split(cfgItem, "#")[1]
	thordHold := strings.Split(thordHolds, ";")
	otherStCodeString := strings.Split(cfgItem, "#")[2]
	otherStCode, _ := strconv.Atoi(otherStCodeString)
	for _, td := range thordHold {
		codeString := strings.Split(td, ",")[0]
		code, _ := strconv.Atoi(codeString)
		thordString := strings.Split(td, ",")[1]
		thord, _ := strconv.ParseFloat(thordString, 64)
		if thord >= data && thord < vUpper {
			StCode = code
			vUpper = thord
		}
	}
	if StCode == 110 {
		StCode = otherStCode
	}
	return StCode
}

// GetGood return good range
func GetGood(cfgItem string) string {
	StCode = 110
	thordHolds := strings.Split(cfgItem, "#")[1]
	thordHold := strings.Split(thordHolds, ";")
	otherStCodeString := strings.Split(cfgItem, "#")[2]
	otherStCode, _ := strconv.Atoi(otherStCodeString)
	var Good string
	var before string
	for index, td := range thordHold {
		codeString := strings.Split(td, ",")[0]
		code, _ := strconv.Atoi(codeString)
		thordString := strings.Split(td, ",")[1]
		//thord, _ := strconv.ParseFloat(thordString, 64)
		if otherStCode == 0 {
			if index == (len(thordHold) - 1) {
				Good = "[" + thordString + " " + "Max]"
			}
		} else {
			if index == 0 && code == 0 {
				Good = "[0" + " " + thordString + "]"
				break
			} else {
				if code != 0 {
					before = thordString
				} else {
					Good = "[" + before + " " + thordString + "]"
					break
				}
			}
		}
	}
	return Good
}
