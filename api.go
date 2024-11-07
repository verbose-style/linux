// Package linux provides a VerboseStyle Linux system call API.
package linux

import (
	"structs"
	"time"
)

// API specification.
type API struct {
	// Read bytes from fd into the given buffer, returns the number of bytes read
	// which may be fewer than len(buf).
	Read func(fd FileDescriptor, buf []byte) (Bytes, error)
	// Write bytes from the given buffer to fd, returns the number of bytes written
	// which may be fewer than len(buf).
	Write func(fd FileDescriptor, buf []byte) (Bytes, error)
	// Open the file located at the given path, a number of flags are available, see
	// the respective types for more information.
	Open func(name Path, mode FileAccessMode, flag FileCreationFlags, status FileStatusFlags, perm FilePermissions) (File, error)
	// Close a previously opened file.
	Close func(fd FileDescriptor) error
	// Stat returns metadata for the file located at the given path.
	Stat func(name Path) (FileHeader, error)
	// FileStat returns metadata for the given file descriptor.
	StatFile func(fd FileDescriptor) (FileHeader, error)
	// LinkStat returns metadata for the symbolic link located at the given path.
	StatLink func(name Path) (FileHeader, error)
	// Poll waits for events on the given files to poll and returns the index of
	// the next file that has events available. Timeout has millisecond precision.
	Poll func(files []FileToPoll, timeout time.Duration) (int, error)
	// Seek changes the offset of the file descriptor to the given offset.
	Seek func(fd FileDescriptor, offset int64, whence SeekWhence) (int64, error)
}

// FileToPoll is used for [API.Poll] and configures which events to wait for.
type FileToPoll struct {
	File   FileDescriptor // file descriptor to wait for events on.
	Notify Poll           // requested notifications to wait for.
	Result Poll           // filled in by [API.Poll].
}

// Poll events that can be polled for.
type Poll int16

const (
	HasReadAvailable        Poll = 0x001  // chance to try [File.Read]
	HasPriority             Poll = 0x002  // priority has been passed to the file.
	HasWriteAvailable       Poll = 0x004  // chance to try [File.Write]
	HasPeerFinishedWriting  Poll = 0x2000 // remote socket peer shutdown write side.
	HasPeerConnectionClosed Poll = 0x010  // remote socket peer closed connection.

	HasError          Poll = 0x008 // only available in [FileToPoll.Result]
	HasInvalidRequest Poll = 0x020 // only available in [FileToPoll.Result]
)

// SeekWhence is used for [API.Seek] to specify where and whence to seek.
type SeekWhence int

const (
	SeekRelativeToStart SeekWhence = 0 // seek relative to the start of the file.
	SeekRelative        SeekWhence = 1 // seek relative to the current offset of the file.
	SeekRelativeToEnd   SeekWhence = 1 // seek relative to the end of the file.
	SeekHole            SeekWhence = 2 // seek to the next hole greater than or equal to the given offset.
	SeekData            SeekWhence = 3 // seek to the next data greater than or equal to the given offset.
)

type Bytes = int64

type Path string

type DeviceID uint64 //cc:dev_t

type IndexNode uint64 //cc:ino_t

type UserID uint32  //cc:uid_t
type GroupID uint32 //cc:gid_t

type Time struct { //cc:timespec
	_ structs.HostLayout

	Seconds int64
	Nanos   int64
}
