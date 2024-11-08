package internal

import (
	"fmt"
	"testing"

	"verbose.style/linux"
)

// #include <linux/fcntl.h>
// #include <linux/stat.h>
// #include <linux/unistd.h>
// #include <linux/mman.h>
import "C"

func assert[T comparable](t *testing.T, a, b T) {
	t.Helper()
	if a != b {
		t.Fatal(fmt.Sprintf("%v != %v", a, b))
	}
}

func Test(t *testing.T) {
	var _ linux.FileCreationFlags
	assert(t, linux.FileCloseOnExecute, C.O_CLOEXEC)
	assert(t, linux.FileCreateIfNeeded, C.O_CREAT)
	assert(t, linux.FileAssertDirectory, C.O_DIRECTORY)
	assert(t, linux.FileAssertCreation, C.O_EXCL)
	assert(t, linux.FileIsNotTheTerminal, C.O_NOCTTY)
	assert(t, linux.FileTrapSymbolicLink, C.O_NOFOLLOW)
	assert(t, linux.FileTemporaryInside, C.O_TMPFILE)
	assert(t, linux.FileTruncatedToZero, C.O_TRUNC)
	var _ linux.FileStatusFlags
	assert(t, linux.FileAppend, C.O_APPEND)
	assert(t, linux.FileAsync, C.FASYNC)
	assert(t, linux.FileDirect, C.O_DIRECT)
	assert(t, linux.FileSyncData, C.O_DSYNC)
	assert(t, linux.FileDoNotUpdateAccessTime, C.O_NOATIME)
	assert(t, linux.FileNonBlocking, C.O_NONBLOCK)
	assert(t, linux.FilePath, C.O_PATH)
	assert(t, linux.FileSync, C.O_SYNC)
	var _ linux.FilePermissions
	assert(t, linux.FileReadableByGroup, C.S_IRGRP)
	assert(t, linux.FileReadableByOthers, C.S_IROTH)
	assert(t, linux.FileWritableByUser, C.S_IWUSR)
	assert(t, linux.FileWritableByGroup, C.S_IWGRP)
	assert(t, linux.FileWritableByOthers, C.S_IWOTH)
	assert(t, linux.FileExecutableByUser, C.S_IXUSR)
	assert(t, linux.FileExecutableByGroup, C.S_IXGRP)
	assert(t, linux.FileExecutableByOthers, C.S_IXOTH)
	assert(t, linux.FileExecutesAsOwner, C.S_ISUID)
	assert(t, linux.FileExecutesAsGroup, C.S_ISGID)
	assert(t, linux.FilesInheritGroup, C.S_ISGID)
	assert(t, linux.FilesLockedToOwner, C.S_ISVTX)
	assert(t, linux.DirectorySearchableByUser, C.S_IXUSR)
	assert(t, linux.DirectorySearchableByGroup, C.S_IXGRP)
	assert(t, linux.DirectorySearchableByOthers, C.S_IXOTH)
	var _ linux.FileAccessMode
	assert(t, linux.FileAccessReadOnly, C.O_RDONLY)
	assert(t, linux.FileAccessWriteOnly, C.O_WRONLY)
	assert(t, linux.FileAccessReadWrite, C.O_RDWR)
	var _ linux.MemoryProtection
	assert(t, linux.MemoryNotAccessible, C.PROT_NONE)
	assert(t, linux.MemoryAllowReads, C.PROT_READ)
	assert(t, linux.MemoryAllowWrites, C.PROT_WRITE)
	assert(t, linux.MemoryAllowExecution, C.PROT_EXEC)
	assert(t, linux.MemoryAllowAtomics, C.PROT_SEM)
	var _ linux.MapType
	assert(t, linux.MapShared, C.MAP_SHARED)
	assert(t, linux.MapPrivate, C.MAP_PRIVATE)
	assert(t, linux.MapSharedValidateFlags, C.MAP_SHARED_VALIDATE)
	var _ linux.Map
	assert(t, linux.MapAnonymous, C.MAP_ANONYMOUS)
	assert(t, linux.Map32Bit, C.MAP_32BIT)
	assert(t, linux.MapExactAddress, C.MAP_FIXED)
	assert(t, linux.MapExactAddressOnce, C.MAP_FIXED_NOREPLACE)
	assert(t, linux.MapGrowsDown, C.MAP_GROWSDOWN)
	assert(t, linux.MapHugeTables, C.MAP_HUGETLB)
	assert(t, linux.MapHuge2MB, C.MAP_HUGE_2MB)
	assert(t, linux.MapHuge1GB, C.MAP_HUGE_1GB)
	assert(t, linux.MapKeepAwayFromSwap, C.MAP_LOCKED)
	assert(t, linux.MapDoNotReserveSwap, C.MAP_NORESERVE)
	assert(t, linux.MapPopulate, C.MAP_POPULATE)
	assert(t, linux.MapStack, C.MAP_STACK)
	assert(t, linux.MapSync, C.MAP_SYNC)
	assert(t, linux.MapUninitialized, C.MAP_UNINITIALIZED)

	assert(t, linux.FileRelativeToWorkingDirectory, C.AT_FDCWD)
}
