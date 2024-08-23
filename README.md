# Traceroute

An implementation of the traceroute command in Go.

## Usage

```bash
# Run the program with the domain name or IP address as an argument
# sudo is required to open raw sockets
$ sudo go run main.go google.com
Traceroute to google.com (142.251.33.110), 64 hops max
 1 192.168.4.1 (192.168.4.1) 2.863625ms 3.210333ms 3.511709ms
 2 br1-belcarra-v100.bb.as54858.net. (64.187.165.1) 4.654792ms 4.662542ms 4.461042ms
 3 * * *
 4 agg2-sea-b-t9-4.bb.spectrumnet.us. (216.243.24.105) 6.537833ms 4.192584ms 4.649875ms
 5 be12.cr3-sea-b.bb.as11404.net. (174.127.151.13) 4.231459ms 4.15975ms 4.276084ms
 6 be10.cr3-sea-a.bb.as11404.net. (65.50.198.62) 4.405041ms 4.898458ms 4.947291ms
 7 be11.cr5-sea.bb.as11404.net. (174.127.151.10) 4.272958ms 3.630458ms 3.4455ms
 8 216.243.15.225 (216.243.15.225) 3.520042ms 4.227208ms 3.443541ms
 9 142.251.70.99 (142.251.70.99) 5.192583ms 6.023208ms 6.332375ms
10 142.251.50.175 (142.251.50.175) 4.450959ms 4.190625ms 4.345667ms
11 sea30s10-in-f14.1e100.net. (142.251.33.110) 4.363208ms 4.145375ms 3.480375ms
```
