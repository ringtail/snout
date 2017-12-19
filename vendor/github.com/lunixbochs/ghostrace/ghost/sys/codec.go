package sys

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/lunixbochs/ghostrace/ghost/memio"
	"github.com/lunixbochs/ghostrace/ghost/sys/call"
	"github.com/lunixbochs/ghostrace/ghost/sys/num"
)

type Codec struct {
	Arch ArchType
	OS   OSType
	Mem  memio.MemIO
}

func NewCodec(arch ArchType, os OSType, mem memio.MemIO) (*Codec, error) {
	if arch != ARCH_X86_64 && os != OS_LINUX {
		return nil, errors.New("unsupported arch/os")
	}
	return &Codec{arch, os, mem}, nil
}

func (c *Codec) DecodeCall(n int, args []uint64) (Syscall, error) {
	return c.decode(n, args, 0, false)
}

func (c *Codec) DecodeRet(n int, args []uint64, ret uint64) (Syscall, error) {
	return c.decode(n, args, ret, true)
}

func (c *Codec) GetName(n int) string {
	name, _ := num.Linux_x86_64[n]
	return name
}

func (c *Codec) decode(n int, args []uint64, ret uint64, done bool) (Syscall, error) {
	name, ok := num.Linux_x86_64[n]
	if !ok {
		return nil, fmt.Errorf("unknown syscall: %d\n", n)
	}
	if !done {
		return nil, errors.New("decoding unfinished syscalls is unimplemented")
	}
	base := call.Generic{n, name, args, ret}
	var out Syscall = &base
	var err error
	switch name {
	case "open":
		path, err := c.Mem.ReadStrAt(args[0])
		if err != nil {
			return nil, err
		}
		out = &call.Open{base, path, int(args[1]), int(args[2]), int(ret)}
	case "close":
		out = &call.Close{base, int(args[0])}
	case "read":
		length := int(int64(ret))
		var data []byte
		if length > 0 {
			data = make([]byte, ret)
			_, err = c.Mem.ReadAt(data, args[1])
		}
		out = &call.Read{base, int(args[0]), data, args[1], args[2], length}
	case "readv":
		length := int(int64(ret))
		var data []byte
		if length > 0 {
			data = make([]byte, length)
			mem := c.Mem.StreamAt(args[1])
			// TODO: platform specific
			var pos uint64
			for _, vec := range iovecRead(mem, args[2], 64, binary.LittleEndian) {
				end := vec.Len
				if int(pos+end) > length {
					end = uint64(length)
				}
				c.Mem.ReadAt(data[pos:end], vec.Base)
				if end == uint64(length) {
					break
				}
			}
		}
		out = &call.Read{base, int(args[0]), data, args[1], args[2], length}
	case "write":
		data := make([]byte, args[2])
		_, err = c.Mem.ReadAt(data, args[1])
		out = &call.Write{base, int(args[0]), data, args[1], args[2], int(int64(ret))}
	case "writev":
		mem := c.Mem.StreamAt(args[1])
		vecs := iovecRead(mem, args[2], 64, binary.LittleEndian)
		var size uint64
		for _, v := range vecs {
			size += v.Len
		}
		data := make([]byte, 0, size)
		for _, vec := range vecs {
			pos := uint64(len(data))
			data = data[:pos+vec.Len]
			c.Mem.ReadAt(data[pos:pos+vec.Len], vec.Base)
		}
		out = &call.Write{base, int(args[0]), data, args[1], args[2], int(int64(ret))}
	case "execve":
		path, _ := c.Mem.ReadStrAt(args[0])
		var readPointers = func(addr uint64) []uint64 {
			var pointers []uint64
			var tmp [8]byte
			stream := c.Mem.StreamAt(addr)
			for {
				_, err := stream.Read(tmp[:])
				if err != nil {
					break
				}
				ptr := binary.LittleEndian.Uint64(tmp[:])
				if ptr == 0 {
					break
				}
				pointers = append(pointers, ptr)
			}
			return pointers
		}
		argvAddrs := readPointers(args[1])
		argv := make([]string, len(argvAddrs))
		for i, addr := range argvAddrs {
			argv[i], _ = c.Mem.ReadStrAt(addr)
		}
		envpAddrs := readPointers(args[2])
		envp := make([]string, len(envpAddrs))
		for i, addr := range envpAddrs {
			envp[i], _ = c.Mem.ReadStrAt(addr)
		}
		out = &call.Execve{base, path, argv, envp}
	}
	return out, err
}
