// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.
//

package ps

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// Nps process struct
type Nps struct {
	// Pid  int32
	Pid  int
	Name string
}

// ToInt convert []int32 to []int
func ToInt(pid []int32) (res []int) {
	for _, v := range pid {
		res = append(res, int(v))
	}

	return
}

// GetPid get the process id
func GetPid() int {
	return os.Getpid()
}

// Pids get the all process id
func Pids() ([]int, error) {
	ids, err := process.Pids()
	return ToInt(ids), err
}

// PidExists determine whether the process exists
func PidExists(pid int) (bool, error) {
	return process.PidExists(int32(pid))
}

// Process get the all process struct
func Process() ([]Nps, error) {
	var npsArr []Nps
	pid, err := process.Pids()
	if err != nil {
		return npsArr, err
	}

	for i := 0; i < len(pid); i++ {
		nps, _ := process.NewProcess(pid[i])
		names, _ := nps.Name()

		np := Nps{
			int(pid[i]),
			names,
		}

		npsArr = append(npsArr, np)
	}

	return npsArr, err
}

// FindName find the process name by the process id
func FindName(pid int) (string, error) {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}

	return nps.Name()
}

// FindNames find the all process name
func FindNames() ([]string, error) {
	var strArr []string
	pid, err := process.Pids()

	if err != nil {
		return strArr, err
	}

	for i := 0; i < len(pid); i++ {
		nps, _ := process.NewProcess(pid[i])
		names, _ := nps.Name()

		strArr = append(strArr, names)
	}

	return strArr, err
}

// FindIds finds the all processes named with a subset
// of "name" (case insensitive),
// return matched IDs.
func FindIds(name string) ([]int, error) {
	var pids []int
	nps, err := Process()
	if err != nil {
		return pids, err
	}

	name = strings.ToLower(name)
	for i := 0; i < len(nps); i++ {
		psname := strings.ToLower(nps[i].Name)
		abool := strings.Contains(psname, name)
		if abool {
			pids = append(pids, nps[i].Pid)
		}
	}

	return pids, err
}

// FindPath find the process path by the process pid
func FindPath(pid int) (string, error) {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}

	return nps.Exe()
}

// Run command shell
func Run(path string) ([]byte, error) {
	cmdName := "/bin/bash"
	params := "-c"
	if runtime.GOOS == "windows" {
		cmdName = "cmd"
		params = "/c"
	}

	cmd := exec.Command(cmdName, params, path)
	output, err := cmd.Output()

	return output, err
}

// IsRun return the process is runing or not
func IsRun(pid int) (bool, error) {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return false, err
	}

	return nps.IsRunning()
}

// Status return the process status
func Status(pid int) ([]string, error) {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return []string{}, err
	}

	return nps.Status()
}

// Kill kill the process by PID
func Kill(pid int) error {
	p := os.Process{Pid: pid}
	return p.Kill()
}
