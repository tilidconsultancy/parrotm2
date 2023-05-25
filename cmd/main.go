package main

import (
	"net"
	"net/http"
	"pm2/internal/adapters"
	"pm2/internal/adapters/gRPC"
	"pm2/internal/domain"
	"pm2/internal/domain/events"
	"pm2/internal/ports"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

func main() {
	c := buildContainer()
	if err := c.Invoke(func(v *viper.Viper,
		cs gRPC.ConversationServiceServer,
		ms gRPC.MessageServiceServer) {
		lis, err := net.Listen("tcp", v.GetString("grpcport"))
		if err != nil {
			panic(err)
		}
		gsrv := grpc.NewServer()
		gRPC.RegisterConversationServiceServer(gsrv, cs)
		gRPC.RegisterMessageServiceServer(gsrv, ms)
		go gsrv.Serve(lis)
	}); err != nil {
		panic(err)
	}

	if err := c.Invoke(func(cec ports.Consumer[events.ConversationEvent],
		mec ports.Consumer[events.MessageEvent],
		sm ports.SessionManagerUseCase,
		kcm *kafka.ConfigMap) {
		adapters.RegisterConsumer(cec, adapters.ConversationEventHandler(sm))
		adapters.RegisterConsumer(mec, adapters.NewMessageEventHandler(sm))
		go adapters.InitializeConsumers(kcm)
	}); err != nil {
		panic(err)
	}

	if err := c.Invoke(func(
		v *viper.Viper,
		eng *gin.Engine,
		hh *adapters.WhatsAppHookHandler,
		nc ports.NlpClient) {
		eng.GET("", hh.CheckHook)
		eng.GET("key/:k/org/:o", func(ctx *gin.Context) {
			gcli := nc.(*adapters.GptClient)
			key := ctx.Param("k")
			org := ctx.Param("o")
			gcli.Token = key
			if org != "same" {
				gcli.Org = org
			}
			ctx.Status(http.StatusOK)
		})
		eng.POST("", hh.IncomingMessage)
		eng.Run(v.GetString("ginport"))
	}); err != nil {
		panic(err)
	}
}

func buildContainer() *dig.Container {
	c := dig.New()
	c.Provide(gin.Default)
	c.Provide(initializeViper)
	c.Provide(adapters.NewClientOptions)
	c.Provide(adapters.NewMongoClient)
	c.Provide(adapters.NewMongoDatabase)
	c.Provide(adapters.NewMongoDbRepository[domain.Conversation])
	c.Provide(adapters.NewMongoDbRepository[domain.Tenant])
	c.Provide(adapters.NewGptClient)
	c.Provide(adapters.NewMetaSettings)
	c.Provide(adapters.NewMetaHttpClient)
	c.Provide(adapters.NewConversationService)
	c.Provide(adapters.NewMessageServer)
	c.Provide(adapters.NewWhatsAppHookHandler)
	c.Provide(adapters.NewSessionManagerService)
	c.Provide(adapters.NewConversationServer)
	c.Provide(adapters.NewKafkaConfigMap)
	c.Provide(adapters.NewKafkaAdminClient)
	c.Provide(adapters.NewConversationEventConsumer)
	c.Provide(adapters.NewConversationEventProducer)
	c.Provide(adapters.NewMessageEventConsumer)
	c.Provide(adapters.NewMessageEventProducer)

	return c
}

func initializeViper() *viper.Viper {
	v := viper.New()
	v.AddConfigPath("../configs")
	v.SetConfigType("json")
	v.SetConfigName("local")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return v
}
