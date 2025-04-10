package server

type Server struct {
}

func Serve(port int) (*Server, error) {
	return &Server{}, nil
}

func (server *Server) Close() {

}
