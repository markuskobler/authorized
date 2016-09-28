package main

import (
	"github.com/markuskobler/authorized"

	_ "github.com/markuskobler/authorized/exec"
)

func main() {
	authorized.Run()
}
