package system

import (
	"os/user"
	"strings"
)

type Service struct{}

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
