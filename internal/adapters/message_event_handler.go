package adapters

import (
	"context"
	"log"
	"pm2/internal/domain/events"
	"pm2/internal/ports"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func NewMessageEventConsumer(v *viper.Viper,
	kcm *kafka.ConfigMap) ports.Consumer[events.MessageEvent] {
	return NewKafkaConsumer[events.MessageEvent](kcm, &KafkaConsumerSettings{
		Topic:             v.GetString("kafka.messageEvent.topic"),
		NumPartitions:     v.GetInt("kafka.messageEvent.numpartitions"),
		ReplicationFactor: v.GetInt("kafka.messageEvent.replicationfactor"),
		GroupId:           uuid.NewString(),
		AutoOffsetReset:   v.GetString("kafka.messageEvent.auto-offset-reset"),
		Retries:           5,
	})
}

func NewMessageEventProducer(v *viper.Viper,
	kcm *kafka.ConfigMap) ports.Producer[events.MessageEvent] {
	return NewKafkaProducer[events.MessageEvent](kcm, &KafkaProducerSettings{
		Topic:             v.GetString("kafka.messageEvent.topic"),
		NumPartitions:     v.GetInt("kafka.messageEvent.numpartitions"),
		ReplicationFactor: v.GetInt("kafka.messageEvent.replicationfactor"),
		Partition:         kafka.PartitionAny,
		Offset:            kafka.OffsetBeginning,
	})
}

func NewMessageEventHandler(sm ports.SessionManagerUseCase) ports.ConsumerFunc[events.MessageEvent] {
	return func(ctx context.Context, _ ports.ConsumerContext, me events.MessageEvent) {
		log.Printf("[MESSAGE-EVENT-CONSUMER] - consume new message event for conversation: %s", me.ConversationId)
		if err := sm.InvokeSessionEvents(func(s *ports.Session) bool {
			return s.Key == me.ConversationId.String()
		}, ctx, me.Message); err != nil {
			panic(err)
		}
		log.Printf("[MESSAGE-EVENT-CONSUMER] - finish publishing message event for conversation: %s", me.ConversationId)
	}
}
