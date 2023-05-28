package mc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

type Package struct {
	ID   int
	Data []byte
}

func NewPackage(ID int) *Package {
	return &Package{
		ID:   ID,
		Data: make([]byte, 0),
	}
}

func (pkg *Package) AddString(str string) {
	data := []byte(str)
	len := ToVarInt(len(data))
	pkg.Data = append(pkg.Data, len...)
	pkg.Data = append(pkg.Data, data...)
}

func (pkg *Package) AddStringBytes(str []byte) {
	len := ToVarInt(len(str))
	pkg.Data = append(pkg.Data, len...)
	pkg.Data = append(pkg.Data, str...)
}

func (pkg *Package) Read(cr Reader) error {

	pkgLen, _, err := ReadVarInt(cr)
	if err != nil {
		return err
	}
	pkgID, length, err := ReadVarInt(cr)
	if err != nil {
		return err
	}

	data, err := ReadData(cr, pkgLen-length)
	if err != nil {
		return err
	}

	pkg.ID = pkgID
	pkg.Data = data

	return nil
}

func (pkg *Package) WriteTo(w io.Writer) (int64, error) {
	written := int64(0)
	pkgID := ToVarInt(pkg.ID)

	pkgLen := ToVarInt(len(pkgID) + len(pkg.Data))

	length, err := w.Write(pkgLen)
	written += int64(length)
	if err != nil {
		return written, err
	}

	length, err = w.Write(pkgID)
	written += int64(length)
	if err != nil {
		return written, err
	}

	length, err = w.Write(pkg.Data)
	written += int64(length)
	if err != nil {
		return written, err
	}

	return written, nil
}

type InboundHandshake struct {
	ProtocolVersion int
	ServerAddress   string
	ServerPort      int
	NextState       int
}

func (pkg *Package) Handshake() (InboundHandshake, error) {
	hs := InboundHandshake{}
	br := bufio.NewReader(bytes.NewReader(pkg.Data))

	protVersion, _, err := ReadVarInt(br)
	if err != nil {
		return hs, err
	}
	hs.ProtocolVersion = protVersion

	strLen, _, err := ReadVarInt(br)
	if err != nil {
		return hs, err
	}
	addr, _ := ReadData(br, strLen)
	hs.ServerAddress = string(addr)

	portData, _ := ReadData(br, 2)
	port := binary.BigEndian.Uint16(portData)
	hs.ServerPort = int(port)

	nextState, _, err := ReadVarInt(br)
	if err != nil {
		return hs, err
	}
	hs.NextState = nextState

	return hs, nil
}
