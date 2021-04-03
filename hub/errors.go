package hub

import "errors"

var (
	ErrNotInHub = errors.New("hub: user not in hub")
	ErrDuplicate = errors.New("hub: user already exists")
	ErrPrivateHub = errors.New("hub: hub is private")
)

func returnError(c chan error, err error) {
	c <- err
	close(c)
}
