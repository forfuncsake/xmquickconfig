# xmquickconfig
WiFi QuickConfig for Xiongmai IP Cameras

This repo contains an implementation of the protocol used to send SSID/Password config to some WiFi IP cameras. Please note that this method of configuration is **INSECURE** and leaks the WiFi network password into the airwaves in an unencrypted form. As such, this tool should ideally only be used for testing, but can be used in place of a prescribed mobile app for some camera models (Verified on Q-See QDB03-AU).

More details on the protocol can be found in this post:  
https://forfuncsake.github.io/post/2020/07/wifi-quick-config/


Many wireless IP Cameras only support 2.4Ghz WiFi. Therefore, this tool should be run while connected to a 2.4Ghz Network AP.

## Build & Run

The `xmquickconfig` tool is written in `go` (golang). With `go` (> v1.14 preferred) installed on your system:  
```
$ git clone https://github.com/forfuncsake/xmquickconfig
$ cd xmquickconfig
$ go build
$ ./xmquickconfig -s <ssid> -p <password> -v
```
