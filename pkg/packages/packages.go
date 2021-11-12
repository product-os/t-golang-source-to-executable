package packages

import (
	"errors"
	"sync"
)

type Installer interface {
	Install(...string) error
}

var (
	installers map[string]Installer
	mx         sync.Mutex
)

func Register(distro string, impl Installer) {
	mx.Lock()
	defer mx.Unlock()
	if installers == nil {
		installers = make(map[string]Installer, 1)
	}
	installers[distro] = impl
}

func Install(distro string, pkgs ...string) error {
	if installer, ok := installers[distro]; ok {
		return installer.Install(pkgs...)
	}
	return errors.New("unknown distro")
}
