package mc

import (
	"encoding/json"
	"fmt"
	"os"
)

type StatusResponse struct {
	Version            Version     `json:"version"`
	Players            Players     `json:"players"`
	Description        Description `json:"description"`
	Favicon            string      `json:"favicon"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}

type Description struct {
	Text   string `json:"text"`
	Online int    `json:"online"`
}

func (s *StatusResponse) AsPackage() *Package {
	pkg := NewPackage(0x00)

	data, err := json.Marshal(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err.Error())
		return pkg
	}

	pkg.AddStringBytes(data)

	return pkg
}
