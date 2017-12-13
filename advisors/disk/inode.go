package disk

import (
	"github.com/ringtail/snout/types"
	"github.com/ringtail/snout/storage"
	"github.com/ringtail/snout/collectors/disk"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"fmt"
)

const (
	INODE_THRESHOLD     = 0.8
	INODE_USAGE_SYMPTOM = "INODE_USAGE_SYMPTOM"
)

func GetInodeSymptom(metrics_tree *storage.MetricsTree) types.Symptom {
	disk_status := metrics_tree.FindSection(disk.DISK_STATUS)
	if disk_status == nil {
		log.Warnf("Failed to get disk status from metrics storage.")
		return nil
	}
	inodeUsed, _ := strconv.Atoi(disk_status.Find("InodesUsed"))
	inodeAll, _ := strconv.Atoi(disk_status.Find("Inodes"))
	log.Debug(inodeUsed, inodeAll)
	if (float32(inodeUsed) / float32(inodeAll)) > INODE_THRESHOLD {
		desc := fmt.Sprintf("Current system can have %d inodes, but used %d inodes", inodeAll, inodeUsed)
		adviseDescs := []string{
			"The inode is a data structure in a Unix-style file system that describes a " +
				"filesystem object such as a file or a directory.Maybe you have too much fragment files",

			"You can use `sysctl -w vm.vfs_cache_pressure=200`.default value is 100, 100 means kernel's reclaim" +
				" thememory versus pagecache and swap. Increasing this value increases the rate at which VFS caches are reclaimed",

			"You can also enlarge the inode settings by remount disk,But we strongly suggest you not to do it.",
		}
		inode_usage_symptom := types.CreateTextDefaultSymptom(INODE_USAGE_SYMPTOM, desc, adviseDescs)
		return inode_usage_symptom
	}
	return nil
}
