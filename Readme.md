# GoMonitor

## Usage
go build

gomonitor -h www.google.com -t 5 -s 30 -m get


## Flags
| Flags | Description                                 |
|:-----:| --------------------------------------------|
| u     | url                                         |
| r     | number of request (default 0 = forever)     |
| s     | seconds to wait between request             |
| m     | HTTP method (just get and post supported)   |