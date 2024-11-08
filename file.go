package linux

import (
	"structs"
	"sync/atomic"
	"syscall"
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

// MapIntoMemory maps the entire file into memory and returns it.
func (f *File) MapIntoMemory(mtype MapType, prot MemoryProtection, flags Map) (MappedMemory, error) {
	head, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return f.Linux.MapFileIntoMemory(nil, int(head.Size), prot, mtype, flags, f.Descriptor, 0)
}

// Close the file.
func (f *File) Close() error {
	if !f.Closed.Swap(true) {
		return f.Linux.Close(f.Descriptor)
	}
	return nil
}

// MaxRead is the maximum number of bytes that can be read in a single call to [File.Read].
const MaxRead Bytes = 0x7ffff000

// FileCreationFlags affect the semantics of the [Kernel.Open] operation.
type FileCreationFlags int

const (
	FileCloseOnExecute   FileCreationFlags = 0x80000  // close the file automatically on [Kernel.Execute].
	FileCreateIfNeeded   FileCreationFlags = 0x40     // create the file if it does not exist.
	FileAssertDirectory  FileCreationFlags = 0x10000  // fail to open if the path is not a directory.
	FileAssertCreation   FileCreationFlags = 0x80     // fail to open if the file already exists.
	FileIsNotTheTerminal FileCreationFlags = 0x100    // if the pathname is a terminal, it shouldn't become the controlling terminal for the process.
	FileTrapSymbolicLink FileCreationFlags = 0x20000  // if the trailing component is a symbolic link, don't follow it, open it directly.
	FileTemporaryInside  FileCreationFlags = 0x410000 // creates an unnamed temporary file inside the provided directory
	FileTruncatedToZero  FileCreationFlags = 0x200    // resets the file to length 0, writes will overwrite any existing content.
)

// FileStatusFlags affect the semantics of subsequent I/O operations. These can be retrieved and (in some cases) modified;
// see [File.Status] for details.
type FileStatusFlags int

const (
	FileAppend                FileStatusFlags = 0x400    // append data to the end of the file when writing.
	FileAsync                 FileStatusFlags = 0x2000   // emit [SignalIO] whenever input or output becomes available.
	FileDirect                FileStatusFlags = 0x4000   // avoid cache where possible and use underlying hardware directly
	FileSyncData              FileStatusFlags = 0x1000   // all [File.Write] operations are automatically followed by a [File.SyncData].
	FileDoNotUpdateAccessTime FileStatusFlags = 0x40000  // request that the access time of the file is not updated on [File.Read]
	FileNonBlocking           FileStatusFlags = 0x800    // return "resource temporarily unavailable" if a read/write would block
	FilePath                  FileStatusFlags = 0x200000 // file is opened as a reference-only, no read/write operations are allowed.
	FileSync                  FileStatusFlags = 0x101000 // all [File.Write] operations are automatically followed by a [File.Sync].
)

// FilePermissions mode bits.
type FilePermissions uint32

const (
	FileReadableByUser   FilePermissions = syscall.S_IRUSR // file is readable by its owner
	FileReadableByGroup  FilePermissions = syscall.S_IRGRP // file is readable by its group
	FileReadableByOthers FilePermissions = syscall.S_IROTH // file is readable by others

	FileWritableByUser   FilePermissions = syscall.S_IWUSR // file is writable by its owner
	FileWritableByGroup  FilePermissions = syscall.S_IWGRP // file is writable by its group
	FileWritableByOthers FilePermissions = syscall.S_IWOTH // file is writable by others

	FileExecutableByUser   FilePermissions = syscall.S_IXUSR // file is executable by its owner
	FileExecutableByGroup  FilePermissions = syscall.S_IXGRP // file is executable by its group
	FileExecutableByOthers FilePermissions = syscall.S_IXOTH // file is executable by others

	FileExecutesAsOwner FilePermissions = syscall.S_ISUID // file will be executed as if it were executed by the owner of the file
	FileExecutesAsGroup FilePermissions = syscall.S_ISGID // file will be executed as if it were executed by the group of the file

	FilesInheritGroup  FilePermissions = syscall.S_ISGID // files created in this directory inherit their group ID from the directory
	FilesLockedToOwner FilePermissions = syscall.S_ISVTX // files in this directory can only be renamed or deleted by owners.

	DirectorySearchableByUser   FilePermissions = syscall.S_IXUSR // directory is searchable by its owner
	DirectorySearchableByGroup  FilePermissions = syscall.S_IXGRP // directory is searchable by its group
	DirectorySearchableByOthers FilePermissions = syscall.S_IXOTH // directory is searchable by others
)

// FileAccessMode request opening the file read-only, write-only, or read/write, respectively.
type FileAccessMode int

const (
	FileAccessReadOnly  FileAccessMode = syscall.O_RDONLY // enable reads
	FileAccessWriteOnly FileAccessMode = syscall.O_WRONLY // enable writes
	FileAccessReadWrite FileAccessMode = syscall.O_RDWR   // enable both reads and writes
)

const FileRelativeToWorkingDirectory = -100 // AT_FDCWD
