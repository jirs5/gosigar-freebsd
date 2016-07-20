// +build linux,cgo

package cgroup

/*
#include <unistd.h>
*/
import "C"

func GetClockTicks() int {
	return int(C.sysconf(C._SC_CLK_TCK))
}
