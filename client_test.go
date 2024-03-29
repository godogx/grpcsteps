package grpcsteps_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/godogx/grpcsteps"
	"github.com/godogx/grpcsteps/internal/grpctest"
	testSrv "github.com/godogx/grpcsteps/internal/test/grpctest"
)

func TestClient_NoServer(t *testing.T) {
	t.Parallel()

	c := grpcsteps.NewClient(
		grpcsteps.RegisterServiceFromInstance("NoServer",
			(*grpctest.ItemServiceServer)(nil),
			grpcsteps.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		),
	)

	runClientSuite(t, c, "features/client/NoServer.feature")
}

func TestClient_GetItem(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		handler  func(ctx context.Context, request *grpctest.GetItemRequest) (*grpctest.Item, error)
	}{
		{
			scenario: "ServerPanic",
			handler: func(ctx context.Context, request *grpctest.GetItemRequest) (*grpctest.Item, error) {
				panic("Internal Server Error")
			},
		},
		{
			scenario: "InternalServerError",
			handler: func(ctx context.Context, request *grpctest.GetItemRequest) (*grpctest.Item, error) {
				return nil, status.Errorf(codes.Internal, "Internal Server Error")
			},
		},
		{
			scenario: "MethodUnimplemented",
			handler: func(ctx context.Context, request *grpctest.GetItemRequest) (*grpctest.Item, error) {
				return nil, status.Errorf(codes.Unimplemented, "GetItem is not implemented")
			},
		},
		{
			scenario: "NotFound",
			handler: func(ctx context.Context, request *grpctest.GetItemRequest) (*grpctest.Item, error) {
				return nil, status.Errorf(codes.NotFound, "Item %d not found", request.GetId())
			},
		},
		{
			scenario: "ErrorDocString",
			handler: func(ctx context.Context, request *grpctest.GetItemRequest) (*grpctest.Item, error) {
				return nil, status.Errorf(codes.FailedPrecondition, `invalid "id"`)
			},
		},
		{
			scenario: "Success",
			handler: func(ctx context.Context, request *grpctest.GetItemRequest) (*grpctest.Item, error) {
				var locale string

				if md, ok := metadata.FromIncomingContext(ctx); ok {
					if locales := md.Get("Locale"); len(locales) > 0 {
						locale = locales[0]
					}
				}

				return &grpctest.Item{
					Id:         42,
					Locale:     locale,
					Name:       "Test",
					CreateTime: timestamppb.Now(),
				}, nil
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			runClientTest(t,
				fmt.Sprintf("GetItem%s", tc.scenario),
				testSrv.GetItem(tc.handler),
			)
		})
	}
}

func TestClient_ListItems(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		handler  func(*grpctest.ListItemsRequest, grpctest.ItemService_ListItemsServer) error
	}{
		{
			scenario: "ServerPanic",
			handler: func(*grpctest.ListItemsRequest, grpctest.ItemService_ListItemsServer) error {
				panic("Internal Server Error")
			},
		},
		{
			scenario: "InternalServerError",
			handler: func(*grpctest.ListItemsRequest, grpctest.ItemService_ListItemsServer) error {
				return status.Errorf(codes.Internal, "Internal Server Error")
			},
		},
		{
			scenario: "MethodUnimplemented",
			handler: func(*grpctest.ListItemsRequest, grpctest.ItemService_ListItemsServer) error {
				return status.Errorf(codes.Unimplemented, "ListItems is not implemented")
			},
		},
		{
			scenario: "ErrorDocString",
			handler: func(*grpctest.ListItemsRequest, grpctest.ItemService_ListItemsServer) error {
				return status.Errorf(codes.FailedPrecondition, `invalid "page_size"`)
			},
		},
		{
			scenario: "Success",
			handler: func(_ *grpctest.ListItemsRequest, srv grpctest.ItemService_ListItemsServer) error {
				var locale string

				if md, ok := metadata.FromIncomingContext(srv.Context()); ok {
					if locales := md.Get("Locale"); len(locales) > 0 {
						locale = locales[0]
					}
				}

				item := &grpctest.Item{
					Id:         42,
					Locale:     locale,
					Name:       "Test",
					CreateTime: timestamppb.Now(),
				}

				return srv.Send(item)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			runClientTest(t,
				fmt.Sprintf("ListItems%s", tc.scenario),
				testSrv.ListItems(tc.handler),
			)
		})
	}
}

