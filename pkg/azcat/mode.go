package azcat

import (
	"reflect"

	"github.com/JeffreyRichter/enum/enum"
)

type Mode uint8

var EMode = Mode(0)

func (Mode) None() Mode   { return Mode(0) }
func (Mode) Read() Mode   { return Mode(1) }
func (Mode) Write() Mode  { return Mode(2) }
func (Mode) List() Mode   { return Mode(3) }
func (Mode) Remove() Mode { return Mode(4) }

func (m Mode) String() string {
	return enum.StringInt(m, reflect.TypeOf(m))
}
