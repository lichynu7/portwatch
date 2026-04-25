package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ProcessInfo holds metadata about the process that owns a socket.
type ProcessInfo struct {
	PID  int
	Name string
	Exe  string
}

// LookupProcess attempts to find the process that owns the given inode by
// scanning /proc/<pid>/fd symlinks and matching against /proc/<pid>/net
// socket inodes. It returns nil when the owning process cannot be determined
// (e.g. insufficient permissions or the socket belongs to the kernel).
func LookupProcess(inode uint64) *ProcessInfo {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			// Not a numeric directory — skip.
			continue
		}

		if ownsInode(pid, inode) {
			return buildProcessInfo(pid)
		}
	}

	return nil
}

// ownsInode reports whether the process with the given PID has an open file
// descriptor whose symlink resolves to socket:[inode].
func ownsInode(pid int, inode uint64) bool {
	fdDir := fmt.Sprintf("/proc/%d/fd", pid)
	entries, err := os.ReadDir(fdDir)
	if err != nil {
		return false
	}

	target := fmt.Sprintf("socket:[%d]", inode)

	for _, entry := range entries {
		link, err := os.Readlink(filepath.Join(fdDir, entry.Name()))
		if err != nil {
			continue
		}
		if link == target {
			return true
		}
	}

	return false
}

// buildProcessInfo reads /proc/<pid>/comm and /proc/<pid>/exe to populate a
// ProcessInfo. Partial results are returned when only some fields are readable.
func buildProcessInfo(pid int) *ProcessInfo {
	info := &ProcessInfo{PID: pid}

	// comm contains the short process name (up to 15 chars, newline-terminated).
	commBytes, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err == nil {
		info.Name = strings.TrimSpace(string(commBytes))
	}

	// exe is a symlink to the executable path.
	exePath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err == nil {
		info.Exe = exePath
	}

	return info
}
