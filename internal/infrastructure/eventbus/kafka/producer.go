/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/wangyingjie930/nexus-pkg/logger"
	"os"
	"os/signal"
	"sirius-lottery/internal/infrastructure/contract/eventbus"
	"syscall"
)

type producerImpl struct {
	topic string
	p     *kafka.Publisher
}

func NewProducer(broker, topic string) (eventbus.Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := kafka.NewPublisher(kafka.PublisherConfig{
		Brokers:               []string{broker},
		Marshaler:             kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: config,
	}, watermill.NewStdLogger(false, false))

	if err != nil {
		return nil, err
	}

	go func(ctx context.Context) {
		waitExit()
		if err := producer.Close(); err != nil {
			logger.Ctx(ctx).Err(err).Msg("close producer error")
		}
	}(context.Background())

	return &producerImpl{
		topic: topic,
		p:     producer,
	}, nil
}

func (r *producerImpl) Send(ctx context.Context, body []byte, opts ...eventbus.SendOpt) error {
	return r.BatchSend(ctx, [][]byte{body}, opts...)
}

func (r *producerImpl) BatchSend(ctx context.Context, bodyArr [][]byte, opts ...eventbus.SendOpt) error {
	option := eventbus.SendOption{}
	for _, opt := range opts {
		opt(&option)
	}

	var msgArr []*message.Message
	for _, body := range bodyArr {
		msg := message.NewMessage(watermill.NewUUID(), body)
		msgArr = append(msgArr, msg)
	}

	err := r.p.Publish(r.topic, msgArr...)
	if err != nil {
		return err
	}

	return nil
}

func waitExit() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-signals
}
