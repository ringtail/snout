package sys

import "github.com/lunixbochs/ghostrace/ghost/sys/call"

type Syscall interface {
	Base() *call.Generic
	String() string
}
