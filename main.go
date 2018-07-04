package MonkeyKiller

import (
	"os/exec"
	"log"
	"strings"
	"os"
	"io/ioutil"
)

func dealErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getAllDevice() []string{
	var resultDeviceList []string
	getDeviceCmd := exec.Command("adb", "devices")
	deviceStr, err := getDeviceCmd.CombinedOutput()
	dealErr(err)
	deviceList := strings.Split(string(deviceStr), "\n")[1:]
	for _, eachLine := range deviceList {
		if eachLine == "" {
			break
		}
		if strings.Contains(eachLine, "device") {
			deviceID := strings.Split(eachLine, "\t")[0]
			resultDeviceList = append(resultDeviceList, deviceID)
		}
	}
	log.Println("All Devices: ", resultDeviceList)
	return resultDeviceList
}

func killMonkey(deviceID string) {
	getPsCmdStr := []string{"-s", deviceID, "shell", "ps"}
	findPsCmdStr := []string{"monkey"}
	getPsCmd := exec.Command("adb", getPsCmdStr...)
	findPsCmd := exec.Command("findstr", findPsCmdStr...)
	log.Println("get process cmd: ", getPsCmd.Args)
	log.Println("find process cmd: ", findPsCmd.Args)

	findPsCmd.Stdin, _ = getPsCmd.StdoutPipe()
	findResultReader, err := findPsCmd.StdoutPipe()
	dealErr(err)
	defer findResultReader.Close()
	findPsCmd.Stderr = os.Stderr

	findPsCmd.Start()
	getPsCmd.Run()
	resultByte, _ := ioutil.ReadAll(findResultReader)
	if len(resultByte) == 0 {
		log.Println("no monkey process detected on ", deviceID)
		return
	}
	resultList := strings.Split(string(resultByte), " ")
	var monkeyPID string
	for _, arg := range resultList[1:] {
		if arg == "" {
			continue
		}
		monkeyPID = arg
		break
	}
	log.Println(monkeyPID)

	killCmdStr := []string{"-s", deviceID, "shell", "kill", "-9", monkeyPID}
	killCmd := exec.Command("adb", killCmdStr...)
	err = killCmd.Run()
	dealErr(err)
	log.Println("monkey dead: ", deviceID)
}

func main() {
	deviceList := getAllDevice()
	for _, deviceID := range deviceList {
		killMonkey(deviceID)
	}
}