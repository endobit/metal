package auth

import (
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username     string
	PasswordHash string
	Admin        bool
}

func newUser(username, password string, admin bool) (*user, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	user := &user{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Admin:        admin,
	}

	return user, nil
}

func (u *user) IsPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

type userStore struct {
	mutex sync.RWMutex
	users map[string]user
}

func newUserStore() *userStore {
	return &userStore{
		users: make(map[string]user),
	}
}

func (s *userStore) Save(user user) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.users[user.Username]; ok {
		return ErrUserExists
	}

	s.users[user.Username] = user

	return nil
}

func (s *userStore) Get(username string) (user, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	u, ok := s.users[username]
	if !ok {
		return user{}, ErrUserNotFound
	}

	return u, nil
}
