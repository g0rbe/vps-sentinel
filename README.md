# vps-sentinel

## The idea

There are a lot of servers which are just hanging on the Internet.

Peoples tries to host his own WordPress, mail server, etc..., but once he set up, never login again to the server, dont care about what happening on his server

Playing with databases, configuring his CMS on the admin panel, but forget the underlying system.

The only goal is take care of the service, the underlying system is not a priority.

## The implementation

The operation is simple:

- Collect informations about the server
- Send a report to the administrator
- If the admin read it, he do more for the server than before

By default `vps-sentinel` runs every day at 4:00 AM with `systemd` as timer.

At this time, only debian based system are supported. 

#### Current features

- Basic informations about the server
    - Average system load
    - Free / total memory
    - free / total swap
    - Uptime in day
- List open ports (as configured)
    - `tcp` = IPv4 TCP
    - `tcp6` = IPv6 TCP
    - `udp` = IPv4 UDP
    - `udp6` = IPv6 UDP
    - List every port which is in listening state
    - Show the port number and it's process name
- Show runnig processes, as a `top` like list (if enabled)

#### TODO

- Run ClamAV on the selected folders
- Parse log files
- Check `systemd` sercvices
- Run auto update if `apt-daily*` is not enabled and reboot if required 

## Reuirements

On build:

```
golang >= 1.10
```

## Install

```
git clone https://github.com/g0rbe/vps-sentinel
cd vps-sentinel
chmod +x install.sh
sudo ./install.sh install
sudo nano /etc/vps-sentinel.conf
```

## Remove

```
sudo ./install.sh remove
```

## Update

Install wil not remove the existing configuration file!

```
git pull
sudo ./install.sh install
```

## Build

To build your own binary from source:

```
./install.sh build
sudo ./install.sh install
```

## Configuration

The command below open the config files with `nano`.

To close without save: `ctrl + x`

To save: `ctrl + s` and `ctrl + x`

```
sudo ./install.sh conf
```

### The report

Exampe report:

```
System informations:
- Average system loads (1/5/15): 0.1, 0.1 0.1
- Free memory: 703.31 MiB (total: 1000 MiB)
- Free swap: 1000 MiB (total: 1000 MiB)
- Uptime: 14.734 day(s)

+-------------------------------------+
| List of interfaces with IPs         |
+------+------------------------------+
| eth0 | 127.0.0.1/24                 |
| eth0 | fe80::e6b9:7aff:fe65:64aa/64 |
+------+------------------------------+

+---------------------+
| Open ports (tcp)    |
+-------+-------------+
| PORT  | PROCESS     |
+-------+-------------+
| 80    | nginx       |
+-------+-------------+

+--------------------+
| Open ports (tcp6)  |
+------+-------------+
| PORT | PROCESS     |
+------+-------------+
| 80    | nginx       |
+------+-------------+

+---------------------+
| Open ports (udp)    |
+-------+-------------+
| PORT  | PROCESS     |
+-------+-------------+
+-------+-------------+

+---------------------+
| Open ports (udp6)   |
+-------+-------------+
| PORT  | PROCESS     |
+-------+-------------+
+-------+-------------+

+-----------------------------------------------------------------+
| List of processes                                               |
+---------+------------------------------+---------+--------------+
| PID     | NAME                         | CPU     | MEMORY (MIB) |
+---------+------------------------------+---------+--------------+
| 1       | systemd                      | 7.293   | 3            |
| 817     | NetworkManager               | 0.089   | 7            |
| 876     | wpa_supplicant               | 0.017   | 4            |
| ...     | ...                          | ...     | ...          |
+---------+------------------------------+---------+--------------+
````