package linux

import "reflect"

type Error[T any] struct{ ErrMethods[T] }

type ErrMethods[T any] byte

func (n ErrMethods[T]) Error() string {
	field := reflect.TypeFor[T]().Field(int(n))
	if field.Tag != "" {
		return string(field.Tag)
	}
	return field.Name
}

func (n ErrMethods[T]) parse(err error) error {
	if err == nil {
		return nil
	}
	var msg = err.Error()
	var types T
	var value = reflect.ValueOf(&types).Elem()
	var rtype = reflect.TypeFor[T]()
	for i := range rtype.NumField() {
		field := rtype.Field(i)
		if string(field.Tag) == msg {
			value.Field(i).Field(0).SetUint(uint64(i))
			return value.Field(i).Interface().(error)
		}
	}
	return err
}

func (n ErrMethods[T]) Types() T {
	var types T
	var value = reflect.ValueOf(&types).Elem()
	for i := range value.NumField() {
		value.Field(i).Field(0).SetUint(uint64(i))
	}
	return types
}

// ReadError returned by [API.Read], [File.Read] operations.
type ReadError Error[struct {
	WouldBlock  ReadError `resource temporarily unavailable` // file requested as non-blocking and the read would block, try again later.
	BadFile     ReadError `bad file descriptor`              // file is not valid.
	Fault       ReadError `bad address`                      // buffer is outside the accessible address space.
	Interrupted ReadError `interrupted system call`          // read was interrupted by a signal.
	Invalid     ReadError `invalid argument`                 // file is not suitable for reading.
	IO          ReadError `I/O error`                        // an I/O error occurred.
	Directory   ReadError `is a directory`                   // directories cannot be read.
}]

// WriteError returned by [API.Write], [File.Write] operations.
type WriteError Error[struct {
	WouldBlock     WriteError `resource temporarily unavailable` // file requested as non-blocking and the write would block, try again later.
	BadFile        WriteError `bad file descriptor`              // file is not valid.
	NoDestination  WriteError `destination address required`     // files is a datagram socket and requires a destination address.
	QuotaExhausted WriteError `disk quota exceeded`              // user's quota of space has run out.
	Fault          WriteError `bad address`                      // buffer is outside the accessible address space.
	TooMuch        WriteError `file too large`                   // file exceeds the maximum file size.
	Interrupted    WriteError `interrupted system call`          // write was interrupted by a signal.
	Invalid        WriteError `invalid argument`                 // file is not suitable for writing.
	IO             WriteError `I/O error`                        // an I/O error occurred.
	NoMoreSpace    WriteError `no space left on device`          // device has no more space.
	NotPermitted   WriteError `operation not permitted`          // file is not open for writing.
	BrokenPipe     WriteError `broken pipe`                      // write to a closed pipe with no readers.
}]

// OpenError returned by [API.Open] operations.
type OpenError Error[struct {
	AccessDenied   OpenError `permission denied`                // one of the directories is missing the search/execute permission bit, or wrong user.
	BadFile        OpenError `bad file descriptor`              // file is not valid.
	Busy           OpenError `device or resource busy`          // file is mounted and cannot be opened.
	QuotaExhausted OpenError `disk quota exceeded`              // user's quota of space has run out.
	AlreadyExists  OpenError `file exists`                      // file already exists and [FileCreateIfNeeded] and [FileAssertCreation] were used.
	Fault          OpenError `bad address`                      // pathname is outside your accessible address space.
	FileTooLarge   OpenError `file too large`                   // file exceeds architecture file size limit.
	NotPermitted   OpenError `operation not permitted`          // permissions missing.
	ReadOnly       OpenError `read-only file system`            // file is on a read-only filesystem and write access was requested.
	FileInUse      OpenError `file in use`                      // file is in use.
	WouldBlock     OpenError `resource temporarily unavailable` // file requested as non-blocking and the open would block, try again later.
}]

