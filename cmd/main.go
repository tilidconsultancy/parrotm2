package main

import (
	"net"
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
		hh *adapters.WhatsAppHookHandler) {
		eng.GET("", hh.CheckHook)
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
	c.Provide(adapters.NewMongoDbRepository[domain.TenantUser])
	c.Provide(adapters.NewMongoDbRepository[domain.IncomingMessage])
	c.Provide(adapters.NewMongoDbRepository[domain.LabelMeaning])
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
	c.Provide(adapters.NewLabelMeaningService)

	c.Provide(adapters.NewgRPCConnection)
	c.Provide(gRPC.NewSpeechServiceStreamClient)
	c.Provide(adapters.NewChattyClient)

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
