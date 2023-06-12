package adapters

import (
	"bytes"
	"context"
	"io"
	"pm2/internal/adapters/gRPC"
	"pm2/internal/domain"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	ChattyClient struct {
		client gRPC.SpeechServiceStreamClient
	}
)

func NewChattyClient(client gRPC.SpeechServiceStreamClient) *ChattyClient {
	return &ChattyClient{
		client: client,
	}
}

func NewgRPCConnection(v *viper.Viper) grpc.ClientConnInterface {
	conn, err := grpc.Dial(v.GetString("chattyaddr"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return conn
}

func (c *ChattyClient) TextToSpeech(ctx context.Context, txt string, t *domain.Tenant) (io.Reader, error) {
	s, err := c.client.TextToSpeech(ctx, &gRPC.TextToSpeechRequest{
		Voice:        t.AccountSettings.Voice,
		Content:      txt,
		OutputFormat: gRPC.AudioFormat_OGG_OPUS,
	})
	if err != nil {
		return nil, err
	}
	buff := &bytes.Buffer{}
	for err != io.EOF {
		abuff, err := s.Recv()
		if err == io.EOF {
			return buff, nil
		}
		if err != nil {
			return nil, err
		}
		if _, err := buff.Write(abuff.Chunk); err != nil {
			return nil, err
		}
	}
	return buff, nil
}

func (c *ChattyClient) SpeechToText(ctx context.Context, filestream io.Reader) (string, error) {
	s, err := c.client.SpeechToText(context.Background())
	if err != nil {
		return "", err
	}
	s.Send(&gRPC.SpeechToTextRequest{
		Payload: &gRPC.SpeechToTextRequest_InputFormat{
			InputFormat: gRPC.AudioFormat_ANY,
		},
	})
	for {
		msg, err := s.Recv()
		if err != nil {
			return "", err
		}
		switch e := msg.Payload.(type) {
		case *gRPC.SpeechToTextResponse_Size:
			b := make([]byte, e.Size)
			n, err := filestream.Read(b)
			if err != nil && err != io.EOF {
				return "", err
			}
			s.Send(&gRPC.SpeechToTextRequest{
				Payload: &gRPC.SpeechToTextRequest_ChunkData{
					ChunkData: &gRPC.AudioBuffer{
						Chunk: b,
						Size:  int32(n),
					},
				},
			})
		case *gRPC.SpeechToTextResponse_Content:
			return e.Content, nil
		}
	}
}
