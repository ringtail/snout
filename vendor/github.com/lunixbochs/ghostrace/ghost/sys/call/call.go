package call

import (
	"fmt"
	"strings"
)

// TODO: where to put errno?

type Generic struct {
	Num  int
	Name string
	Args []uint64
	Ret  uint64
}

func (c *Generic) Base() *Generic {
	return c
}

func (c *Generic) String() string {
	strArgs := make([]string, len(c.Args))
	for i, v := range c.Args {
		strArgs[i] = fmt.Sprintf("0x%x", v)
	}
	args := strings.Join(strArgs, ", ")
	return fmt.Sprintf("%s(%s) = 0x%x", c.Name, args, c.Ret)
}

type Open struct {
	Generic
	Path        string
	Mode, Flags int
	Fd          int
}

func (c *Open) String() string {
	return fmt.Sprintf("open(%#v, %d, %d) = %d", c.Path, c.Mode, c.Flags, c.Fd)
}

type Close struct {
	Generic
	Fd int
}

func (c *Close) String() string {
	return fmt.Sprintf("close(%d)", c.Fd)
}

type Read struct {
	Generic
	Fd        int
	Data      []byte
	Buf, Size uint64
	Ret       int
}

func (c *Read) String() string {
	return fmt.Sprintf("read(%d, 0x%x) = (%d) %#v", c.Fd, c.Size, c.Ret, string(c.Data))
}

type Write struct {
	Generic
	Fd        int
	Data      []byte
	Buf, Size uint64
	Ret       int
}

func (c *Write) String() string {
	return fmt.Sprintf("write(%d, %#v) = %d", c.Fd, string(c.Data), c.Ret)
}

type Readv struct {
	Generic
	Fd           int
	Iovec, Count uint64
}

type Writev struct {
	Generic
	Fd           int
	Data         []byte
	Iovec, Count uint64
}

type Execve struct {
	Generic
	Path string
	Argv []string
	Envp []string
}

func (c *Execve) String() string {
	return fmt.Sprintf("execve(%s, %+v, %+v)", c.Path, c.Argv, c.Envp)
}
