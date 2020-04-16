package main

import (
	"fmt"
	"os"

	"github.com/g0rbe/vps-sentinel/clamav"

	"github.com/g0rbe/vps-sentinel/netinfo"

	"github.com/g0rbe/vps-sentinel/process"

	"github.com/go-mail/mail"

	"github.com/g0rbe/vps-sentinel/sysinfo"

	"github.com/g0rbe/vps-sentinel/configparser"
	"github.com/g0rbe/vps-sentinel/port"
)

func main() {

	fmt.Printf("Parsing configuration file...\n")

	conf, err := configparser.Parse("/etc/vps-sentinel.conf")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse configuration file: %s\n", err)
		os.Exit(1)
	}

	report := ""

	fmt.Printf("Getting system informations...\n")

	// Get system informations
	sInfo, err := sysinfo.GetSysInfo()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get system informations: %s\n", err)
		os.Exit(1)
	}

	report += sInfo

	fmt.Printf("Getting network informations...\n")

	// Get network infomations
	netInfo, err := netinfo.GetNetinfo()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get network informations: %s\n", err)
		os.Exit(1)
	}

	report += netInfo

	// Iterate over the given protocols to get a report of listening ports
	for _, protocol := range conf.PortProtocol {

		fmt.Printf("getting open ports of %s...\n", protocol)

		ports, err := port.GetListeningPorts(protocol)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse %s: %s\n", protocol, err)
			os.Exit(1)
		}

		report += ports
	}

	// Get a report of existing processes on the system
	if conf.ProcessEnable {

		fmt.Printf("Generating a list of processes...\n")

		procList, err := process.GetReport(conf.ProcessSort)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to list processes: %s\n", err)
			os.Exit(1)
		}

		report += procList
	}

	// Scan with ClamAV in the given paths
	for _, path := range conf.ClamAVPath {

		fmt.Printf("Running ClamAV on %s...\n", path)

		out, err := clamav.RunClamAV(path)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run clamav on %s: %s\n", path, err)
			os.Exit(1)
		}

		report += out
	}

	fmt.Printf("Sending report...\n")

	subj := fmt.Sprintf("[%s] Daily report from vps-sentinel", sysinfo.GetFqdn())

	m := mail.NewMessage()
	m.SetHeader("From", conf.SMTPUser)
	m.SetHeader("To", conf.SMTPRecipient)
	m.SetHeader("Subject", subj)
	m.SetBody("text/plain", report)

	d := mail.NewDialer(conf.SMTPServer, conf.SMTPPort, conf.SMTPUser, conf.SMTPPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send mail: %s\n", err)
		os.Exit(1)
	}
}
