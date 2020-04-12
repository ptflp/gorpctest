package server

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

type Server struct {
	Port     uint
	UseHttp  bool
	UseJson  bool
	Sleep    time.Duration
	listener net.Listener
}

func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) Start() error {
	if s.Port <= 0 {
		return errors.New("port must be specified")
	}

	err := rpc.Register(&Handler{
		Sleep: s.Sleep,
	})

	if err != nil {
		return nil
	}

	s.listener, err = net.Listen("tcp", ":"+strconv.Itoa(int(s.Port)))
	if err != nil {
		return err
	}

	if s.UseHttp {
		rpc.HandleHTTP()
		return http.Serve(s.listener, nil)
	}

	if s.UseJson {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				return err
			}

			jsonrpc.ServeConn(conn)
		}
	}

	rpc.Accept(s.listener)

	return nil
}