func TestClient_CreateItems(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		handler  func(itemsServer grpctest.ItemService_CreateItemsServer) error
	}{
		{
			scenario: "ServerPanic",
			handler: func(itemsServer grpctest.ItemService_CreateItemsServer) error {
				panic("Internal Server Error")
			},
		},
		{
			scenario: "InternalServerError",
			handler: func(itemsServer grpctest.ItemService_CreateItemsServer) error {
				return status.Errorf(codes.Internal, "Internal Server Error")
			},
		},
		{
			scenario: "MethodUnimplemented",
			handler: func(itemsServer grpctest.ItemService_CreateItemsServer) error {
				return status.Errorf(codes.Unimplemented, "CreateItems is not implemented")
			},
		},
		{
			scenario: "ErrorDocString",
			handler: func(itemsServer grpctest.ItemService_CreateItemsServer) error {
				return status.Errorf(codes.FailedPrecondition, `invalid "name"`)
			},
		},
		{
			scenario: "Success",
			handler: func(srv grpctest.ItemService_CreateItemsServer) error {
				var numItems int64

				for {
					_, err := srv.Recv()

					if errors.Is(err, io.EOF) {
						break
					}

					if err != nil {
						return err
					}

					numItems++
				}

				return srv.SendAndClose(&grpctest.CreateItemsResponse{NumItems: numItems})
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			runClientTest(t,
				fmt.Sprintf("CreateItems%s", tc.scenario),
				testSrv.CreateItems(tc.handler),
			)
		})
	}
}

func TestClient_TransformItems(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		handler  func(itemsServer grpctest.ItemService_TransformItemsServer) error
	}{
		{
			scenario: "ServerPanic",
			handler: func(itemsServer grpctest.ItemService_TransformItemsServer) error {
				panic("Internal Server Error")
			},
		},
		{
			scenario: "InternalServerError",
			handler: func(itemsServer grpctest.ItemService_TransformItemsServer) error {
				return status.Errorf(codes.Internal, "Internal Server Error")
			},
		},
		{
			scenario: "MethodUnimplemented",
			handler: func(itemsServer grpctest.ItemService_TransformItemsServer) error {
				return status.Errorf(codes.Unimplemented, "TransformItems is not implemented")
			},
		},
		{
			scenario: "ErrorDocString",
			handler: func(itemsServer grpctest.ItemService_TransformItemsServer) error {
				return status.Errorf(codes.FailedPrecondition, `invalid "name"`)
			},
		},
		{
			scenario: "Success",
			handler: func(srv grpctest.ItemService_TransformItemsServer) error {
				for {
					item, err := srv.Recv()

					if errors.Is(err, io.EOF) {
						break
					}

					if err != nil {
						return err
					}

					item.Name = fmt.Sprintf("Modified %s", item.GetName())

					if err := srv.Send(item); err != nil {
						return err
					}
				}

				return nil
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			runClientTest(t,
				fmt.Sprintf("TransformItems%s", tc.scenario),
				testSrv.TransformItems(tc.handler),
			)
		})
	}
}

func runClientTest(t *testing.T, scenario string, opts ...testSrv.ServiceOption) {
	t.Helper()

	dialer := testSrv.StartServer(t, opts...)

	c := grpcsteps.NewClient(
		grpcsteps.WithDefaultServiceOptions(
			grpcsteps.WithDialOptions(
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(dialer),
			),
		),
		grpcsteps.RegisterService(grpctest.RegisterItemServiceServer),
	)

	runClientSuite(t, c, fmt.Sprintf("features/client/%s.feature", scenario))
}

func runClientSuite(t suiteT, c *grpcsteps.Client, paths ...string) {
	runSuite(t,
		initScenario(c.RegisterContext),
		featureFiles(paths...),
	)
}
