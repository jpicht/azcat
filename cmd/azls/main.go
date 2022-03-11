package main

import (
	"github.com/jpicht/azcat/internal"
	"github.com/jpicht/azcat/pkg/azcat"
)

func main() {
	internal.Main(azcat.EMode.List())
}
