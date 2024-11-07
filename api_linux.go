package linux

import (
	"syscall"
	"time"
	"unsafe"
)

func Native() *API {
	var os = new(API)
	*os = API{
		Read: func(f FileDescriptor, buf []byte) (Bytes, error) {
			count, _, err := syscall.Syscall(syscall.SYS_READ, uintptr(f), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
			return Bytes(count), new(ReadError).parse(syscall.Errno(err))
		},
		Write: func(f FileDescriptor, buf []byte) (Bytes, error) {
			count, _, err := syscall.Syscall(syscall.SYS_WRITE, uintptr(f), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
			return Bytes(count), new(WriteError).parse(syscall.Errno(err))
		},
		Open: func(path Path, access FileAccessMode, creation FileCreationFlags, status FileStatusFlags, perm FilePermissions) (File, error) {
			fd, err := syscall.Open(string(path), int(access)|int(creation)|int(status), uint32(perm))
			return File{Linux: os, Descriptor: FileDescriptor(fd)}, new(OpenError).parse(err)
		},
		Close: func(f FileDescriptor) error {
			err := syscall.Close(int(f))
			return new(CloseError).parse(err.(syscall.Errno))
		},
		Stat: func(path Path) (FileHeader, error) {
			var header FileHeader
			err := syscall.Stat(string(path), (*syscall.Stat_t)(unsafe.Pointer(&header)))
			return header, new(StatError).parse(err)
		},
		StatFile: func(fd FileDescriptor) (FileHeader, error) {
			var header FileHeader
			err := syscall.Fstat(int(fd), (*syscall.Stat_t)(unsafe.Pointer(&header)))
			return header, new(StatError).parse(err)
		},
		StatLink: func(name Path) (FileHeader, error) {
			var header FileHeader
			err := syscall.Lstat(string(name), (*syscall.Stat_t)(unsafe.Pointer(&header)))
			return header, new(StatError).parse(err)
		},
		Poll: func(files []FileToPoll, timeout time.Duration) (int, error) {
			if len(files) == 0 {
				return 0, new(PollError).Types().Fault
			}
			i, _, err := syscall.Syscall(syscall.SYS_POLL, uintptr(unsafe.Pointer(&files[0])), uintptr(len(files)), uintptr(timeout))
			return int(i), new(PollError).parse(err)
		},
		Seek: func(fd FileDescriptor, offset int64, whence SeekWhence) (int64, error) {
			o, err := syscall.Seek(int(fd), offset, int(whence))
			return int64(o), new(SeekError).parse(err)
		},
	}
	return os
}
