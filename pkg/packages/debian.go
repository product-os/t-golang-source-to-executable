package packages

import (
	"os"
	"strings"

	"github.com/product-os/t-golang-source-to-executable/pkg/shell"
)

func init() {
	Register("debian", new(Debian))
}

type Debian struct{}

func (d *Debian) Install(pkgs ...string) error {
	_, err := shell.Run("apt-get", []string{"install", strings.Join(pkgs, " ")}, nil, os.Stdout, nil)
	return err
}
