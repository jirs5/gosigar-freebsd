// Copyright (c) 2012 VMware, Inc.

package sigar

// #include <stdlib.h>
// #include <windows.h>
import "C"

import (
	"fmt"
	"unsafe"
)

func init() {
}

func (self *LoadAverage) Get() error {
	return nil
}

func (self *Uptime) Get() error {
	return nil
}

func (self *Mem) Get() error {
	var statex C.MEMORYSTATUSEX
	statex.dwLength = C.DWORD(unsafe.Sizeof(statex))

	succeeded := C.GlobalMemoryStatusEx(&statex)
	if succeeded == C.FALSE {
		lastError := C.GetLastError()
		return fmt.Errorf("GlobalMemoryStatusEx failed with error: %d", int(lastError))
	}

	self.Total = uint64(statex.ullTotalPhys)
	self.Free = uint64(statex.ullAvailPhys)
	self.Used = self.Total - self.Free
	vtotal := uint64(statex.ullTotalVirtual)
	self.ActualFree = uint64(statex.ullAvailVirtual)
	self.ActualUsed = vtotal - self.ActualFree

	return nil
}

func (self *Swap) Get() error {
	return notImplemented()
}

func (self *Cpu) Get() error {

	var lpIdleTime, lpKernelTime, lpUserTime C.FILETIME

	succeeded := C.GetSystemTimes(&lpIdleTime, &lpKernelTime, &lpUserTime)
	if succeeded == C.FALSE {
		lastError := C.GetLastError()
		return fmt.Errorf("GetSystemTime failed with error: %d", int(lastError))
	}

	LOT := float64(0.0000001)
	HIT := (LOT * 4294967296.0)

	idle := ((HIT * float64(lpIdleTime.dwHighDateTime)) + (LOT * float64(lpIdleTime.dwLowDateTime)))
	user := ((HIT * float64(lpUserTime.dwHighDateTime)) + (LOT * float64(lpUserTime.dwLowDateTime)))
	kernel := ((HIT * float64(lpKernelTime.dwHighDateTime)) + (LOT * float64(lpKernelTime.dwLowDateTime)))
	system := (kernel - idle)

	self.Idle = uint64(idle)
	self.User = uint64(user)
	self.Sys = uint64(system)
	return nil
}

func (self *CpuList) Get() error {
	return notImplemented()
}

func (self *FileSystemList) Get() error {
	return notImplemented()
}

func (self *ProcList) Get() error {
	return notImplemented()
}

func (self *ProcState) Get(pid int) error {
	return notImplemented()
}

func (self *ProcMem) Get(pid int) error {
	return notImplemented()
}

func (self *ProcTime) Get(pid int) error {
	return notImplemented()
}

func (self *ProcArgs) Get(pid int) error {
	return notImplemented()
}

func (self *ProcExe) Get(pid int) error {
	return notImplemented()
}

func (self *FileSystemUsage) Get(path string) error {
	var availableBytes C.ULARGE_INTEGER
	var totalBytes C.ULARGE_INTEGER
	var totalFreeBytes C.ULARGE_INTEGER

	pathChars := C.CString(path)
	defer C.free(unsafe.Pointer(pathChars))

	succeeded := C.GetDiskFreeSpaceEx((*C.CHAR)(pathChars), &availableBytes, &totalBytes, &totalFreeBytes)
	if succeeded == C.FALSE {
		lastError := C.GetLastError()
		return fmt.Errorf("GetDiskFreeSpaceEx failed with error: %d", int(lastError))
	}

	self.Total = *(*uint64)(unsafe.Pointer(&totalBytes))
	return nil
}

func notImplemented() error {
	panic("Not Implemented")
	return nil
}
