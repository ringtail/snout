package statfs

import "golang.org/x/sys/unix"

type DiskInfo struct {
	Available  int64
	Capacity   int64
	Usage      int64
	Inodes     int64
	InodesFree int64
	InodesUsed int64
}

func GetDiskInfo(path string) (*DiskInfo, error) {
	statfs := &unix.Statfs_t{}
	err := unix.Statfs(path, statfs)
	if err != nil {
		return nil, err
	}

	// Available is blocks available * fragment size
	available := int64(statfs.Bavail) * int64(statfs.Bsize)

	// Capacity is total block count * fragment size
	capacity := int64(statfs.Blocks) * int64(statfs.Bsize)

	// Usage is block being used * fragment size (aka block size).
	usage := (int64(statfs.Blocks) - int64(statfs.Bfree)) * int64(statfs.Bsize)

	inodes := int64(statfs.Files)
	inodesFree := int64(statfs.Ffree)
	inodesUsed := inodes - inodesFree

	df := &DiskInfo{
		Available:  available,
		Capacity:   capacity,
		Usage:      usage,
		Inodes:     inodes,
		InodesFree: inodesFree,
		InodesUsed: inodesUsed,
	}

	return df, nil
}

