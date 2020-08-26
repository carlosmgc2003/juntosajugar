package models

import (
	"encoding/json"
	"errors"
	"regexp"
)

var (
	InvalidUsername = errors.New("User Model: Invalid Username")
	InvalidEmail    = errors.New("User Model: Invalid Email")
	InvalidFilename = errors.New("User Model: Invalid Filename")
)

func (U *User) FromJson(requestBody []byte) error {
	var temp_user User
	err := json.Unmarshal(requestBody, &temp_user)
	if err != nil {
		return err
	}
	U.ID = temp_user.ID
	if !validUsername(temp_user.Name) {
		return InvalidUsername
	}
	U.Name = temp_user.Name
	if !validEmail(temp_user.Email) {
		return InvalidEmail
	}
	U.Email = temp_user.Email
	if !validFilename(temp_user.Display_pic) {
		return InvalidFilename
	}
	U.Display_pic = temp_user.Display_pic
	return err
}

/*
Only contains alphanumeric characters, underscore and dot.
Underscore and dot can't be at the end or start of a username (e.g _username / username_ / .username / username.).
Underscore and dot can't be next to each other (e.g user_.name).
Underscore or dot can't be used multiple times in a row (e.g user__name / user..name).
*/

func validUsername(username string) bool {
	var re = regexp.MustCompile(`(?m)^(?:[a-zA-Z0-9]+[._]?[a-zA-Z0-9]+)+$`)
	return re.MatchString(username) && len(username) <= 30 && len(username) >= 6
}

func validEmail(email string) bool {
	var re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email) && len(email) <= 100
}

func validFilename(filename string) bool {
	return len(filename) >= 10 && len(filename) <= 50
}
