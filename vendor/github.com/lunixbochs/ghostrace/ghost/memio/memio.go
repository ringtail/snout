package memio

import (
	"bytes"
	"io"
)

type MemIO interface {
	ReadAt(p []byte, addr uint64) (int, error)
	ReadStrAt(addr uint64) (string, error)
	WriteAt(p []byte, addr uint64) (int, error)
	StreamAt(addr uint64) io.ReadWriter
}

type callback func(p []byte, addr uint64) (int, error)

type memIO struct {
	read, write callback
}

func NewMemIO(read, write callback) MemIO {
	return &memIO{read, write}
}

func (m *memIO) ReadAt(p []byte, addr uint64) (int, error) {
	return m.read(p, addr)
}

// TODO: this probably makes way too many syscalls.
// could do a large page-aligned read but that's not as portable
func (m *memIO) ReadStrAt(addr uint64) (string, error) {
	var tmp = [4]byte{1, 1, 1, 1}
	var ret []byte
	nul := []byte{0}
	for !bytes.Contains(tmp[:], nul) {
		n, err := m.ReadAt(tmp[:], addr)
		if err != nil {
			return "", err
		}
		addr += uint64(n)
		ret = append(ret, tmp[:]...)
	}
	split := bytes.Index(ret, nul)
	return string(ret[:split]), nil
}

func (m *memIO) WriteAt(p []byte, addr uint64) (int, error) {
	return m.write(p, addr)
}

func (m *memIO) StreamAt(addr uint64) io.ReadWriter {
	return &memIOStream{*m, addr}
}

type memIOStream struct {
	memIO
	addr uint64
}

func (m *memIOStream) Read(p []byte) (int, error) {
	n, err := m.read(p, m.addr)
	m.addr += uint64(n)
	return n, err
}

func (m *memIOStream) Write(p []byte) (int, error) {
	n, err := m.write(p, m.addr)
	m.addr += uint64(n)
	return n, err
}
