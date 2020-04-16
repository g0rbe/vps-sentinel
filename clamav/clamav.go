package clamav

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

// freshclam checks the status of the freshclam service.
// If freshclam daemon is not running, then update the db manually
func freschclam() error {

	cmdCheck := exec.Command("/bin/systemctl", "-q", "is-active", "clamav-freshclam.service")

	if err := cmdCheck.Run(); err == nil {
		return nil
	}

	cmdUpdate := exec.Command("/usr/bin/freshclam", "--quiet")

	if out, err := cmdUpdate.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run freshclam: %s, %s", out, err)
	}

	return nil
}

// RunClamAV runs clamscan on the given path
// Returns the report
func RunClamAV(path string) (string, error) {

	// Update ClamAV-s datavase
	if err := freschclam(); err != nil {
		return "", fmt.Errorf("faile to update database with freshclam: %s", err)
	}

	cmd := exec.Command("/usr/bin/clamscan", "-i", "-r", path)

	// Pipe stdout
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to gte stdout pipe: %s", err)
	}

	// Pipe stderr
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stderr pipe: %s", err)
	}

	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start ClamAV: %s", err)
	}

	cmdErr, err := ioutil.ReadAll(stdErr)
	if err != nil {
		return "", fmt.Errorf("failed to read from cmd's stderr: %s", err)
	}

	cmdOut, err := ioutil.ReadAll(stdOut)
	if err != nil {
		return "", fmt.Errorf("failed to read form cmd's stdout: %s", err)
	}

	if err = cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok &&
			exitErr.ExitCode() != 1 && exitErr.ExitCode() != 0 {

			return "", fmt.Errorf("failed to scan %s: %s", path, cmdErr)
		}
	}

	report := fmt.Sprintf("Report of scanning: %s\n", path)

	report += string(cmdOut) + "\n\n"

	return report, nil
}
