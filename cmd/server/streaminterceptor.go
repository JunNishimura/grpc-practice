package main

import (
	"log"

	"google.golang.org/grpc"
)

func myStreamServerInterceptor1(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre] my stream server interceptor 1: ", info.FullMethod)

	err := handler(srv, &myServerStreamWrapper{ss})

	log.Println("[post] my stream server interceptor 1: ", info.FullMethod)
	return err
}

type myServerStreamWrapper struct {
	grpc.ServerStream
}

func (w *myServerStreamWrapper) RecvMsg(m any) error {
	err := w.ServerStream.RecvMsg(m)
	log.Println("[after recv] my stream server interceptor 1: ", m)
	return err
}

func (w *myServerStreamWrapper) SendMsg(m any) error {
	log.Println("[before send] my stream server interceptor 1: ", m)
	return w.ServerStream.SendMsg(m)
}
