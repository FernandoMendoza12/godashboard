package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Model data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	//write log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Model.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	//everything comes correct
	return &logs.LogResponse{
		Result: "Logged",
	}, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to lsiten for gRPC :%v", err)
	}
	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Model: app.Models})

	log.Printf("gRPC server started on port: %s ", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for gRPC :%v", err)
	}

}
