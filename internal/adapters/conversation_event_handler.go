package adapters

import (
	"context"
	"log"
	"pm2/internal/domain/events"
	"pm2/internal/ports"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/ledongthuc/goterators"
	"github.com/spf13/viper"
)

func NewConversationEventConsumer(v *viper.Viper,
	kcm *kafka.ConfigMap) ports.Consumer[events.ConversationEvent] {
	return NewKafkaConsumer[events.ConversationEvent](kcm, &KafkaConsumerSettings{
		Topic:             v.GetString("kafka.conversationEvent.topic"),
		NumPartitions:     v.GetInt("kafka.conversationEvent.numpartitions"),
		ReplicationFactor: v.GetInt("kafka.conversationEvent.replicationfactor"),
		GroupId:           uuid.NewString(),
		AutoOffsetReset:   v.GetString("kafka.conversationEvent.auto-offset-reset"),
		Retries:           5,
	})
}

func NewConversationEventProducer(v *viper.Viper,
	kcm *kafka.ConfigMap) ports.Producer[events.ConversationEvent] {
	return NewKafkaProducer[events.ConversationEvent](kcm, &KafkaProducerSettings{
		Topic:             v.GetString("kafka.conversationEvent.topic"),
		NumPartitions:     v.GetInt("kafka.conversationEvent.numpartitions"),
		ReplicationFactor: v.GetInt("kafka.conversationEvent.replicationfactor"),
		Partition:         kafka.PartitionAny,
		Offset:            kafka.OffsetBeginning,
	})
}

func ConversationEventHandler(sm ports.SessionManagerUseCase) ports.ConsumerFunc[events.ConversationEvent] {
	return func(ctx context.Context, _ ports.ConsumerContext, cc2 events.ConversationEvent) {
		log.Printf("[CONVERSATION-EVENT-CONSUMER] - consume new conversation event: %s", cc2.Conversation.Id)
		if err := sm.InvokeSessionEvents(func(s *ports.Session) bool {
			_, _, err := goterators.Find(s.Keys, func(item string) bool {
				return item == cc2.Conversation.Tenant.Id.String()
			})
			return err == nil
		}, ctx, cc2.Conversation); err != nil {
			panic(err)
		}
		log.Printf("[CONVERSATION-EVENT-CONSUMER] - finish publishing conversation event: %s", cc2.Conversation.Id)
	}
}
