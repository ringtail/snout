# snout
<p align="center">
    <img width="400"  src="http://snout.oss-cn-beijing.aliyuncs.com/snout.png">
  <p align="center">sniff bad smell about performance</p>
</p>

[![Build Status](https://travis-ci.org/ringtail/snout.svg?branch=master)](https://travis-ci.org/ringtail/snout)
[![Codecov](https://codecov.io/gh/ringtail/snout/branch/master/graph/badge.svg)](https://codecov.io/gh/ringtail/snout)
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

# What is snout

snout is command line tool to collect system metrics and give you some advice about performance

# Usage
```
=====================================================

     ███████╗███╗   ██╗ ██████╗ ██╗   ██╗████████╗
     ██╔════╝████╗  ██║██╔═══██╗██║   ██║╚══██╔══╝
     ███████╗██╔██╗ ██║██║   ██║██║   ██║   ██║
     ╚════██║██║╚██╗██║██║   ██║██║   ██║   ██║
     ███████║██║ ╚████║╚██████╔╝╚██████╔╝   ██║
     ╚══════╝╚═╝  ╚═══╝ ╚═════╝  ╚═════╝    ╚═╝

  snout is tool to improve your system performance
=====================================================


Usage:
	snout [commands|flags]

The commands & flags are:
	help 				print snout help
	version 			print the version to stdout
	--debug 			switch on debug mode

Examples:
	# start in normal mode
	snout

	# start in debug mode
	snout  --debug
```
# TIME_WAIT DEMO
```
+--------------------+--------------------------------+----------------------------------------------------------------------------------+
|      SYMPTOM       |          DESCRIPTION           |                                     ADVISES                                      |
+--------------------+--------------------------------+----------------------------------------------------------------------------------+
| TIME_WAIT_TOO_MUCH | tcp connection state           | You can reuse tcp connection by set `keepalive` in http client,set               |
|                    | `TIME_WAIT` is too much,       | `fastcgi_keep_conn` in php-fpm settings                                          |
|                    | current amount is 23           |                                                                                  |
+                    +                                +----------------------------------------------------------------------------------+
|                    |                                | You can accelerate the `TIME_WAIT` connection recycle by sysctl: sysclt -w       |
|                    |                                | net.ipv4.tcp_syncookies = 1;sysclt -w net.ipv4.tcp_tw_reuse = 1;sysclt -w        |
|                    |                                | net.ipv4.tcp_tw_recycle = 1;sysclt -w net.ipv4.tcp_fin_timeout = 30              |
+--------------------+--------------------------------+----------------------------------------------------------------------------------+
```

# Related Project
`statfs(df)` implement in golang  (<a href="https://github.com/ringtail/go-statfs">https://github.com/ringtail/go-statfs</a>)
`sysctl` implement in golang  (<a href="https://github.com/ringtail/sysctl">https://github.com/ringtail/sysctl</a>)
`netstat` implement in golang  (<a href="https://github.com/ringtail/GOnetstat">https://github.com/ringtail/GOnetstat</a>)