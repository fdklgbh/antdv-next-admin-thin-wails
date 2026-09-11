package system

import (
	"os"
	"os/user"
	"runtime"
	"strings"
)

type Service struct{}

// IsUbuntu22 checks the running system, not the machine used to build the app.
func (*Service) IsUbuntu22() (bool, error) {
	if runtime.GOOS != "linux" {
		return false, nil
	}
	data, err := os.ReadFile("/etc/os-release")
	if os.IsNotExist(err) {
		data, err = os.ReadFile("/usr/lib/os-release")
	}
	if err != nil {
		return false, err
	}
	var id, version string
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		value = strings.Trim(value, "\"'")
		switch key {
		case "ID":
			id = value
		case "VERSION_ID":
			version = value
		}
	}
	return id == "ubuntu" && (version == "22" || strings.HasPrefix(version, "22.")), nil
}

func (*Service) CurrentUsername() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	username := currentUser.Username
	if separator := strings.LastIndexByte(username, '\\'); separator >= 0 {
		username = username[separator+1:]
	}

	return username, nil
}
