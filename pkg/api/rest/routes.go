package rest

//func (s *server) routesFc() {
//routes := []struct {
//	method  string
//	pattern string
//	//handlerFunc http.HandlerFunc
//	handler http.Handler
//}{
//	{
//		"GET",
//		"/categories",
//		s.handleCategoriesGet(),
//	},
//	{
//		"GET",
//		"/categories/:name",
//		s.handleCategoryGet(),
//	},
//	{
//		"POST",
//		"/categories",
//		createHandler,
//	},
//	{
//		"PUT",
//		"/categories/:name",
//		s.handleCategoryUpdate(),
//	},
//	{
//		"DELETE",
//		"/categories/:name",
//		s.handleCategoryDelete(),
//	},
//{
//	"GET",
//	"/genres",
//	s.handleGenresGet(),
//},
//{
//	"GET",
//	"/genres/:name",
//	s.handleGenreGet(),
//},
//{
//	"POST",
//	"/genres",
//	s.handleGenreCreate(),
//},
//{
//	"PUT",
//	"/genres/:name",
//	s.handleGenreUpdate(),
//},
//{
//	"DELETE",
//	"/genres/:name",
//	s.handleGenreDelete(),
//},
//{
//	"GET",
//	"/cast_members",
//	s.handleCastMembersGet(),
//},
//{
//	"GET",
//	"/cast_members/:name",
//	s.handleCastMemberGet(),
//},
//{
//	"POST",
//	"/cast_members",
//	s.handleCastMemberCreate(),
//},
//{
//	"PUT",
//	"/cast_members/:name",
//	s.handleCastMemberUpdate(),
//},
//{
//	"DELETE",
//	"/cast_members/:name",
//	s.handleCastMemberDelete(),
//},
//{
//	"GET",
//	"/videos",
//	s.handleVideosGet(),
//},
//{
//	"GET",
//	"/videos/:title",
//	s.handleVideoGet(),
//},
//{
//	"POST",
//	"/videos",
//	s.handleVideoCreate(),
//},
//{
//	"PUT",
//	"/videos/:title",
//	s.handleVideoUpdate(),
//},
//{
//	"DELETE",
//	"/videos/:title",
//	s.handleVideoDelete(),
//},
//}

//for _, route := range routes {
//	s.router.Handler(route.method, route.pattern, route.handler)
//	//s.router.HandlerFunc(route.method, route.pattern, route.handlerFunc)
//}
//}
