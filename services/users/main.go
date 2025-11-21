package main

import (
	"context"
	"fmt"
	"sync"

	"elproxy.cloud/elproxy-server/common/model"
	"go-micro.dev/v5"
)

type UsersService struct {
	mu    sync.Mutex
	users map[string]*model.User
}

func (s *UsersService) Get(ctx context.Context, req *model.GetUserRequest, rsp *model.GetUserResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[req.ID]
	if !exists {
		return fmt.Errorf("user with ID %s not found", req.ID)
	}

	*rsp = model.GetUserResponse{User: *user}
	return nil
}

func main() {
	svc := micro.NewService(
		micro.Name("users"),
		micro.Version("1.0.0"),
	)
	svc.Init()

	usersService := newUsersService()
	svc.Handle(usersService)

	if err := svc.Run(); err != nil {
		fmt.Println("Error running service:", err)
	}
}

func newUsersService() *UsersService {
	users := []*model.User{
		{ID: "1", Name: "Alice", Email: "alice@user.service"},
		{ID: "2", Name: "Bob", Email: "bob@user.service"},
	}

	userMap := make(map[string]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	return &UsersService{
		users: userMap,
	}
}