// CloseError returned by [API.Close] operations.
type CloseError Error[struct {
	BadFile        CloseError `bad file descriptor`     // file is not valid.
	Interrupted    CloseError `interrupted system call` // close was interrupted by a signal.
	IO             CloseError `I/O error`               // an I/O error occurred.
	QuotaExhausted CloseError `disk quota exceeded`     // user's quota of space has run out, can be returned on close when IO is being buffered.
	NoMoreSpace    CloseError `no space left on device` // device has no more space, can be returned on close when IO is being buffered.
}]

// StatError returned by [API.Stat], [API.StatLink] and [API.StatFile] operations.
type StatError Error[struct {
	DoesNotExist     StatError `no such file or directory`             // an element in the path does not exist.
	AccessDenied     StatError `permission denied`                     // one of the directories is missing the search/execute permission bit.
	BadFile          StatError `bad file descriptor`                   // file is not valid.
	Fault            StatError `bad address`                           // path string is corrupted.
	Invalid          StatError `invalid argument`                      // invalid flags
	Loop             StatError `too many levels of symbolic links`     // recursion limit reached.
	NameTooLong      StatError `file name too long`                    // unsupported file name
	OutOfMemory      StatError `cannot allocate memory`                // kernel is out of memory
	NotDirectory     StatError `not a directory`                       // a component of the path prefix is not a directory.
	StatFileTooLarge StatError `value too large for defined data type` // file size is 64 bits and the system is 32 bits.
}]

// PollError returned by [API.Poll] operations.
type PollError Error[struct {
	Fault       PollError `bad address`             // files to poll is nil or points out of the accessible address space.
	Interrupted PollError `interrupted system call` // poll was interrupted by a signal.
	Invalid     PollError `invalid argument`        // too many files to poll or the timeout is invalid.
	OutOfMemory PollError `cannot allocate memory`  // kernel is out of memory
}]

// SeekError returned by [API.Seek] operations.
type SeekError Error[struct {
	BadFile  SeekError `bad file descriptor`                   // file is not valid.
	Invalid  SeekError `invalid argument`                      // invalid whence or resulting offset out of bounds.
	NotFound SeekError `no such device or address`             // [SeekData] or [SeekHole] could not find a suitable offset within bounds.
	Overflow SeekError `value too large for defined data type` // resulting offset is too large to fit in an int64.
	Illegal  SeekError `illegal seek`                          // pipes/sockets are not seekable.
}]

// MapError returned by [API.MapFileIntoMemory] operations.
type MapError Error[struct {
	AccessDenied  MapError `permission denied`                     // non-regular file or [FileAccessMode] is incompatible with [MemoryProtection]
	Locked        MapError `resource temporarily unavailable`      // file is locked, or too much locked memory in-use.
	BadFile       MapError `bad file descriptor`                   // file is not valid and [MapAnonymous] not set.
	AlreadyExists MapError `file exists`                           // [MapExactAddressOnce] is set and the address is already mapped.
	Invalid       MapError `invalid argument`                      // addr, length or offset is invalid and/or [MapType] missing.
	TooManyFiles  MapError `too many open files`                   // process has too many files open.
	Unsupported   MapError `no such device`                        // file system does not support mapping.
	OutOfMemory   MapError `cannot allocate memory`                // kernel is out of memory, and/or addr exceeds the virtual address space.
	Overflow      MapError `value too large for defined data type` // resulting offset is too large to fit in an int64.
	NotPermitted  MapError `operation not permitted`               // file is not readable or writable, and/or process huge page capabilities.
}]

// ProtectMemoryError returned by [API.ProtectMemory] operations.
type ProtectMemoryError Error[struct {
	AccessDenied ProtectMemoryError `permission denied`      // mapped file [FileAccessMode] is incompatible with [MemoryProtection].
	Invalid      ProtectMemoryError `invalid argument`       // addr is not aligned to page size, is invalid or flags are invalid.
	OutOfMemory  ProtectMemoryError `cannot allocate memory` // kernel is out of memory
}]

// HeapError returned by [API.Heap] operations.
type HeapError Error[struct {
	OutOfMemory HeapError `cannot allocate memory` // no more memory available.
}]
