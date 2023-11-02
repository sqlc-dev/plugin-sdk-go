package rpc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sqlc-dev/sqlc-go/plugin"
	pb "github.com/sqlc-dev/sqlc-go/plugin"
)

func Handle(server pb.CodegenServiceServer) {
	if err := handle(server); err != nil {
		fmt.Fprintf(os.Stderr, "error generating output: %s", err)
		os.Exit(2)
	}
}

func handle(server pb.CodegenServiceServer) error {
	handler := newStdioRPCHandler()
	pb.RegisterCodegenServiceServer(handler, server)
	return handler.Handle()
}

type stdioRPCHandler struct {
	services map[string]*serviceInfo
}

func newStdioRPCHandler() *stdioRPCHandler {
	return &stdioRPCHandler{services: map[string]*serviceInfo{}}
}

type serviceInfo struct {
	serviceImpl any
	methods     map[string]*grpc.MethodDesc
}

func (s *stdioRPCHandler) RegisterService(sd *grpc.ServiceDesc, ss any) {
	// TODO some type checking, see e.g. grpc server.RegisterService()
	info := &serviceInfo{
		serviceImpl: ss,
		methods:     make(map[string]*grpc.MethodDesc),
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}
	s.services[sd.ServiceName] = info
}

func (s *stdioRPCHandler) Handle() error {
	var methodArg string
	if len(os.Args) < 2 {
		// For backwards compatibility with sqlc before v1.24.0
		methodArg = fmt.Sprintf("/%s/%s", pb.CodegenService_ServiceDesc.ServiceName, "Generate")
	} else {
		methodArg = os.Args[1]
	}

	// Adapted from grpc server handleStream()

	sm := methodArg
	if sm != "" && sm[0] == '/' {
		sm = sm[1:]
	}
	pos := strings.LastIndex(sm, "/")
	if pos == -1 {
		errDesc := fmt.Sprintf("malformed method name: %q", methodArg)
		return status.Error(codes.Unimplemented, errDesc)
	}
	service := sm[:pos]
	method := sm[pos+1:]

	srv, knownService := s.services[service]
	if knownService {
		if md, ok := srv.methods[method]; ok {
			return s.processUnaryRPC(srv, md)
		}
	}

	// Unknown service, or known server unknown method.
	var errDesc string
	if !knownService {
		errDesc = fmt.Sprintf("unknown service %v", service)
	} else {
		errDesc = fmt.Sprintf("unknown method %v for service %v", method, service)
	}

	return status.Error(codes.Unimplemented, errDesc)
}

func (s *stdioRPCHandler) processUnaryRPC(srv *serviceInfo, md *grpc.MethodDesc) error {
	reqBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var resp protoreflect.ProtoMessage

	// TODO make this generic
	switch md.MethodName {
	case "Generate":
		var req plugin.GenerateRequest
		if err := proto.Unmarshal(reqBytes, &req); err != nil {
			return err
		}
		service, ok := srv.serviceImpl.(pb.CodegenServiceServer)
		if !ok {
			return status.Errorf(codes.Internal, codes.Internal.String())
		}
		resp, err = service.Generate(context.Background(), &req)
		if err != nil {
			return err
		}
	}

	respBytes, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(os.Stdout)
	if _, err := w.Write(respBytes); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
