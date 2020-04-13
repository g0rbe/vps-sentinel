package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

// ProcInfo stores informations about one process
type ProcInfo struct {
	Pid         string
	Name        string
	CPUUsage    string
	MemoryUsage string
}

// getNameFromPid returns the process name that associated with the given pid
func getNameFromPid(pid int) (string, error) {

	statFilePath := fmt.Sprintf("/proc/%d/stat", pid)

	statFile, err := os.Open(statFilePath)

	if err != nil {
		return "", fmt.Errorf("failed to open %s: %s", statFilePath, err)
	}

	defer statFile.Close()

	content, err := ioutil.ReadAll(statFile)

	if err != nil {
		return "", fmt.Errorf("failed to read %s: %s", statFilePath, err)
	}

	contentSlice := strings.Fields(string(content))

	return strings.Trim(contentSlice[1], "()"), nil

}

// ListProcesses returns the processes name and its pid
func listProcesses() ([]ProcInfo, error) {

	procInfoArray := make([]ProcInfo, 0)

	pids, err := ioutil.ReadDir("/proc")

	if err != nil {
		return nil, fmt.Errorf("failed to list /proc: %s", err)
	}

	for _, pid := range pids {

		pidInt, err := strconv.Atoi(pid.Name())

		if err != nil {
			continue
		}

		procName, err := getNameFromPid(pidInt)

		if err != nil {
			return nil, fmt.Errorf("failed to get process name from pid: %s", err)
		}

		cpuUsage, err := getCPUUsageFromPid(pidInt)

		if err != nil {
			return nil, fmt.Errorf("failed to get cpu usage of pid %d: %s", pidInt, err)
		}

		memoryUsage, err := getMemoryUsageFromPid(pidInt)

		if err != nil {
			return nil, fmt.Errorf("failed to get memory usage: %s", err)
		}

		procInfo := ProcInfo{
			Pid:         pid.Name(),
			Name:        procName,
			CPUUsage:    cpuUsage,
			MemoryUsage: memoryUsage}

		procInfoArray = append(procInfoArray, procInfo)
	}

	return procInfoArray, nil
}

// GetReport returns the report of processes
func GetReport(sortField string) (string, error) {

	procInfos, err := listProcesses()

	if err != nil {
		return "", fmt.Errorf("failed to list processes: %s", err)
	}

	var sort []table.SortBy

	switch sortField {
	case "pid":
		sort = append(sort, table.SortBy{Name: "Pid", Mode: table.AscNumeric})
	case "name":
		sort = append(sort, table.SortBy{Name: "Name", Mode: table.Asc})
	case "cpu":
		sort = append(sort, table.SortBy{Name: "CPU", Mode: table.DscNumeric})
	case "memory":
		sort = append(sort, table.SortBy{Name: "Memory (MiB)", Mode: table.DscNumeric})

	}

	t := table.NewWriter()

	t.SetTitle("List of processes")

	t.AppendHeader(table.Row{"Pid", "Name", "CPU", "Memory (MiB)"})

	for _, procInfo := range procInfos {

		t.AppendRow(table.Row{procInfo.Pid, procInfo.Name,
			procInfo.CPUUsage, procInfo.MemoryUsage})

	}

	t.SortBy(sort)

	report := t.Render() + "\n\n"

	return report, nil

}
