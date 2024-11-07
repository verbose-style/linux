package linux

import "syscall"

// MaxRead is the maximum number of bytes that can be read in a single call to [File.Read].
const MaxRead Bytes = 0x7ffff000

// FileCreationFlags affect the semantics of the [Kernel.Open] operation.
type FileCreationFlags int

const (
	// Enable the close-on-exec flag for the new file descriptor. Specifying this flag permits a program to avoid additional
	// [File.SetStatus] operations to set the [FileClosesOnExecute] flag.
	//
	// Note that the use of this flag is essential in some multithreaded programs, because using a separate [File.SetStatus]
	// operation to set the [FileClosesOnExecute] flag does not suffice to avoid race conditions where one thread
	// opens a file descriptor and attempts to set its close-on-exec flag using [File.SetStatus] at the same time as another
	// thread does a [Kernel.Fork] plus [Kernel.Execute]. Depending on the order of execution, the race may lead to the
	// file descriptor returned by [Kernel.Open] being unintentionally leaked to the program // executed by the child process
	// created by [Kernel.Fork]. (This kind of race is in principle possible for any system call that creates a file descriptor
	// whose close-on-exec flag should be set, and various other Linux system calls provide an equivalent of the [FileCloseOnExecute]
	// flag to deal with this problem.)
	FileCloseOnExecute FileCreationFlags = syscall.O_CLOEXEC

	// FileCreateIfNeeded creates a regular file if pathname does not exist. The owner (user ID) of the new file is set to the
	// effective user ID of the process. The group ownership (group ID) of the new file is set either to the effective group ID
	// of the process (System V semantics) or to the group ID of the parent directory (BSD semantics). The behavior depends on
	// whether the set-group-ID mode bit is set on the parent directory: if that bit is set, then BSD semantics apply; otherwise,
	// System V semantics apply. For some filesystems, the behavior also depends on the bsdgroups and sysvgroups mount options
	// described in [Kernel.Mount].
	FileCreateIfNeeded FileCreationFlags = syscall.O_CREAT

	// FileAssertDirectory will check if the pathname is not a directory, in which case the open will fail. This flag was added
	// in Linux 2.1.126, to avoid denial-of-service problems if [Kernel.OpenDirectory] is called on a FIFO or tape device.
	FileAssertDirectory FileCreationFlags = syscall.O_DIRECTORY

	// FileAssertCreation ensures that this call creates the file: if this flag is specified in conjunction with [FileCreateIfNeeded],
	// and pathname already exists, then [Kernel.Open] fails with the error [ErrOpenAlreadyExists].
	//
	// When these two flags are specified, symbolic links are not followed: if pathname is a symbolic link, then [Kernel.Open] fails
	// regardless of where the symbolic link points. In general, the behavior of [FileAssertCreation] is undefined if it is used without
	// [FileCreateIfNeeded]. There is one exception: on Linux 2.6 and later, [FileAssertCreation] can be used without [FileCreateIfNeeded]
	// if pathname refers to a block device. If the block device is in use by the system (e.g., mounted), [Kernel.Open] fails with the
	// error [ErrOpenBusy]. On NFS Filesystems, [FileAssertCreation] is supported only when using NFSv3 or later on kernel 2.6 or later.
	// In NFS environments wheren[FileAssertCreation] support is not provided, programs that rely on it for performing locking tasks will
	// contain a race condition. Portable programs that want to perform atomic file locking using a lockfile, and need to avoid reliance
	// on NFS support for [FileAssertCreation], can create a unique file on the same filesystem (e.g., incorporating hostname and PID), and
	// use [Kernel.Link] to make a link to the lockfile. If [Kernel.Link] returns a nil error, the lock is successful. Otherwise, use
	// [Kernel.Stat] on the unique file to check if its link count has increased to 2, in which case the lock is also successful.
	FileAssertCreation FileCreationFlags = syscall.O_EXCL

	// FileIsNotTheTerminal, if pathname refers to a terminal device it will not become the process's controlling terminal even
	// if the process does not have one.
	FileIsNotTheTerminal FileCreationFlags = syscall.O_NOCTTY

	// If the trailing component (i.e., basename) of pathname is a symbolic link, then the open fails, with the error [ErrOpenLoop]. Symbolic
	// links in earlier components of the pathname will still be followed. (Note that the [ErrOpenLoop] error that can occur in this case is
	// indistinguishable from the case where an open fails because there are too many symbolic links found while resolving components in the
	// prefix part of the pathname.)
	// This flag is a FreeBSD extension, which was added in Linux 2.1.126, and has subsequently been standardized in POSIX.1-2008.
	FileTrapSymbolicLink FileCreationFlags = syscall.O_NOFOLLOW

	// FileTemporaryInside creates an unnamed temporary regular file. The pathname argument specifies a directory; an unnamed inode will be
	// created in that directory's filesystem. Anything written to the resulting file will be lost when the last file descriptor is closed,
	// unless the file is given a name.
	FileTemporaryInside FileCreationFlags = 020000000

	// FileTruncatedToZero will truncated the file to length 0 if the file already exists and is a regular file and the access mode allows writing
	// (i.e., is [FileAccessReadWrite] or [FileAccessWriteOnly]). If the file is a FIFO or terminal device file, the flag is ignored. Otherwise,
	// the effect of [FileTruncatedToZero] is unspecified.
	FileTruncatedToZero FileCreationFlags = syscall.O_TRUNC
)

