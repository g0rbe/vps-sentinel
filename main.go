package main

import (
	"fmt"
	"os"

	"github.com/g0rbe/vps-sentinel/logparser"

	"github.com/g0rbe/vps-sentinel/clamav"
	"github.com/g0rbe/vps-sentinel/ipinfo"

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

	var report string

	for _, feature := range conf.ReportStructure {

		switch feature {
		case "system":

			report += "############## " + "System informations" + " ##############\n\n"

			fmt.Printf("Getting system informations...\n")

			if sInfo, err := sysinfo.GetSysInfo(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get system informations: %s\n", err)
				report += fmt.Sprintf("Failed to get system informations: %s\n", err)
			} else {
				report += sInfo
			}
		case "ip":

			report += "###### " + "List of interfaces and its IP addresses" + " ######\n\n"

			fmt.Printf("Getting ip informations...\n")

			if netInfo, err := ipinfo.GetIPInfo(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get network informations: %s\n", err)
				report += fmt.Sprintf("Failed to get network informations: %s\n", err)
			} else {
				report += netInfo
			}
		case "port":

			// Iterate over the given protocols to get a report of listening ports
			for _, protocol := range conf.PortProtocol {

				report += "##### " + "Open ports (" + protocol + ")" + " #####\n\n"

				fmt.Printf("Getting open ports of %s...\n", protocol)

				if ports, err := port.GetListeningPorts(protocol); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to parse %s: %s\n", protocol, err)
					report += fmt.Sprintf("Failed to parse %s: %s\n", protocol, err)
				} else {
					report += ports
				}
			}
		case "process":

			report += "################################## " +
				"List of processes" + " ##################################\n\n"

			fmt.Printf("Generating a list of processes...\n")

			if procList, err := process.GetReport(conf.ProcessSort); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to list processes: %s\n", err)
				report += fmt.Sprintf("Failed to list processes: %s\n", err)
			} else {
				report += procList
			}
		case "clamav":

			// Scan with ClamAV in the given paths
			for _, path := range conf.ClamAVPath {

				report += "###################### " +
					"ClamAV scan in " + path + " #######################\n\n"

				fmt.Printf("Running ClamAV in %s...\n", path)

				if out, err := clamav.RunClamAV(path); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to run clamav on %s: %s\n", path, err)
					report += fmt.Sprintf("Failed to run clamav on %s: %s\n", path, err)
				} else {
					report += out
				}
			}
		case "log.ssh":

			report += "#################### " +
				"Accepted SSH logins" + " #####################\n\n"

			fmt.Printf("Finding accepted SSH logins...\n")

			if out, err := logparser.GetAcceptedLogins(conf.SSHLogPath); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get accepted SSH logins: %s", err)
				report += fmt.Sprintf("Failed to get accepted SSH logins: %s", err)
			} else {
				report += out
			}

			if conf.SSHParseFailed {
				report += "######## " +
					"Failed SSH logins" + " ########\n\n"

				fmt.Printf("Finding failed SSH logins...\n")

				if out, err := logparser.GetFailedLogins(conf.SSHLogPath); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to get failed SSH logins: %s", err)
					report += fmt.Sprintf("Failed to get failed SSH logins: %s", err)
				} else {
					report += out
				}
			}
		}
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
