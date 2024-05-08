package main

import (
	"log"

	"google.golang.org/grpc"
)

type EdgeServerStream struct {
	grpc.ServerStream
}

var (
	serviceLoad = 0
)

func (e *EdgeServerStream) RecvMsg(m interface{}) error {
	// should do some calculation per stream here
	if err := e.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

func StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		serviceLoad++
		log.Printf("Service load: %d\n", serviceLoad)
		return handler(srv, ss)
	}
}
