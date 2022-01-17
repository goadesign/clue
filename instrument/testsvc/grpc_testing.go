package testsvc

import (
	"context"
	"net"
	"testing"

	"github.com/crossnokaye/micro/instrument/testsvc/gen/grpc/test/client"
	testpb "github.com/crossnokaye/micro/instrument/testsvc/gen/grpc/test/pb"
	"github.com/crossnokaye/micro/instrument/testsvc/gen/grpc/test/server"
	"github.com/crossnokaye/micro/instrument/testsvc/gen/test"
	"google.golang.org/grpc"
)

type (
	// GRPClient is a test service gRPC client.
	GRPClient interface {
		GRPCMethod(context.Context, *Fields) (*Fields, error)
		GRPCStream(context.Context) (Stream, error)
	}

	// GRPCOption is a function that can be used to configure the gRPC server.
	GRPCOption func(*grpcOptions)

	grpcc struct {
		genc *client.Client
	}

	grpcOptions struct {
		grpcfn   UnaryFunc
		streamfn StreamFunc
		unary    grpc.UnaryServerInterceptor
		stream   grpc.StreamServerInterceptor
	}
)

// SetupGRPC instantiates the test service with the given options and starts a
// gRPC server to host it. The results are a ready-to-use client and a function
// used to stop the server.
func SetupGRPC(t *testing.T, opts ...GRPCOption) (c GRPClient, stop func()) {
	t.Helper()

	// Configure options
	options := &grpcOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Create test gRPC server
	s := svc{grpcfn: options.grpcfn, streamfn: options.streamfn}
	endpoints := test.NewEndpoints(&s)
	svr := server.New(endpoints, nil, nil)
	server := grpc.NewServer(grpc.UnaryInterceptor(options.unary), grpc.StreamInterceptor(options.stream))
	testpb.RegisterTestServer(server, svr)

	// Start listen loop
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	go server.Serve(l)

	// Connect to in-memory server
	conn, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	c = grpcc{client.NewClient(conn)}

	// Cleanup
	stop = func() {
		server.Stop()
		conn.Close()
		l.Close()
	}

	return
}

// WithUnaryFunc provides the implementation for the gRPC unary method.
func WithUnaryFunc(fn UnaryFunc) GRPCOption {
	return func(opt *grpcOptions) {
		opt.grpcfn = fn
	}
}

// WithStreamFunc provides the implementation for the gRPC streaming method.
func WithStreamFunc(fn StreamFunc) GRPCOption {
	return func(opt *grpcOptions) {
		opt.streamfn = fn
	}
}

// WithUnaryInterceptor provides the unary interceptor for the gRPC server.
func WithUnaryInterceptor(fn grpc.UnaryServerInterceptor) GRPCOption {
	return func(opt *grpcOptions) {
		opt.unary = fn
	}
}

// WithStreamInterceptor provides the stream interceptor for the gRPC server.
func WithStreamInterceptor(fn grpc.StreamServerInterceptor) GRPCOption {
	return func(opt *grpcOptions) {
		opt.stream = fn
	}
}

// GRPCMethod implements the gRPC method.
func (c grpcc) GRPCMethod(ctx context.Context, req *Fields) (res *Fields, err error) {
	rq := &test.Fields{}
	if req != nil {
		*rq = test.Fields(*req)
	}
	resp, err := c.genc.GrpcMethod()(ctx, rq)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		res = &Fields{}
		*res = Fields(*(resp.(*test.Fields)))
	}

	return
}

// GRPCStream implements the gRPC stream method.
func (c grpcc) GRPCStream(ctx context.Context) (Stream, error) {
	res, err := c.genc.GrpcStream()(ctx, nil)
	if err != nil {
		return nil, err
	}
	return adapter{res.(test.GrpcStreamClientStream)}, nil
}
