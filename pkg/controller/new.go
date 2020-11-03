package controller

import (
	"io"
)

type Controller struct {
	Stdout io.Writer
	Stderr io.Writer
}
