package rest

import (
	"net/http"
)

func (s *server) routes() {
	routes := []struct {
		method      string
		pattern     string
		handlerFunc http.HandlerFunc
	}{
		{
			"GET",
			"/categories",
			s.handleCategoriesGet(),
		},
		{
			"GET",
			"/categories/:name",
			s.handleCategoryGet(),
		},
		{
			"POST",
			"/categories",
			s.handleCategoryCreate(),
		},
		{
			"PUT",
			"/categories/:name",
			s.handleCategoryUpdate(),
		},
		{
			"DELETE",
			"/categories/:name",
			s.handleCategoryDelete(),
		},
		{
			"GET",
			"/genres",
			s.handleGenresGet(),
		},
		{
			"GET",
			"/genres/:name",
			s.handleGenreGet(),
		},
		{
			"POST",
			"/genres",
			s.handleGenreCreate(),
		},
		{
			"PUT",
			"/genres/:name",
			s.handleGenreUpdate(),
		},
		{
			"DELETE",
			"/genres/:name",
			s.handleGenreDelete(),
		},
		{
			"GET",
			"/cast_members",
			s.handleCastMembersGet(),
		},
		{
			"GET",
			"/cast_members/:name",
			s.handleCastMemberGet(),
		},
		{
			"POST",
			"/cast_members",
			s.handleCastMemberCreate(),
		},
		{
			"PUT",
			"/cast_members/:name",
			s.handleCastMemberUpdate(),
		},
		{
			"DELETE",
			"/cast_members/:name",
			s.handleCastMemberDelete(),
		},
	}

	for _, route := range routes {
		s.router.HandlerFunc(route.method, route.pattern, route.handlerFunc)
	}
}
