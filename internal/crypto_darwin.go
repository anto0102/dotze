//go:build darwin

package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	ioctlReadTermios  = 0x40487413 // TIOCGETA su macOS
	ioctlWriteTermios = 0x80487414 // TIOCSETA su macOS
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
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), ioctlReadTermios, uintptr(unsafe.Pointer(&oldState))); err != 0 {
		return "", err
	}
	newState := oldState
	newState.Lflag &^= syscall.ECHO
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), ioctlWriteTermios, uintptr(unsafe.Pointer(&newState))); err != 0 {
		return "", err
	}
	defer syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), ioctlWriteTermios, uintptr(unsafe.Pointer(&oldState)))

	reader := bufio.NewReader(tty)
	password, _ := reader.ReadString('\n')
	fmt.Println()
	return strings.TrimRight(password, "\r\n"), nil
}
