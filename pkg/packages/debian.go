package packages

import (
	"os"

	"github.com/product-os/t-golang-source-to-executable/pkg/shell"
)

func init() {
	Register("debian", new(Debian))
}

type Debian struct{}

func (d *Debian) Install(pkgs ...string) error {
	var err error
	_, err = shell.Run("apt-get", []string{"update"}, nil, os.Stdout, nil)
	installArgs := []string{"install", "-y", "--no-install-recommends"}
	installArgs = append(installArgs, pkgs...)
	_, err = shell.Run("apt-get", installArgs, nil, os.Stdout, nil)
	return err
}
