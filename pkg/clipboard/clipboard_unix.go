// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || linux || netbsd || openbsd || solaris || dragonfly
// +build freebsd linux netbsd openbsd solaris dragonfly

package clipboard

import (
	"errors"
	"log/slog"
	"os/exec"
)

const (
	xsel        = "xsel"
	xclip       = "xclip"
	flatpakPath = "/run/host/bin/"
)

var (
	// Primary choose primary mode on unix
	Primary bool

	pasteCmdArgs, copyCmdArgs []string

	xselPasteArgs = []string{xsel, "--output", "--clipboard"}
	xselCopyArgs  = []string{xsel, "--input", "--clipboard"}

	xselPasteFlatpakArgs = []string{flatpakPath + xsel, "--output", "--clipboard"}
	xselCopyFlatpakArgs  = []string{flatpakPath + xsel, "--input", "--clipboard"}

	xclipPasteArgs = []string{xclip, "-out", "-selection", "clipboard"}
	xclipCopyArgs  = []string{xclip, "-in", "-selection", "clipboard"}

	xclipPasteFlatpakArgs = []string{flatpakPath + xclip, "-out", "-selection", "clipboard"}
	xclipCopyFlatpakArgs  = []string{flatpakPath + xclip, "-in", "-selection", "clipboard"}

	errMissingCommands = errors.New("no clipboard utilities available. Please install xsel or xclip")
)

func init() {
	pasteCmdArgs = xclipPasteArgs
	copyCmdArgs = xclipCopyArgs

	if _, err := exec.LookPath(xclip); err == nil {
		slog.Info("Using xclip")
		return
	} else {
		_, err = exec.LookPath(flatpakPath + xclip)
		if err == nil {
			pasteCmdArgs = xclipPasteFlatpakArgs
			copyCmdArgs = xclipCopyFlatpakArgs
			slog.Info("Using xclip from flatpak")
			return
		}
	}

	pasteCmdArgs = xselPasteArgs
	copyCmdArgs = xselCopyArgs

	if _, err := exec.LookPath(xsel); err == nil {
		slog.Info("Using xsel")
		return
	} else {
		_, err = exec.LookPath(flatpakPath + xsel)
		if err == nil {
			pasteCmdArgs = xselPasteFlatpakArgs
			copyCmdArgs = xselCopyFlatpakArgs
			slog.Info("Using xsel from flatpak")
			return
		}
	}

	slog.Error("No clipboard utilities available. Please install xsel or xclip")
	Unsupported = true
}

func getPasteCommand() *exec.Cmd {
	if Primary {
		pasteCmdArgs = pasteCmdArgs[:1]
	}
	return exec.Command(pasteCmdArgs[0], pasteCmdArgs[1:]...)
}

func getCopyCommand() *exec.Cmd {
	if Primary {
		copyCmdArgs = copyCmdArgs[:1]
	}
	return exec.Command(copyCmdArgs[0], copyCmdArgs[1:]...)
}

func readAll() (string, error) {
	if Unsupported {
		return "", errMissingCommands
	}

	pasteCmd := getPasteCommand()
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func writeAll(text string) error {
	if Unsupported {
		return errMissingCommands
	}
	copyCmd := getCopyCommand()
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}

	return copyCmd.Wait()
}
