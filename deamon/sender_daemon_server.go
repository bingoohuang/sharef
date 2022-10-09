package deamon

import (
	context "context"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/emiraganov/sharef/api"
	"github.com/emiraganov/sharef/streamer"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type SenderDaemonServer struct {
	sender *streamer.Sender
}

func StartSenderDaemonServer(sender *streamer.Sender, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	api.RegisterSenderServer(grpcServer, &SenderDaemonServer{sender: sender})
	// determine whether to use TLS
	return grpcServer.Serve(lis)
}

func (*SenderDaemonServer) Hello(context.Context, *api.HelloRequest) (*api.HelloReply, error) {
	return &api.HelloReply{}, nil
}

func (s *SenderDaemonServer) SendFile(req *api.SendFileRequest, stream api.Sender_SendFileServer) error {
	fi, err := os.Stat(req.Filename)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	fs := s.sender.NewFileStreamer(req.Filename)
	reader, writer := io.Pipe()
	fs.SetOutput(writer)

	fs.AsyncStream(fi)
	for {
		select {
		case <-fs.Done:
			return nil
		default:
		}

		data := make([]byte, 4096)
		n, err := reader.Read(data)
		if err != nil && err != io.EOF {
			return status.Errorf(codes.Internal, err.Error())
		}

		out := &api.STDOutput{
			Line: string(data[:n]),
		}

		if err := stream.Send(out); err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
}
