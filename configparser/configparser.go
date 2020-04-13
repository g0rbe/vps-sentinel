// Package configparser parse and check the configurations to vps-sentinel
package configparser

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// Config is the tructure to store configuration settings
type Config struct {
	PortProtocol  []string
	ProcessEnable bool
	ProcessSort   string
	SMTPServer    string
	SMTPPort      int
	SMTPUser      string
	SMTPPassword  string
	SMTPRecipient string
}

// Parse used to parse and check the configurations in the gven config file
func Parse(path string) (Config, error) {

	conf := Config{}

	// Load the config file
	cfg, err := ini.Load(path)

	if err != nil {
		return conf, fmt.Errorf("failed to parse %s: %s", path, err)
	}

	// Parse ports protocol
	conf.PortProtocol = cfg.Section("port").Key("protocol").Strings(",")
	if len(conf.PortProtocol) != 0 {

		for _, v := range conf.PortProtocol {

			if v != "tcp" && v != "tcp6" && v != "udp" && v != "udp6" {
				return conf, fmt.Errorf("failed to parse port->protocol: invalid option: %s", v)
			}
		}
	}

	// Parse process->enable
	conf.ProcessEnable, err = cfg.Section("process").Key("enable").Bool()
	if err != nil {
		return conf, fmt.Errorf("failed to parse process->enable: invalid option %s", err)
	}

	conf.ProcessSort = cfg.Section("process").Key("sort").String()

	if conf.ProcessSort != "pid" && conf.ProcessSort != "name" &&
		conf.ProcessSort != "cpu" && conf.ProcessSort != "memory" {

		return conf, fmt.Errorf("failed to parse process->sort: invalid option: %s",
			conf.ProcessSort)
	}

	// Parse smtp->server
	conf.SMTPServer = cfg.Section("smtp").Key("server").String()

	if conf.SMTPServer == "" {
		return conf, fmt.Errorf("failed to parse 'smtp->server': empty or not exist")
	}

	// Parse smtp->port
	conf.SMTPPort, err = cfg.Section("smtp").Key("port").Int()

	if err != nil {
		return conf, fmt.Errorf("failed to parse 'smtp->port': %s", err)
	} else if conf.SMTPPort > 65535 || conf.SMTPPort < 0 {
		return conf, fmt.Errorf("invalid port number in smtp->port: %d", conf.SMTPPort)
	}

	// Parse smtp->user
	conf.SMTPUser = cfg.Section("smtp").Key("user").String()

	if conf.SMTPUser == "" {
		return conf, fmt.Errorf("failed to parse 'smtp->user': empty or not exist")
	}

	// Parse smtp->password
	conf.SMTPPassword = cfg.Section("smtp").Key("password").String()

	if conf.SMTPPassword == "" {
		return conf, fmt.Errorf("failed to parse 'smtp->password': empty or not exist")
	}

	// Parse smtp->recipient
	conf.SMTPRecipient = cfg.Section("smtp").Key("recipient").String()

	if conf.SMTPRecipient == "" {
		return conf, fmt.Errorf("failed to parse 'smtp->recipient': empty or not exist")
	}

	return conf, nil
}
