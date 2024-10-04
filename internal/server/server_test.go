package server

import (
	"context"
	"testing"

	pb "github.com/FlutterDizaster/EncryNest/api/generated"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestServer_Register(t *testing.T) {
	type test struct {
		name     string
		username string
		password string
		email    string
		wantErr  bool
	}

	tests := []test{
		{
			name:     "success",
			username: "username",
			password: "password",
			email:    "user@example.com",
			wantErr:  false,
		},
		{
			name:     "same user registration",
			username: "username",
			password: "password",
			email:    "user@example.com",
			wantErr:  true,
		},
		{
			name:     "blank fields",
			username: "",
			password: "",
			email:    "",
			wantErr:  true,
		},
	}

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	settings := Settings{
		Addr: "",
		Port: "50555",
	}
	srv := NewServer(settings)

	go func() {
		serr := srv.Run(ctx)
		if serr != nil {
			t.Error(serr)
		}
	}()

	conn, cerr := grpc.NewClient(
		":50555",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, cerr)
	defer conn.Close()

	client := pb.NewEncryNestUserServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.RegisterUserRequest{
				Username: tt.username,
				Password: tt.password,
				Email:    tt.email,
			}

			resp, err := client.RegisterUser(ctx, req)
			if tt.wantErr {
				require.Error(t, err, "test name = %s", tt.name)
			} else {
				require.NoError(t, err, "test name = %s", tt.name)
				assert.NotEmpty(t, resp.GetToken(), "test name = %s", tt.name)
			}
		})
	}
}

func TestServer_Authenticate(t *testing.T) {
	type test struct {
		name             string
		username         string
		password         string
		wantErr          bool
		needRegistration bool
	}

	tests := []test{
		{
			name:             "success",
			username:         "username",
			password:         "password",
			wantErr:          false,
			needRegistration: true,
		},
		{
			name:             "unregistered user auth",
			username:         "username1",
			password:         "password1",
			wantErr:          true,
			needRegistration: false,
		},
	}

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	settings := Settings{
		Addr: "",
		Port: "50555",
	}
	srv := NewServer(settings)

	go func() {
		serr := srv.Run(ctx)
		if serr != nil {
			t.Error(serr)
		}
	}()

	conn, cerr := grpc.NewClient(
		":50555",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, cerr)
	defer conn.Close()

	client := pb.NewEncryNestUserServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.needRegistration {
				registerReq := &pb.RegisterUserRequest{
					Username: tt.username,
					Password: tt.password,
					Email:    tt.username,
				}

				_, err := client.RegisterUser(ctx, registerReq)
				require.NoError(t, err, "test name = %s", tt.name)
			}

			req := &pb.AuthenticateUserRequest{
				Username: tt.username,
				Password: tt.password,
			}

			resp, err := client.AuthenticateUser(ctx, req)
			if tt.wantErr {
				require.Error(t, err, "test name = %s", tt.name)
			} else {
				require.NoError(t, err, "test name = %s", tt.name)
				assert.NotEmpty(t, resp.GetToken(), "test name = %s", tt.name)
			}
		})
	}
}
