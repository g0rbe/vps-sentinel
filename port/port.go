package port

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/jedib0t/go-pretty/table"
)

// sockinfo holds informations about listening port and its related process name
type portInfo struct {
	PortNo   string
	ProcName string
}

// inodeToProc returns the process name that belongs to the given socket inode
// Dont remove the '()', because it is the part of the final report
// Returns '?' if nothing found
func inodeToProc(inode int) (string, error) {

	// List /proc/*
	pids, err := ioutil.ReadDir("/proc")

	if err != nil {
		return "", fmt.Errorf("failed to open /proc: %s", err)
	}

	for _, pid := range pids {

		// List /proc/<pid>/fd/*
		fdDir := fmt.Sprintf("/proc/%s/fd", pid.Name())

		fds, err := ioutil.ReadDir(fdDir)

		// Skip errors, because /proc contains not just pid folders, but /proc/uptime, etc..
		if err != nil {
			continue
		}

		for _, fd := range fds {

			// Explicit use of Stat(),
			// because the inner Sys() interface not follow the fd's symlink
			var stat syscall.Stat_t
			fdPath := fmt.Sprintf("/proc/%s/fd/%s", pid.Name(), fd.Name())

			err := syscall.Stat(fdPath, &stat)

			if err != nil {
				continue
			}

			if stat.Ino == uint64(inode) {

				statPath := fmt.Sprintf("/proc/%s/stat", pid.Name())

				statFile, err := os.Open(statPath)

				if err != nil {
					return "", fmt.Errorf("failed to open %s: %s", statPath, err)
				}

				defer statFile.Close()

				content, err := ioutil.ReadAll(statFile)

				if err != nil {
					return "", fmt.Errorf("failed to read %s: %s", statPath, err)
				}

				procName := strings.Fields(string(content))[1]

				procName = strings.Trim(procName, "()")

				return procName, nil
			}
		}
	}

	return "?", nil
}

// parsePorts parses the open ports and the related process name
// Returns a 2d array of strings, ["port", "procname"]
func parsePorts(protocol string) ([]portInfo, error) {

	var listenState string

	result := make([]portInfo, 0)

	switch protocol {
	case "tcp", "tcp6":
		listenState = "0A"
	case "udp", "udp6":
		listenState = "07"
	}

	file, err := os.Open("/proc/net/" + protocol)

	if err != nil {
		return nil, fmt.Errorf("Failed to open /proc/net/%s: %s", protocol, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())

		if fields[3] == listenState {

			inode, err := strconv.Atoi(fields[9])

			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to int: %s", fields[9], err)
			}

			procName, err := inodeToProc(inode)

			if err != nil {
				return nil, fmt.Errorf("failed to get process name from inode: %s", err)
			}

			portStrHex := strings.Split(fields[1], ":")[1]

			portInt, err := strconv.ParseInt(portStrHex, 16, 0)

			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to int: %s", portStrHex, err)
			}

			portStr := strconv.Itoa(int(portInt))

			// Check wether the current port is exist in the list to disbale duplication
			isExist := false

			for _, v := range result {
				if v.PortNo == portStr {
					isExist = true
				}
			}

			if !isExist {
				result = append(result, portInfo{PortNo: portStr, ProcName: procName})
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error in scanner: %s", err)
	}

	return result, nil
}

// GetListeningPorts generates a table report of open ports and its related process
func GetListeningPorts(protocol string) (string, error) {

	ports, err := parsePorts(protocol)

	if err != nil {
		return "", fmt.Errorf("failed to parse ports: %s", err)
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Port", "Process"})

	for _, port := range ports {
		t.AppendRow(table.Row{port.PortNo, port.ProcName})
	}

	sort := []table.SortBy{table.SortBy{Name: "Port", Mode: table.AscNumeric}}

	t.SortBy(sort)

	result := t.Render() + "\n\n"

	return result, nil
}
