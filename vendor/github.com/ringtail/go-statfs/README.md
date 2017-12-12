# go-statfs
[![Build Status](https://travis-ci.org/ringtail/go-statfs.svg?branch=master)](https://travis-ci.org/ringtail/go-statfs)
[![Codecov](https://codecov.io/gh/ringtail/go-statfs/branch/master/graph/badge.svg)](https://codecov.io/gh/ringtail/go-statfs)
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)   
Statfs implement in golang 

## Description
go-statfs is similar as command `df` and `df -i`，you can get such a struct from it.

```
type DiskInfo struct {
	Available  int64
	Capacity   int64
	Usage      int64
	Inodes     int64
	InodesFree int64
	InodesUsed int64
}
```

## Usage

```
func TestGetDiskInfo(t *testing.T) {
	disUsage, err := GetDiskInfo("/")
	if err != nil {
		t.Errorf("Failed to get disk info,because of %s", err.Error())
	}
	t.Logf("available %d, capacity %d, usage %d, inodes %d, inodesFree %d, inodesUsed %d",
		disUsage.Available, disUsage.Capacity, disUsage.Usage, disUsage.Inodes, disUsage.InodesFree, disUsage.InodesUsed)
}
```

## License
This software is released under the Apache 2.0 license.