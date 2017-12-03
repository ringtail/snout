# sysctl
[![Build Status](https://travis-ci.org/ringtail/sysctl.svg?branch=master)](https://travis-ci.org/ringtail/sysctl)
[![Codecov](https://codecov.io/gh/ringtail/sysctl/branch/master/graph/badge.svg)](https://codecov.io/gh/ringtail/sysctl)
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)       

sysctl implementation in Go

# Usage Test Cases
case1: List all system kernel settings

```
func TestAll(t *testing.T) {
	kernel_settings := All()
	if kernel_settings != nil {
		t.Log("fetch kernel_settings successfully \n")
		t.Log(kernel_settings)
		return
	}
	t.Error("Failed to fetch any kernel_setttings")
}
```
result:

```
     map[net.ipv4.tcp_moderate_rcvbuf:1
		 net.ipv4.conf.all.igmpv3_unsolicited_report_interval:1000
		 net.ipv4.conf.default.rp_filter:0
		 net.ipv6.route.max_size:4096
		 net.ipv4.conf.eth1.arp_notify:0
		 net.dccp.default.retries1:3
		 net.ipv4.conf.default.forwarding:1
		 net.ipv6.conf.all.proxy_ndp:0
		 ...
		 net.ipv4.conf.default.proxy_arp:0
		 net.ipv6.neigh.eth1.base_reachable_time_ms:30000
		 net.ipv6.neigh.eth0.app_solicit:0]
```

---

case2: Find one specific kernel settings

```
func TestFind(t *testing.T) {
	kernel_settings := Find("dev.cdrom.debug")
	if kernel_settings != nil {
		t.Log("fetch kernel_settings successfully \n")
		t.Log(kernel_settings)
		return
	}
	t.Error("Failed to fetch any kernel_setttings")
}
```
result:

```
    map[dev.cdrom.debug:0]
```

---

case3: Find one more specific kernel settings

```
func TestFindOneMore(t *testing.T) {
	kernel_settings := Find("dev.cdrom.debug", "vm.laptop_mode")
	if kernel_settings != nil {
		t.Log("fetch kernel_settings successfully \n")
		t.Log(kernel_settings)
		return
	}
	t.Error("Failed to fetch any kernel_setttings")
}
```
result:

```
    map[vm.laptop_mode:1 dev.cdrom.debug:0]
```

---

case4: Apply one kernel settings

```
func TestApply(t *testing.T) {
	err := Apply("vm.laptop_mode", "1")
	if err != nil {
		log.Errorf("Failed to apply kernel settings %v", err.Error())
		return
	}
	t.Log("Apply kernel settings successfully")
}
```
result:

```
    Apply kernel settings successfully
```
