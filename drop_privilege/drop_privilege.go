package drop_privilege

/*
#include <sys/types.h>
#include <pwd.h>
#include <stdlib.h>
*/
import "C" //strange, it cannot be declared together with the others import.. is it a golang bug?

import (
	"syscall"
	"unsafe"
)

type Passwd struct {
	Uid   uint32
	Gid   uint32
	Dir   string
	Shell string
}

func Getpwnam(name string) *Passwd {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cpw := C.getpwnam(cname)
	return &Passwd{
		Uid:   uint32(cpw.pw_uid),
		Gid:   uint32(cpw.pw_gid),
		Dir:   C.GoString(cpw.pw_dir),
		Shell: C.GoString(cpw.pw_shell),
	}
}

func DropPrivileges(uid int, gid int) (error, bool) {
	if err := syscall.Setuid(uid); err != nil {
		return err, false
	}
	if err := syscall.Setgid(gid); err != nil {
		return err, false
	}
	return nil, true
}
