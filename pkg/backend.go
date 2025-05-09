package pkg

import (
	"os/exec"
	"syscall"
)

type WallpaperBackend interface {
	SetImage(imagepath string) error
	IsInstalled() bool
}

type SwayBg struct {
}

func (s *SwayBg) SetImage(imagepath string) error {
	// This is where you'd call swaybg or something similar
	exec.Command("pkill", "swaybg").Run()
	cmd := exec.Command("swaybg", "-i", imagepath)

	// Redirect std pipes to /dev/null
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Detach from parent
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// Start the process
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	// Do not wait for it
	_ = cmd.Process.Release()

	println("swaybg detached with PID", cmd.Process.Pid)

	return err
}

func (s *SwayBg) IsInstalled() bool {
	// This is where you'd call swaybg or something similar
	return true
}

func InitBackend(map[string]string) WallpaperBackend {
	//TODO prolly here a backend validator func
	return &SwayBg{}
}
