package helloworld

import (
	seccomp "github.com/seccomp/libseccomp-golang"
)

func init() {
	println(seccomp.GetLibraryVersion())
}
