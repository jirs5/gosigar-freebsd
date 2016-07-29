package cgroup

import (
	"path/filepath"
)

// Stats contains metrics and limits from each of the cgroup subsystems.
type Stats struct {
	CPU           *CPUSubsystem           `json:"cpu"`
	CPUAccounting *CPUAccountingSubsystem `json:"cpuacct"`
	Memory        *MemorySubsystem        `json:"memory"`
	BlockIO       *BlockIOSubsystem       `json:"blkio"`
}

// Reader reads cgroup metrics and limits.
type Reader struct {
	// Mountpoint of the root filesystem. Defaults to / if not set. This can be
	// useful for example if you mount / as /rootfs inside of a container.
	rootfsMountpoint  string
	ignoreRootCgroups bool              // Ignore a cgroup when its path is "/".
	cgroupMountpoints map[string]string // Mountpoints for each subsystem (e.g. cpu, cpuacct, memory, blkio).
}

// NewReader creates and returns a new Reader.
func NewReader(rootfsMountpoint string, ignoreRootCgroups bool) (*Reader, error) {
	if rootfsMountpoint == "" {
		rootfsMountpoint = "/"
	}

	// Determine what subsystems are supported by the kernel.
	subsystems, err := SupportedSubsystems(rootfsMountpoint)
	if err != nil {
		return nil, err
	}

	// Locate the mountpoints of those subsystems.
	mountpoints, err := SubsystemMountpoints(rootfsMountpoint, subsystems)
	if err != nil {
		return nil, err
	}

	return &Reader{
		rootfsMountpoint:  rootfsMountpoint,
		ignoreRootCgroups: ignoreRootCgroups,
		cgroupMountpoints: mountpoints,
	}, nil
}

// GetStatsForProcess returns cgroup metrics and limits associated with a process.
func (r *Reader) GetStatsForProcess(pid int) (*Stats, error) {
	// Read /proc/[pid]/cgroup to get the paths to the cgroup metrics.
	paths, err := ProcessCgroupPaths(r.rootfsMountpoint, pid)
	if err != nil {
		return nil, err
	}

	// Build the full path for the subsystems we are interested in.
	cgroupsPaths := map[string]string{}
	for _, interestedSubsystem := range []string{"blkio", "cpu", "cpuacct", "memory"} {
		path, found := paths[interestedSubsystem]
		if !found {
			continue
		}

		if path == "/" && r.ignoreRootCgroups {
			continue
		}

		subsystemMount, found := r.cgroupMountpoints[interestedSubsystem]
		if !found {
			continue
		}

		cgroupsPaths[interestedSubsystem] = filepath.Join(r.rootfsMountpoint, subsystemMount, path)
	}

	// Collect stats from each cgroup subsystem associated with the task.
	stats := Stats{}
	if path, found := cgroupsPaths["blkio"]; found {
		stats.BlockIO = &BlockIOSubsystem{}
		err := stats.BlockIO.Get(path)
		if err != nil {
			return nil, err
		}
	}
	if path, found := cgroupsPaths["cpu"]; found {
		stats.CPU = &CPUSubsystem{}
		err := stats.CPU.Get(path)
		if err != nil {
			return nil, err
		}
	}
	if path, found := cgroupsPaths["cpuacct"]; found {
		stats.CPUAccounting = &CPUAccountingSubsystem{}
		err := stats.CPUAccounting.Get(path)
		if err != nil {
			return nil, err
		}
	}
	if path, found := cgroupsPaths["memory"]; found {
		stats.Memory = &MemorySubsystem{}
		err := stats.Memory.Get(path)
		if err != nil {
			return nil, err
		}
	}

	// Return nil if no metrics were collected.
	if stats.BlockIO == nil && stats.CPU == nil && stats.CPUAccounting == nil && stats.Memory == nil {
		return nil, nil
	}

	return &stats, nil
}
