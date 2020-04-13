package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// #include <unistd.h>
import "C"

/*
 * The calculation is based on
 * https://stackoverflow.com/questions/16726779/how-do-i-get-the-total-cpu-usage-of-an-application-from-proc-pid-stat/16736599#16736599
 */

// getHertz returns the the hertz of the system using cgo
func getHertz() float64 {

	hertz := C.sysconf(C._SC_CLK_TCK)

	return float64(hertz)

}

// getUpTime returns the systems uptime in seconds
func getUpTime() (float64, error) {

	file, err := os.Open("/proc/uptime")
	if err != nil {
		return 0, fmt.Errorf("failed to open /proc/uptime: %s", err)
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("failed to read /proc/meminfo: %s", err)
	}

	uptimeStr := strings.Fields(string(content))[0]

	uptime, err := strconv.ParseFloat(uptimeStr, 64)

	if err != nil {
		return 0, fmt.Errorf("failed to convert %s to float64: %s", uptimeStr, err)
	}

	return uptime, nil
}

// getTimeStatFromPid parse /proc/<pid>/stat and get the process's times
func getTimeStatFromPid(pid int) (utime, stime, cutime, cstime, starttime float64, err error) {

	statFilepath := fmt.Sprintf("/proc/%d/stat", pid)
	statFile, err := os.Open(statFilepath)

	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to open %s: %s", statFilepath, err)
	}

	defer statFile.Close()

	content, err := ioutil.ReadAll(statFile)

	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to read %s: %s", statFilepath, err)
	}

	contentFields := strings.Fields(string(content))

	utimeStr := contentFields[13]
	stimeStr := contentFields[14]
	cutimeStr := contentFields[15]
	cstimeStr := contentFields[16]
	starttimeStr := contentFields[21]

	utime, err = strconv.ParseFloat(utimeStr, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to convert utime's string %s: %s",
			utimeStr, err)
	}

	stime, err = strconv.ParseFloat(stimeStr, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to convert stime's string %s: %s",
			stimeStr, err)
	}

	cutime, err = strconv.ParseFloat(cutimeStr, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to convert cutime's string %s: %s",
			cutimeStr, err)
	}

	cstime, err = strconv.ParseFloat(utimeStr, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to convert cstime's string %s: %s",
			cstimeStr, err)
	}

	starttime, err = strconv.ParseFloat(starttimeStr, 64)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to convert starttime's string %s: %s",
			starttimeStr, err)
	}

	return utime, stime, cutime, cstime, starttime, nil
}

// getCPUUsageFromPid calcultes the CPU usage of the given process.
func getCPUUsageFromPid(pid int) (string, error) {

	uptime, err := getUpTime()

	if err != nil {
		return "", fmt.Errorf("failed to get uptime: %s", err)
	}

	utime, stime, cutime, cstime, starttime, err := getTimeStatFromPid(pid)

	if err != nil {
		return "", fmt.Errorf("failed to get time stats: %s", err)
	}

	hertz := getHertz()

	totalTime := utime + stime + cutime + cstime

	seconds := uptime - (starttime / hertz)

	cpuUsage := fmt.Sprintf("%.3f", 100*((totalTime/hertz)/seconds))

	return cpuUsage, nil
}
