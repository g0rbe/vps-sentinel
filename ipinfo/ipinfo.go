package ipinfo

import (
	"fmt"
	"net"

	"github.com/jedib0t/go-pretty/table"
)

// ifaceIP holds the iface name and it IP address
type ifaceIP struct {
	Name string
	IP   string
}

// Get an array of ifaceIP
func getIPs() ([]ifaceIP, error) {

	ifaceIPs := make([]ifaceIP, 0)

	ifaces, err := net.Interfaces()

	if err != nil {
		return nil, fmt.Errorf("failed to get interfaces list: %s", err)
	}

	for _, iface := range ifaces {

		if iface.Name == "lo" {
			continue
		}

		addrs, err := iface.Addrs()

		if err != nil {
			return nil, fmt.Errorf("failed to get %s's address: %s", iface.Name, err)
		}

		for _, addr := range addrs {

			ifaceIPs = append(ifaceIPs, ifaceIP{Name: iface.Name, IP: addr.String()})

		}
	}

	return ifaceIPs, nil
}

// GetIPInfo creates a report of interfaces and its associated IP addresses
func GetIPInfo() (string, error) {

	ips, err := getIPs()

	if err != nil {
		return "", fmt.Errorf("failed to ip addresses: %s", err)
	}

	t := table.NewWriter()

	for _, ip := range ips {
		t.AppendRow(table.Row{ip.Name, ip.IP})
	}

	report := t.Render() + "\n\n"

	return report, nil

}
