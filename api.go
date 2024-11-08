// Package linux provides a VerboseStyle Linux system call API.
package linux

import (
	"io"
	"structs"
	"time"
	"unsafe"
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
	// MapFileIntoMemory maps the specified file into memory, using the optionally
	// specified pointer as a hint on where to map it in. If addr is nil, a suitable
	// address is chosen by the kernel. Offset must be a multiple of the system's
	// page size. When files are large, this will be more efficient than reading the
	// file into memory.
	MapFileIntoMemory func(addr unsafe.Pointer, length int, prot MemoryProtection, mtype MapType, flags Map, fd FileDescriptor, offset uintptr) (MappedMemory, error)
}

type MemoryProtection int // use by [API.MapFileIntoMemory] abd [API.ProtectMemory]

const (
	MemoryNotAccessible  MemoryProtection = 0x0 // no access allowed.
	MemoryAllowReads     MemoryProtection = 0x1 // read access allowed.
	MemoryAllowWrites    MemoryProtection = 0x2 // write access allowed.
	MemoryAllowExecution MemoryProtection = 0x4 // execute access allowed.
	MemoryAllowAtomics   MemoryProtection = 0x8 // atomic operations allowed.
)

type MapType int // used by [API.MapFileIntoMemory]

const (
	MapShared              MapType = 0x01 // persist writes back to the file.
	MapPrivate             MapType = 0x02 // copy-on-write memory.
	MapSharedValidateFlags MapType = 0x03 // [MapShared] + validate flags.
)

type Map int // used by [API.MapFileIntoMemory]

const (
	MapAnonymous        Map = 0x20       // file must be -1, just allocate anonymous memory.
	Map32Bit            Map = 0x40       // allocate memory in the first 4GB of the address space.
	MapExactAddress     Map = 0x10       // addr must be page-aligned and will be used directly.
	MapExactAddressOnce Map = 0x100000   // like [MapExactAddress] but goroutine-safe.
	MapGrowsDown        Map = 0x100      // touching the first page will grow the mapping down by a single page/
	MapHugeTables       Map = 0x40000    // use huge pages.
	MapHuge2MB          Map = 0x54000000 // use 2MB huge pages.
	MapHuge1GB          Map = 0x78000000 // use 1GB huge pages.
	MapKeepAwayFromSwap Map = 0x2000     // lock the pages in physical memory (do not swap).
	MapDoNotReserveSwap Map = 0x4000     // do not reserve swap space.
	MapPopulate         Map = 0x8000     // eagerly load the file into the map.
	MapStack            Map = 0x20000    // ensure memory is suitably setup to use for a stack.
	MapSync             Map = 0x80000    // for files that support direct mapping of persistent memory.
	MapUninitialized    Map = 0x4000000  // don't zero out pages, subject to the system secutiry policy.
)

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

// MappedMemory from a [File].
type MappedMemory interface {
	io.ReaderAt
	io.WriterAt
	io.Closer

	Len() int

	UnsafePointer() unsafe.Pointer
}
