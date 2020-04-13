# vps-sentinel

## The idea

There are a lot of servers which are just hanging on the Internet.
Peoples tries to host his own WordPress, mail server, etc..., but once he set up, never login again to the server.
Playing with databases, configuring his CMS on the admin panel, but forget the underlying system.
The only goal is take care of the service, the underlying system is not a priority.

Thats the point of the `vps-sentinel`: generate a report of the system on daily basis, and if you read it, you do more for the security than before.

## Reuirements

On build:

```
golang >= 1.10
```

## Install

```
chmod +x install.sh
sudo ./install.sh install
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

Configure `bin/vps-sentinel.conf` before the first install, or `/etc/vps-sentinel.conf` after install

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