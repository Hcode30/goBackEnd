package main

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Password struct {
	salt []byte
	hash []byte
}

type User struct {
	// todo: add friends,avatar implementation
	ID        uint64
	UserName  string
	Password  Password
	Email     string
	FirstName string
	LastName  string
	Friends   []uint64
	Avatar    string
	Role      string
	Messages  []Message
}

func MakeUser(email, firstName, lastName, password string) (*User, error) {
	id := rand.Uint64()
	// validate first firstName
	err := checkName(firstName)
	if err != nil {
		return nil, err
	}
	// validate last lastName
	err = checkName(lastName)
	if err != nil {
		return nil, err
	}
	// validate password	username := ""
	username := ""
	for _, char := range email {
		if char == '@' {
			break
		}
		username += string(char)
	}
	username = strings.ToLower(username)
	// validate email
	username, err = checkEmail(email, username)
	if err != nil {
		return nil, err
	}
	err = checkPassword(password, firstName, lastName, username)
	if err != nil {
		return nil, err
	}
	// hash password
	pass, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        id,
		UserName:  username,
		Password:  pass,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Friends:   make([]uint64, 0),
		Messages:  make([]Message, 0),
		Role:      "user",
	}, nil
}

func checkName(name string) error {
	if len(name) <= 2 {
		return errors.New("name must be at least 3 characters long")
	}
	if len(name) > 48 {
		return errors.New("name must be at most 48 characters long")
	}
	if checkSpecialChars(name) {
		return errors.New("name must not contain special characters")
	}
	return nil
}

func checkSpecialChars(str string) bool {
	specialChars := "!@#$%^&*()_+{}[]|:;<>,.?/~`"
	specialCount := 0
	for _, char := range specialChars {
		if strings.Contains(str, string(char)) {
			specialCount++
		}
	}
	if specialCount == 0 {
		return false
	}
	return true
}

func usernameDoesExist(username string) bool {
	/*for _, user := range users {
	    if user.UserName == username {
	      return true
	    }
	  }
	*/
	// for lsp satisfaction
	username = username + ""
	return false
}
func emailDoesExist(email string) bool {
	/*for _, user := range
	  users {
	    if user.Email == email {
	      return true
	    }
	  }
	*/
	// for lsp satisfaction
	email = email + ""
	return false
}

func checkEmail(email, username string) (string, error) {
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matches := pattern.MatchString(email)
	if !matches {
		return "", errors.New("invalid email address")
	}
	if emailDoesExist(email) {
		return "", errors.New("email already exists")
	}
	for usernameDoesExist(username) {
		username = username + strconv.Itoa(rand.Intn(10))
	}
	return username, nil
}

func checkPassword(password, firstName, lastName, username string) error {
	lowerPassword := strings.ToLower(password)
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if password == lowerPassword {
		return errors.New("password must contain at least one uppercase letter")
	}

	if password == strings.ToUpper(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if _, err := strconv.Atoi(password); err == nil {
		return errors.New("password must contain at least one letter")
	}
	if !checkSpecialChars(password) {
		return errors.New("password must contain at least one special character")
	}
	firstName = strings.ToLower(firstName)
	lastName = strings.ToLower(lastName)

	if strings.Contains(lowerPassword, firstName) ||
		strings.Contains(lowerPassword, lastName) {
		return errors.New("password must not contain first or last name")
	}

	if strings.Contains(lowerPassword, username) {
		return errors.New("password must not contain username")
	}
	return nil
}

func hashPassword(password string) (Password, error) {
	hsh := Argon2idHash{
		time:    3,
		saltLen: 16,
		memory:  12288,
		threads: 1,
		keyLen:  32,
	}
	salt, err := randomSecret(hsh.saltLen)
	if err != nil {
		return Password{}, err
	}
	hashSalt, err := hsh.GenerateHash([]byte(password), salt)
	if err != nil {
		return Password{}, err
	}
	password = string(hashSalt.Hash)
	return Password{
		salt: salt,
		hash: hashSalt.Hash,
	}, nil
}


func (u *User) AddFriend(friend *User) {
  u.Friends = append(u.Friends, friend.ID)
}


func (u *User) SendMessage(receiver *User, content string) {
  message := Message{
    ID:         rand.Uint64(),
    SenderID:   u.ID,
    ReceiverID: receiver.ID,
    Content:    content,
    Time:       time.Now(),
  }
  u.Messages = append(u.Messages, message)
  receiver.Messages = append(receiver.Messages, message)
}


func (u *User) DeleteFriend(friend *User) {
  for i, id := range u.Friends {
    if id == friend.ID {
      u.Friends = append(u.Friends[:i], u.Friends[i+1:]...)
      break
    }
  }
}
func GetUser(email string) (*User, error) {
  for _, user := range users {
    if user.Email == email {
      return &user, nil
    }
  }
  return nil, errors.New("user not found")
}
