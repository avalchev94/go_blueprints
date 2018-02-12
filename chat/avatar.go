package main

import (
	"errors"
	"io/ioutil"
	"path"
)

// ErrNoAvatarURL is the error when the Avatar instance is not able
// to provide an avatar.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatarURL from the client, or returns error
	// if something goes wrong. Returns ErrNoAvatarURL on error.
	GetAvatarURL(c *client) (string, error)
}

// Gravatar Avatar
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userid"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return "//www.gravatar.com/avatar/" + userIdStr, nil
		}
	}
	return "", nil
}

// File System Avatar
type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userid"]; ok {
		if userIdStr, ok := userId.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					if match, _ := path.Match(userIdStr+"*", file.Name()); match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
