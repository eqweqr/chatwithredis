package service

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// делаем с помощью inmemory db, в дальнейшем код будет перенесен в redis.

var (
	ErrUserMessageQueueClosed = errors.New("channel is closed")
	ErrRoomAlreadyExists      = errors.New("room with this name already exists")
	ErrRoomDoesntExists       = errors.New("room with this name doesnt exists")
	ErrInvalidRoomPassword    = errors.New("invalid password")
	ErrUserAlreadyInRoom      = errors.New("user already in room")
	ErrUserNotFound           = errors.New("user with follow name doesnt exists")
	ErrNotOwner               = errors.New("user not have permission")
)

type Message struct {
	Sender string
	Body   string
}

type RoomService interface {
	RoomExists(roomId string) bool
	CheckPasswor(roomId string, password string) error
	CreateRoom(name, password, owner string) error
	DeleteRoom(roomId string, userName string) error
	AddUserToRoom(roomId string, username string) error
	RemoveUserFromRoom(roomId string, username string) error
	GetRoomUsers(roomId string) ([]string, error)
	IsUserInRoom(roomId string, username string) (bool, error)
	BroadcastMessageToRoom(roomId string, message *Message) error
	GetUserMessage(roomId string, username string) (*Message, error)
}

type RoomServiceImpl struct {
	mu                  sync.RWMutex
	rooms               map[string]*room
	maxMessageQueueSize int
}

type room struct {
	roomId   string
	name     string
	password string
	owner    string
	users    map[string]*user
}

type user struct {
	name         string
	messageQueue chan *Message
}

func NewRoomService(maxMessageQueueSize int) *RoomServiceImpl {
	return &RoomServiceImpl{
		rooms:               make(map[string]*room),
		maxMessageQueueSize: maxMessageQueueSize,
	}
}

func (service *RoomServiceImpl) RoomExists(roomId string) bool {
	service.mu.RLock()
	_, ok := service.rooms[roomId]
	service.mu.RUnlock()

	return ok
}

func (service *RoomServiceImpl) CheckPasswor(roomId string, password string) error {
	service.mu.RLock()
	defer service.mu.RUnlock()
	if room, ok := service.rooms[roomId]; ok {
		if room.password == password {
			return nil
		}
		return ErrInvalidRoomPassword
	}
	return ErrRoomDoesntExists
}

func (service *RoomServiceImpl) CreateRoom(name, password, owner string) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	roomId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	if _, ok := service.rooms[roomId.String()]; ok {
		return ErrRoomAlreadyExists
	}

	service.rooms[roomId.String()] = &room{
		roomId:   roomId.String(),
		name:     name,
		password: password,
		owner:    owner,
		users:    make(map[string]*user),
	}

	return nil

}

func (service *RoomServiceImpl) DeleteRoom(roomId string, userName string) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	room, ok := service.rooms[roomId]
	if !ok {
		return ErrRoomDoesntExists
	}
	if userName != room.owner {
		return ErrNotOwner
	}

	for _, u := range room.users {
		close(u.messageQueue)
	}
	delete(service.rooms, roomId)

	return nil
}

func (service *RoomServiceImpl) AddUserToRoom(roomId string, username string) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	room, ok := service.rooms[roomId]
	if !ok {
		return ErrRoomDoesntExists
	}

	if _, ok := room.users[username]; ok {
		return ErrUserAlreadyInRoom
	}

	room.users[username] = &user{name: username, messageQueue: make(chan *Message, service.maxMessageQueueSize)}
	return nil

}

func (service *RoomServiceImpl) RemoveUserFromRoom(roomId string, username string) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	room, ok := service.rooms[roomId]
	if !ok {
		return ErrRoomDoesntExists
	}

	user, ok := room.users[username]
	if !ok {
		return ErrUserDosentExists
	}

	close(user.messageQueue)
	delete(room.users, username)
	return nil
}
func (service *RoomServiceImpl) GetRoomUsers(roomId string) ([]string, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	room, ok := service.rooms[roomId]
	if !ok {
		return nil, ErrRoomDoesntExists
	}
	var usersNames []string
	for _, user := range room.users {
		usersNames = append(usersNames, user.name)
	}

	return usersNames, nil
}
func (service *RoomServiceImpl) IsUserInRoom(roomId string, username string) (bool, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	room, ok := service.rooms[roomId]
	if !ok {
		return false, ErrRoomDoesntExists
	}

	_, ok = room.users[username]
	return ok, nil
}

func (service *RoomServiceImpl) BroadcastMessageToRoom(roomId string, message *Message) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	room, ok := service.rooms[roomId]
	if !ok {
		return ErrRoomDoesntExists
	}

	for _, user := range room.users {
		if user.name != message.Sender {
			user.messageQueue <- message
		}
	}

	return nil
}
func (service *RoomServiceImpl) GetUserMessage(roomId string, username string) (*Message, error) {
	service.mu.RLock()

	room, ok := service.rooms[roomId]
	if !ok {
		service.mu.RUnlock()
		return nil, ErrRoomDoesntExists
	}

	user, ok := room.users[username]
	if !ok {
		service.mu.RUnlock()
		return nil, ErrUserDosentExists
	}
	service.mu.RUnlock()

	msg, ok := <-user.messageQueue
	if !ok {
		return nil, ErrUserMessageQueueClosed
	}

	return msg, nil
}
