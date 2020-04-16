# vps-sentinel

## The idea

There are lot of servers which are just hanging on the Internet.

People try to host their own WordPress sites, mail servers, etc..., but once everything has set up, they never login again. They just simply don't care, what is happening on their server.

Playing with databases, configuring their CMS on the admin panel, but forget the underlying system. The only goal is take care of the service, the underlying system is not a priority.

## The implementation

The operation is simple:

- Collect informations about the server.
- Send a report to the administrator.
- If the admin reads it, he does more for the server than before.

By default `vps-sentinel` runs every day at 4:00 AM with `systemd` as timer.

At this time, only debian based systems are supported. 

#### Current features

- Basic informations about the server
    - Average system load
    - Free / total memory
    - Free / total swap
    - Uptime in day
- List open ports (as configured)
    - `tcp` = IPv4 TCP
    - `tcp6` = IPv6 TCP
    - `udp` = IPv4 UDP
    - `udp6` = IPv6 UDP
    - List every port which is in listening state
    - Show the port number and its process name
- Show runnig processes, as a `top` like list (if enabled)
    - Attributes:
        - pid
        - Executable name
        - Name who runs the process
        - CPU usage
        - Memory usage
    - Sort the list (as configured)
- Run ClamAV in the selected folders

#### TODO

- Parse log files
- Check `systemd` sercvices
- Run auto update if `apt-daily*` is not enabled and reboot if required 

## Reuirements

To build:

```
golang >= 1.10
```

To run with ClamAV:

```
clamav
```

To run without ClamAV:

```

```

## Install

```
git clone https://github.com/g0rbe/vps-sentinel
cd vps-sentinel
chmod +x install.sh
sudo ./install.sh conf
sudo ./install.sh install
```

## Remove

```
sudo ./install.sh remove
```

## Update

Update and replace the existing binary.
Install will not overwrite the existing configuration file!
**Please watch out for new configurations in `vps-sentinel.conf`!**

```
git pull
sudo ./install.sh install
```

## Build

To build your own binary from source:

```
sudo ./install.sh build
sudo ./install.sh install
```

## Configuration

The command below opens the config files with `nano`.

If `vps-sentinel` is installed, then `conf` opens the installed files, else
opens the repositorie's files (thats will be moved to the right place at `install`)

To close without save: `ctrl + x`

To save: `ctrl + s` and `ctrl + x`

```
sudo ./install.sh conf
```

# The report

Example report:

```
############## System informations ##############

- Average system loads (1/5/15): 0.11, 0.11 0.04
- Free memory: 979.27 MiB (total: 1945.09 MiB)
- Free swap: 0.00 MiB (total: 0.00 MiB)
- Uptime: 16.409 day(s)

###### List of interfaces and its IP addresses ######

+------+----------------------------+
| eth0 | 127.0.0.1/32               |
+------+----------------------------+

##### Open ports (tcp) #####

+------+-----------------+
| PORT | PROCESS         |
+------+-----------------+
| 22   | sshd            |
| 80   | nginx           |
| 443  | nginx           |
+------+-----------------+

##### Open ports (tcp6) #####

+------+---------+
| PORT | PROCESS |
+------+---------+
| 22   | sshd    |
| 80   | nginx   |
| 443  | nginx   |
+------+---------+

##### Open ports (udp) #####

+------+-----------------+
| PORT | PROCESS         |
+------+-----------------+
+------+-----------------+

##### Open ports (udp6) #####

+------+---------+
| PORT | PROCESS |
+------+---------+
+------+---------+

################################# List of processes #################################

+-------+-----------------+-----------------+--------+--------------+
| PID   | NAME            | USER            | CPU    | MEMORY (MIB) |
+-------+-----------------+-----------------+--------+--------------+
| 1     | vps-sentinel    | root            | 90.000 | 11           |
| 2     | system          | user            | 1.00   | 12           |
+-------+-----------------+-----------------+--------+--------------+

###################### ClamAV scan in /tmp #######################

/tmp/virus: Eicar-Signature FOUND

----------- SCAN SUMMARY -----------
Known viruses: 6822011
Engine version: 0.102.2
Scanned directories: 1
Scanned files: 0
Infected files: 0
Data scanned: 0.00 MB
Data read: 0.00 MB (ratio 0.00:1)
Time: 20.289 sec (0 m 20 s)
```