// FileStatusFlags affect the semantics of subsequent I/O operations. These can be retrieved and (in some cases) modified;
// see [File.Status] for details.
type FileStatusFlags int

const (
	// FileAppend opens the file in append mode. Before each [File.Write], the file offset is positioned at the end of the file,
	// as if with [File.Seek]. The modification of the file offset and the write operation are performed as a single atomic step.
	// FileAppend may lead to corrupted files on NFS filesystems if more than one process appends data to a file at once. This
	// is because NFS filesystems do not support appending to a file, so the client [Kernel] has to simulate it, which can't
	// be done without a race condition.
	FileAppend FileStatusFlags = syscall.O_APPEND

	// FileAsync enables signal-driven I/O: generate a signal ([SignalIO] by default, but this can be changed via [File.SetStatus])
	// when input or output becomes possible on this file descriptor. This feature is available only for terminals, pseudoterminals,
	// sockets, and (since Linux 2.6) pipes and FIFOs. See [File.Status] for further details.
	FileAsync FileStatusFlags = syscall.O_ASYNC

	// FileDirect tries to minimize cache effects of the I/O to and from this file. In general this will degrade performance, but
	// it is useful in special situations, such as when applications do their own caching. File I/O is done directly to/from
	// user-space buffers. The [FileDirect] flag on its own makes an effort to transfer data synchronously, but does not give the
	// guarantees of the [FileSync] flag that data and necessary metadata are transferred. To guarantee synchronous I/O, [FileSync]
	// must be used in addition to [FileDirect].
	//
	// A semantically similar (but deprecated) interface for block devices is described in [Kernel.Raw].
	FileDirect FileStatusFlags = syscall.O_DIRECT

	// FileSyncData means write operations on the file will complete according to the requirements of synchronized I/O data integrity
	// completion. By the time [File.Write] (and similar) return, the output data has been transferred to the underlying hardware,
	// along with any file metadata that would be required to retrieve that data (i.e., as though each [File.Write] was followed by a
	// call to [File.SyncData]).
	FileSyncData FileStatusFlags = syscall.O_DSYNC

	// Allow files whose sizes cannot be represented in an int (but can be represented in an int64) to be opened.
	FileSize64 = syscall.O_LARGEFILE

	// FileDoNotUpdateAccessTime doesn't update the file last access time (st_atime in the inode) when the file is [File.Read].
	// This flag can be employed only if one of the following conditions is true:
	//
	//  - The effective UID of the process matches the owner UID of the file.
	//  - The calling process has the [CapabilityOwner] in its user namespace and the owner UID of the file has a
	//    mapping in the namespace.
	//
	// This flag is intended for use by indexing or backup programs, where its use can significantly reduce the amount of disk activity.
	// This flag may not be effective on all filesystems. One example is NFS, where the server maintains the access time.
	FileDoNotUpdateAccessTime FileStatusFlags = syscall.O_NOATIME

	// FileNonBlocking means the file is opened in nonblocking mode where possible. Neither the [Kernel.Open] nor any subsequent I/O
	// operations on the file descriptor which is returned will cause the calling process to wait.
	//
	// Note that the setting of this flag has no effect on the operation of [Kernel.Poll], [Kernel.Select], [Kernel.EventPoll], and
	// similar, since those interfaces merely inform the caller about whether a file descriptor is "ready", meaning that an I/O
	// operation performed on the file descriptor with the [FileNonBlocking] flag clear would not block.
	// Note that this flag has no effect for regular files and block devices; that is, I/O operations will (briefly) block when device
	// activity is required, regardless of whether [FileNonBlocking] is set. Since [FileNonBlocking] semantics might eventually be
	// implemented, applications should not depend upon blocking behavior when specifying this flag for regular files and block devices.
	FileNonBlocking = syscall.O_NONBLOCK

	// FilePath obtains a [File] that can be used for two purposes: to indicate a location in the filesystem tree and to perform operations
	// that act purely at the file descriptor level. The file itself is not opened, and many file operations will fail with an error.
	FilePath = 010000000

	// Write operations on the file will complete according to the requirements of synchronized I/O file integrity completion (by contrast
	// with the synchronized I/O data integrity completion provided by [FileDataSync].)
	//
	// By the time [File.Write] (or similar) returns, the output data and associated file metadata have been transferred to the underlying
	// hardware (i.e., as though each [File.Write] was followed by a call to [File.Sync]).
	FileSync = syscall.O_SYNC
)

