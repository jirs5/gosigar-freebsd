// +build !cgo !linux

package cgroup

func GetClockTicks() int {
	return 100
}
