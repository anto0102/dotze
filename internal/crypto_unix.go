//go:build !windows

package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		// fallback
		var password string
		_, err := fmt.Scanln(&password)
		return password, err
	}
	defer tty.Close()

	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, tty.Fd(), syscall.TCGETS, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		return "", err
	}
	newState := oldState
	newState.Lflag &^= syscall.ECHO
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, tty.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return "", err
	}
	defer syscall.Syscall6(syscall.SYS_IOCTL, tty.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)

	reader := bufio.NewReader(tty)
	password, _ := reader.ReadString('\n')
	fmt.Println()
	return strings.TrimRight(password, "\r\n"), nil
}