// FilePermissions mode bits.
type FilePermissions uint32

const (
	FileReadableByUser   FilePermissions = syscall.S_IRUSR // file is readable by its owner
	FileWritableByUser   FilePermissions = syscall.S_IWUSR // file is writable by its owner
	FileExecutableByUser FilePermissions = syscall.S_IXUSR // file is executable by its owner

	FileReadableByGroup   FilePermissions = syscall.S_IRGRP // file is readable by its group
	FileWritableByGroup   FilePermissions = syscall.S_IWGRP // file is writable by its group
	FileExecutableByGroup FilePermissions = syscall.S_IXGRP // file is executable by its group

	FileReadableByOthers   FilePermissions = syscall.S_IROTH // file is readable by others
	FileWritableByOthers   FilePermissions = syscall.S_IWOTH // file is writable by others
	FileExecutableByOthers FilePermissions = syscall.S_IXOTH // file is executable by others

	FileExecutesAsOwner FilePermissions = syscall.S_ISUID // file will be executed as if it were executed by the owner of the file
	FileExecutesAsGroup FilePermissions = syscall.S_ISGID // file will be executed as if it were executed by the group of the file

	// FilesInheritGroup is only applicable to directories and means files created there inherit their group ID from the
	// directory, not from the effective group ID of the creating process
	FilesInheritGroup FilePermissions = syscall.S_ISGID

	// FilesLockedToOwner is only applicable to directories and means means that a file in that directory can be
	// renamed or deleted only by the owner of the file, by the owner of the directory, and by a privileged process.
	FilesLockedToOwner FilePermissions = syscall.S_ISVTX

	DirectorySearchableByUser   FilePermissions = syscall.S_IXUSR // directory is searchable by its owner
	DirectorySearchableByGroup  FilePermissions = syscall.S_IXGRP // directory is searchable by its group
	DirectorySearchableByOthers FilePermissions = syscall.S_IXOTH // directory is searchable by others
)

// FileAccessMode request opening the file read-only, write-only, or read/write, respectively.
type FileAccessMode int

const (
	FileAccessReadOnly  FileAccessMode = syscall.O_RDONLY // open the file read-only
	FileAccessWriteOnly FileAccessMode = syscall.O_WRONLY // open the file write-only
	FileAccessReadWrite FileAccessMode = syscall.O_RDWR   // open the file read-write
)

const FileRelativeToWorkingDirectory = -100 // AT_FDCWD
