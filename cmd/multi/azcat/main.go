package main

import (
	"github.com/jpicht/azcat/actions"
	"github.com/jpicht/azcat/internal"
)

func main() {
	internal.Main(actions.EMode.Read())
}
