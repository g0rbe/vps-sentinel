package sysinfo

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

// MemInfo holds informations about the systems memory
type memInfo struct {
	MemTotal  float64
	MemFree   float64
	SwapTotal float64
	SwapFree  float64
}

// getMemInfo return a MemInfo struct, holding informations about the systems memory
// The returned informations is in byte
func getMemInfo() (memInfo, error) {

	var info memInfo

	file, err := os.Open("/proc/meminfo")

	if err != nil {
		return info, fmt.Errorf("failed to open /proc/meminfo: %s", err)
	}

	defer file.Close()

	lines := bufio.NewScanner(file)

	for lines.Scan() {

		elems := strings.Fields(lines.Text())

		switch elems[0] {
		case "MemTotal:":
			numInKb, err := strconv.ParseFloat(elems[1], 64)
			if err != nil {
				return info, fmt.Errorf("failed to convert %s to int: %s", elems[1], err)
			}
			info.MemTotal = numInKb * 1024
		case "MemFree:":
			numInKb, err := strconv.ParseFloat(elems[1], 64)
			if err != nil {
				return info, fmt.Errorf("failed to convert %s to int: %s", elems[1], err)
			}
			info.MemFree = numInKb * 1024
		case "SwapTotal:":
			numInKb, err := strconv.ParseFloat(elems[1], 64)
			if err != nil {
				return info, fmt.Errorf("failed to convert %s to int: %s", elems[1], err)
			}
			info.SwapTotal = numInKb * 1024
		case "SwapFree:":
			numInKb, err := strconv.ParseFloat(elems[1], 64)
			if err != nil {
				return info, fmt.Errorf("failed to convert %s to int: %s", elems[1], err)
			}
			info.SwapFree = numInKb * 1024
		}
	}

	if err := lines.Err(); err != nil {
		return info, fmt.Errorf("error while reading /proc/meminfo: %s", err)
	}

	return info, nil
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

// getSystemLoad reutrn the average system loads, 1, 5, 15 minutes respectively
func getSystemLoad() ([]float64, error) {

	loads := []float64{0, 0, 0}

	file, err := os.Open("/proc/loadavg")

	if err != nil {
		return loads, fmt.Errorf("failed to open /proc/loadavg: %s", err)
	}

	defer file.Close()

	content, err := ioutil.ReadAll(file)

	if err != nil {
		return loads, fmt.Errorf("failed to read /proc/loadavg: %s", err)
	}

	loadsSplitStr := strings.Fields(string(content))

	loadsStr := []string{loadsSplitStr[0], loadsSplitStr[1], loadsSplitStr[2]}

	for index := range loadsStr {

		loads[index], err = strconv.ParseFloat(loadsStr[index], 64)

		if err != nil {
			return loads, fmt.Errorf("failed to convert %s to float: %s", loadsStr[index], err)
		}
	}

	return loads, nil
}

// GetFqdn returns the fqdn of the system
// This function is copied from https://github.com/Showmax/go-fqdn
func GetFqdn() string {

	hostname, err := os.Hostname()

	if err != nil {
		return "?"
	}

	addrs, err := net.LookupIP(hostname)

	if err != nil {
		return hostname
	}

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {

			ip, err := ipv4.MarshalText()

			if err != nil {
				return hostname
			}

			hosts, err := net.LookupAddr(string(ip))

			if err != nil || len(hosts) == 0 {
				return hostname
			}

			fqdn := hosts[0]

			return strings.TrimSuffix(fqdn, ".") // return fqdn without trailing dot
		}
	}

	return hostname
}

// GetSysInfo returns a report with system informations
// Current informations: system load, free/total memory/swap, uptime
func GetSysInfo() (string, error) {

	loads, err := getSystemLoad()

	if err != nil {
		return "", fmt.Errorf("failed to get system loads: %s", err)
	}

	memInfo, err := getMemInfo()

	if err != nil {
		return "", fmt.Errorf("failed to get memory informations: %s", err)
	}

	uptime, err := getUpTime()

	if err != nil {
		return "", fmt.Errorf("failed to get uptime: %s", err)
	}

	report := "System informations:\n"

	report += fmt.Sprintf("- Average system loads (1/5/15): %.2f, %.2f %.2f\n",
		loads[0], loads[1], loads[2])

	report += fmt.Sprintf("- Free memory: %.2f MiB (total: %.2f MiB)\n",
		memInfo.MemFree/1048576.0, memInfo.MemTotal/1048576.0)

	report += fmt.Sprintf("- Free swap: %.2f MiB (total: %.2f MiB)\n",
		memInfo.SwapFree/1048576.0, memInfo.SwapTotal/1048576.0)

	report += fmt.Sprintf("- Uptime: %.3f day(s)\n", uptime/86400.0)

	report += "\n"

	return report, nil
}
