package hub

import "errors"

var (
	ErrNotInHub = errors.New("hub: user not in hub")
	ErrDuplicate = errors.New("hub: user already exists")
	ErrPrivateHub = errors.New("hub: hub is private")
	ErrNotAdmin = errors.New("hub: action requires admin privilege")
)

func returnError(c chan error, err error) {
	c <- err
	close(c)
}
