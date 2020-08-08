package rest

//func TestCategoryCtx(t *testing.T) {
//	type args struct {
//		s crud.Service
//	}
//	tests := []struct {
//		name string
//		args args
//		want func(next http.Handler) http.Handler
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := CategoryCtx(tt.args.s); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("CategoryCtx() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestGetCategory(t *testing.T) {
//	type args struct {
//		w http.ResponseWriter
//		r *http.Request
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//		})
//	}
//}

//func Test_server_handleCategoriesGet(t *testing.T) {
//	type fields struct {
//		router *chi.Mux
//		svc      crud.Service
//		m      modifying.Service
//		log    zerolog.Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   http.HandlerFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &server{
//				router: tt.fields.router,
//				svc:      tt.fields.svc,
//				m:      tt.fields.m,
//				log:    tt.fields.log,
//			}
//			if got := s.handleCategoriesGet(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("handleCategoriesGet() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_server_handleCategoryCreate(t *testing.T) {
//	type fields struct {
//		router *chi.Mux
//		svc      crud.Service
//		m      modifying.Service
//		log    zerolog.Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   http.HandlerFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &server{
//				router: tt.fields.router,
//				svc:      tt.fields.svc,
//				m:      tt.fields.m,
//				log:    tt.fields.log,
//			}
//			if got := s.handleCategoryCreate(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("handleCategoryCreate() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_server_handleCategoryDelete(t *testing.T) {
//	type fields struct {
//		router *chi.Mux
//		svc      crud.Service
//		m      modifying.Service
//		log    zerolog.Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   http.HandlerFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &server{
//				router: tt.fields.router,
//				svc:      tt.fields.svc,
//				m:      tt.fields.m,
//				log:    tt.fields.log,
//			}
//			if got := s.handleCategoryDelete(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("handleCategoryDelete() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_server_handleCategoryUpdate(t *testing.T) {
//	type fields struct {
//		router *chi.Mux
//		svc      crud.Service
//		m      modifying.Service
//		log    zerolog.Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   http.HandlerFunc
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &server{
//				router: tt.fields.router,
//				svc:      tt.fields.svc,
//				m:      tt.fields.m,
//				log:    tt.fields.log,
//			}
//			if got := s.handleCategoryUpdate(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("handleCategoryUpdate() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
