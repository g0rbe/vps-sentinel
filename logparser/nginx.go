package logparser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

// httpError holds informations about client/server errors (4XX and 5XX)
type httpError struct {
	IP        string // Remote address
	Date      string // Time of request
	Request   string // request string
	Status    int    // Reponse code
	UserAgent string // Client's user agent
}

// parseClientErrors searches for 4XX errors
func parseClientErrors(path string) ([]httpError, error) {

	httpErrorArray := make([]httpError, 0)

	logFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %s", path, err)
	}

	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())
		filedsWithQuote := strings.Split(scanner.Text(), "\"")

		codeStr := strings.Fields(filedsWithQuote[2])[0]
		code, err := strconv.Atoi(codeStr)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to int: %s", codeStr, err)
		}

		if code < 400 || code > 500 {
			continue
		}

		httpErrorArray = append(httpErrorArray, httpError{
			IP:        fields[0],
			Date:      strings.Trim(fields[3], "["),
			Request:   filedsWithQuote[1],
			Status:    code,
			UserAgent: filedsWithQuote[5]})
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read from %s: %s", path, err)
	}

	return httpErrorArray, nil
}

// parseServerErrors searches for 5XX errors
func parseServerErrors(path string) ([]httpError, error) {

	httpErrorArray := make([]httpError, 0)

	logFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %s", path, err)
	}

	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())
		filedsWithQuote := strings.Split(scanner.Text(), "\"")

		codeStr := strings.Fields(filedsWithQuote[2])[0]
		code, err := strconv.Atoi(codeStr)

		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to int: %s", codeStr, err)
		}

		// I know that there is no 6XX errors, justwant to be sure that
		// the code is in the 500 < code < 600 interval
		if code < 500 || code > 600 {
			continue
		}

		httpErrorArray = append(httpErrorArray, httpError{
			IP:        fields[0],
			Date:      strings.Trim(fields[3], "["),
			Request:   filedsWithQuote[1],
			Status:    code,
			UserAgent: filedsWithQuote[5]})
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read from %s: %s", path, err)
	}

	return httpErrorArray, nil
}

// GetNginxClientErrors creates a report from client errors
func GetNginxClientErrors(path string) (string, error) {

	clientErrors, err := parseClientErrors(path)

	if err != nil {
		return "", fmt.Errorf("failed to parse client errors: %s", err)
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Date", "IP", "Status", "User Agent", "Request"})

	for _, clientError := range clientErrors {

		t.AppendRow(table.Row{clientError.Date, clientError.IP, clientError.Status,
			clientError.UserAgent, clientError.Request})
	}

	report := t.Render() + "\n\n"

	return report, nil
}

// GetNginxServerErrors creates a report from sefrver errors
func GetNginxServerErrors(path string) (string, error) {

	clientErrors, err := parseServerErrors(path)

	if err != nil {
		return "", fmt.Errorf("failed to parse server errors: %s", err)
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Date", "IP", "Status", "User Agent", "Request"})

	for _, clientError := range clientErrors {

		t.AppendRow(table.Row{clientError.Date, clientError.IP, clientError.Status,
			clientError.UserAgent, clientError.Request})
	}

	report := t.Render() + "\n\n"

	return report, nil
}
