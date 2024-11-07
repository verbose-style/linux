package linux

import (
	"structs"
	"sync/atomic"
)

// File opened with [API].
type File struct {
	Linux      *API
	Descriptor FileDescriptor
	Closed     atomic.Bool
}

// FileHeader returned by [API.Stat] provides a representation of the metadata that
// the filesystem records on the file.
type FileHeader struct { //cc:stat
	_ structs.HostLayout

	Device      DeviceID
	IndexNode   IndexNode
	HardLinks   uint64
	Permissions FilePermissions
	User        UserID
	Group       GroupID
	_           int32
	Special     DeviceID
	Size        Bytes
	BlockSize   Bytes
	BlockCount  int64

	AccessedAt         Time
	ModifiedAt         Time
	ModifiedMetadataAt Time
	_                  [3]int64
}

// FileDescriptor identifies an open file for the process.
type FileDescriptor int32

// Read implements [io.Reader]
func (f *File) Read(p []byte) (int, error) {
	n, err := f.Linux.Read(f.Descriptor, p)
	return int(n), err
}

// Write implements [io.Writer]
func (f *File) Write(p []byte) (int, error) {
	n, err := f.Linux.Write(f.Descriptor, p)
	return int(n), err
}

// Seek implements [io.Seeker]
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.Linux.Seek(f.Descriptor, offset, SeekWhence(whence))
}

// Stat returns metadata for the file located at the given path.
func (f *File) Stat() (FileHeader, error) { return f.Linux.StatFile(f.Descriptor) }

// Close the file.
func (f *File) Close() error {
	if !f.Closed.Swap(true) {
		return f.Linux.Close(f.Descriptor)
	}
	return nil
}
