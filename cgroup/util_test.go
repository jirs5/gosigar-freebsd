package cgroup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportedSubsystems(t *testing.T) {
	subsystems, err := SupportedSubsystems("testdata")
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, subsystems, 12)
	assertContains(t, subsystems, "cpuset")
	assertContains(t, subsystems, "cpu")
	assertContains(t, subsystems, "cpuacct")
	assertContains(t, subsystems, "blkio")
	assertContains(t, subsystems, "memory")
	assertContains(t, subsystems, "devices")
	assertContains(t, subsystems, "freezer")
	assertContains(t, subsystems, "net_cls")
	assertContains(t, subsystems, "perf_event")
	assertContains(t, subsystems, "net_prio")
	assertContains(t, subsystems, "hugetlb")
	assertContains(t, subsystems, "pids")
}

func TestSubsystemMountpoints(t *testing.T) {
	subsystems := map[string]struct{}{}
	subsystems["blkio"] = struct{}{}
	subsystems["cpu"] = struct{}{}
	subsystems["cpuacct"] = struct{}{}
	subsystems["cpuset"] = struct{}{}
	subsystems["devices"] = struct{}{}
	subsystems["freezer"] = struct{}{}
	subsystems["hugetlb"] = struct{}{}
	subsystems["memory"] = struct{}{}
	subsystems["perf_event"] = struct{}{}

	mountpoints, err := SubsystemMountpoints("testdata", subsystems)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "/sys/fs/cgroup/blkio", mountpoints["blkio"])
	assert.Equal(t, "/sys/fs/cgroup/cpu", mountpoints["cpu"])
	assert.Equal(t, "/sys/fs/cgroup/cpuacct", mountpoints["cpuacct"])
	assert.Equal(t, "/sys/fs/cgroup/cpuset", mountpoints["cpuset"])
	assert.Equal(t, "/sys/fs/cgroup/devices", mountpoints["devices"])
	assert.Equal(t, "/sys/fs/cgroup/freezer", mountpoints["freezer"])
	assert.Equal(t, "/sys/fs/cgroup/hugetlb", mountpoints["hugetlb"])
	assert.Equal(t, "/sys/fs/cgroup/memory", mountpoints["memory"])
	assert.Equal(t, "/sys/fs/cgroup/perf_event", mountpoints["perf_event"])
}

func TestProcessCgroupPaths(t *testing.T) {
	paths, err := ProcessCgroupPaths("testdata", 985)
	if err != nil {
		t.Fatal(err)
	}

	path := "/docker/b29faf21b7eff959f64b4192c34d5d67a707fe8561e9eaa608cb27693fba4242"
	assert.Equal(t, path, paths["blkio"])
	assert.Equal(t, path, paths["cpu"])
	assert.Equal(t, path, paths["cpuacct"])
	assert.Equal(t, path, paths["cpuset"])
	assert.Equal(t, path, paths["devices"])
	assert.Equal(t, path, paths["freezer"])
	assert.Equal(t, path, paths["memory"])
	assert.Equal(t, path, paths["net_cls"])
	assert.Equal(t, path, paths["net_prio"])
	assert.Equal(t, path, paths["perf_event"])
	assert.Len(t, paths, 10)
}

func assertContains(t testing.TB, m map[string]struct{}, key string) {
	_, contains := m[key]
	if !contains {
		t.Errorf("map is missing key %v, map=%+v", key, m)
	}
}
