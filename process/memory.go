package process

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// getMemoryUsageFromPid calculates the given PID's memory usage based on PSS.
// Returns the memory usage in byte.
func getMemoryUsageFromPid(pid string) (string, error) {

	totalSize := 0

	smapsFilePath := fmt.Sprintf("/proc/%s/smaps", pid)

	smapsFile, err := os.Open(smapsFilePath)

	if err != nil {
		return "", fmt.Errorf("failed to open %s: %s", smapsFilePath, err)
	}

	defer smapsFile.Close()

	scanner := bufio.NewScanner(smapsFile)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())

		if fields[0] == "Pss:" {

			partSize, err := strconv.Atoi(fields[1])

			if err != nil {
				return "", fmt.Errorf("failed to convert %s to int: %s", fields[1], err)
			}

			totalSize += partSize
		}
	}

	return strconv.Itoa(totalSize / 1024), nil
}
