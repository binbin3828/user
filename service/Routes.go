package service

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(s *Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/user/{uid}", s.responseHandler(s.GetUser))
	r.Post("/user", s.responseHandler(s.CreateUser))
	r.Put("/user", s.responseHandler(s.ModifyUser))
	r.Delete("/user/{uid}", s.responseHandler(s.DeleteUser))

	r.Post("/friends", s.responseHandler(s.AddFriend))
	r.Get("/friends/{uid}", s.responseHandler(s.GetFriendsList))
	r.Get("/nearbyfriends/{uid}", s.responseHandler(s.GetNearbyFriend))

	return r
}
