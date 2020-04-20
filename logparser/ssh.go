package logparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

// acceptedSSHLogin hold inofrmations about accepted SSH logins
type acceptedLogin struct {
	Time     string
	User     string // The username which the user logged in
	IP       string // The remote IP of the logged in user
	AuthType string // The type of the authentication: pubkey, password, etc...
}

// failedSSHLogin holds informations about failed SSH logins, counts the failed login per IP
type failedLogin struct {
	IP    string
	Count int
}

// parseAcceptedLogin searches the log file for accepted logins
func parseAcceptedLogins(path string) ([]acceptedLogin, error) {

	acceptedLoginArray := make([]acceptedLogin, 0)

	logFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %s", path, err)
	}

	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())

		// Parse accepted logins
		if strings.Contains(fields[4], "sshd") && fields[5] == "Accepted" {

			timeStr := strings.Join(fields[:3], " ")

			acceptedLoginArray = append(acceptedLoginArray,
				acceptedLogin{User: fields[8], IP: fields[10], AuthType: fields[6],
					Time: timeStr})

		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read from %s: %s", path, err)
	}

	return acceptedLoginArray, nil
}

// parseFailedLogin searches the log file for failed logins
func parseFailedLogins(path string) ([]failedLogin, error) {

	failedLoginArray := make([]failedLogin, 0)

	logFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %s", path, err)
	}

	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())

		// Parse failed logins
		if strings.Contains(fields[4], "sshd") && fields[5] == "Failed" {

			var ip string

			// Saw some username in logs, which is not only one string.
			// Like: "X.X.X.X - SSH-2.0-Ope.SSH_7.4\r"
			// This username breaks the simple splitting.
			// But IP is always before "port"
			for num, value := range fields[10:] {
				if value == "port" {
					ip = fields[10+num-1]
					break
				}
			}

			newIP := true

			// Check for existing IPs
			for _, failed := range failedLoginArray {
				if failed.IP == ip {
					failed.Count++
					newIP = false
					break
				}
			}

			if newIP {
				failedLoginArray = append(failedLoginArray, failedLogin{IP: ip, Count: 1})
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read from %s: %s", path, err)
	}

	return failedLoginArray, nil
}

// GetAcceptedLogins generates a report of the accepted logins
func GetAcceptedLogins(path string) (string, error) {

	logins, err := parseAcceptedLogins(path)

	if err != nil {
		return "", fmt.Errorf("failed to parse accepted logins: %s", err)
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Time", "User", "IP", "Authentication type"})

	for _, login := range logins {
		t.AppendRow(table.Row{login.Time, login.User, login.IP, login.AuthType})
	}

	report := t.Render() + "\n\n"

	return report, nil
}

// GetFailedLogins generates a report of the failed logins
func GetFailedLogins(path string, multiple bool) (string, error) {

	logins, err := parseFailedLogins(path)

	if err != nil {
		return "", fmt.Errorf("failed to parse failed logins: %s", err)
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"IP", "Count"})

	for _, login := range logins {

		if multiple && login.Count == 1 {
			continue
		}

		t.AppendRow(table.Row{login.IP, login.Count})
	}

	sort := []table.SortBy{table.SortBy{Name: "Count", Mode: table.DscNumeric}}

	t.SortBy(sort)

	report := t.Render() + "\n\n"

	return report, nil
}
