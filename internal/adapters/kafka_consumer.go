package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pm2/internal/ports"

	"reflect"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/exp/slices"
)

type (
	KafkaConsumerSettings struct {
		Topic             string
		NumPartitions     int
		ReplicationFactor int
		GroupId           string
		AutoOffsetReset   string
		Retries           uint16
	}

	KafkaConsumerRegister struct {
		kcs         *KafkaConsumerSettings
		topic       string
		consumerref interface{}
		fn          interface{}
		cfn         interface{}
	}

	KafkaConsumer[T interface{}] struct {
		Kcs *KafkaConsumerSettings
		kcm *kafka.ConfigMap
	}
)

func NewKafkaConsumer[T interface{}](kcm *kafka.ConfigMap,
	kcs *KafkaConsumerSettings) ports.Consumer[T] {
	cmt := *kcm
	cmt.SetKey("group.id", kcs.GroupId)
	cmt.SetKey("auto.offset.reset", kcs.AutoOffsetReset)
	kc := &KafkaConsumer[T]{
		Kcs: kcs,
		kcm: &cmt,
	}
	return kc
}

var kcms []KafkaConsumerRegister

func RegisterConsumer[T interface{}](cc ports.Consumer[T], fn ports.ConsumerFunc[T]) {
	kfc := cc.(*KafkaConsumer[T])
	kcms = append(kcms, KafkaConsumerRegister{
		kcs:         kfc.Kcs,
		topic:       kfc.Kcs.Topic,
		consumerref: kfc,
		fn:          fn,
		cfn:         KafkaCallFnWithResilence[T],
	})
}

func InitializeConsumers(kcm *kafka.ConfigMap) {
	tps := []string{}
	for _, v := range kcms {
		CreateKafkaTopic(context.Background(), kcm, &TopicConfiguration{
			Topic:             v.topic,
			NumPartitions:     v.kcs.NumPartitions,
			ReplicationFactor: v.kcs.ReplicationFactor,
		})
		tps = append(tps, v.topic)
	}
	c, err := kafka.NewConsumer(kcm)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
			InitializeConsumers(kcm)
		}
	}()
	c.SubscribeTopics(tps, rebalanceCallback)
	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			panic(err)
		}
		cm := kcms[slices.IndexFunc(kcms, func(k KafkaConsumerRegister) bool {
			return k.topic == *msg.TopicPartition.Topic
		})]
		cfn := reflect.ValueOf(cm.cfn) //KafkaCallFnWithResilence
		cfn.Call([]reflect.Value{
			reflect.ValueOf(msg),
			reflect.ValueOf(kcm),
			reflect.ValueOf(cm.consumerref).Elem().FieldByName("Kcs").Elem(),
			reflect.ValueOf(cm.fn),
		})
		c.CommitMessage(msg)
	}
}

func rebalanceCallback(c *kafka.Consumer, event kafka.Event) error {
	switch ev := event.(type) {
	case kafka.AssignedPartitions:
		log.Printf("%% %s rebalance: %d new partition(s) assigned: %v\n", c.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)
		err := c.IncrementalAssign(ev.Partitions)
		if err != nil {
			panic(err)
		}
	case kafka.RevokedPartitions:
		log.Printf("%% %s rebalance: %d partition(s) revoked: %v\n", c.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)
		if c.AssignmentLost() {
			fmt.Fprintf(os.Stderr, "%% Current assignment lost!\n")
		}
	}
	return nil
}

func KafkaCallFnWithResilence[T interface{}](msg *kafka.Message,
	kcm *kafka.ConfigMap,
	kcs KafkaConsumerSettings,
	fn ports.ConsumerFunc[T]) {

	ctx := context.Background()
	cctx := ports.ConsumerContext{RemainingRetries: kcs.Retries, Faulted: kcs.Retries == 0}
	var payload T

	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		kafkaSendToDlq(ctx, &kcs, kcm, msg, err)
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			fmt.Println(err.Error())
			if kcs.Retries > 1 {
				kcs.Retries--
				KafkaCallFnWithResilence(msg, kcm, kcs, fn)
				return
			}

			kafkaSendToDlq(ctx, &kcs, kcm, msg, err)
		}
	}()

	fn(ctx, cctx, payload)
}

func kafkaSendToDlq(
	ctx context.Context,
	kcs *KafkaConsumerSettings,
	kcm *kafka.ConfigMap,
	msg *kafka.Message,
	merr error) {
	p, err := kafka.NewProducer(kcm)

	if err != nil {
		panic(err)
	}

	defer p.Close()

	tpn := *msg.TopicPartition.Topic + "_error"
	msg.Headers = append(msg.Headers, kafka.Header{
		Key:   "error",
		Value: []byte(merr.Error()),
	})

	CreateKafkaTopic(ctx, kcm, &TopicConfiguration{
		Topic:             tpn,
		NumPartitions:     1,
		ReplicationFactor: kcs.ReplicationFactor,
	})

	dy := make(chan kafka.Event)
	msg.TopicPartition.Topic = &tpn
	msg.TopicPartition.Partition = kafka.PartitionAny
	if err = p.Produce(msg, dy); err != nil {
		panic(err)
	}
	<-dy
}
