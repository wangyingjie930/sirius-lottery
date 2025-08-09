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
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"sirius-lottery/internal/infrastructure/contract/eventbus"
	"time"
)

type consumerImpl struct {
	broker        string
	topic         string
	groupID       string
	handler       eventbus.ConsumerHandler
	consumerGroup sarama.ConsumerGroup
}

func RegisterConsumer(broker string, topic, groupID string, handler eventbus.ConsumerHandler, opts ...eventbus.ConsumerOpt) error {
	logger := watermill.NewStdLogger(false, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return err
	}

	router.AddPlugin(plugin.SignalsHandler)

	// Router level middleware are executed for every message sent to the router
	router.AddMiddleware(
		// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
		middleware.CorrelationID,

		// The handler function is retried if it returns an error.
		// After MaxRetries, the message is Nacked and it's up to the PubSub to resend it.
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,

		// Recoverer handles panics from handlers.
		// In this case, it passes them as errors to the Retry middleware.
		middleware.Recoverer,
	)

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // Start consuming from the earliest message
	config.Consumer.Group.Session.Timeout = 30 * time.Second

	subscriber, _ := kafka.NewSubscriber(kafka.SubscriberConfig{
		Brokers:               []string{broker},
		Unmarshaler:           kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: config,
		ConsumerGroup:         groupID,
	}, logger)

	ctx := context.Background()
	router.AddNoPublisherHandler("test", topic, subscriber, func(msg *message.Message) error {
		return handler.HandleMessage(ctx, &eventbus.Message{
			Topic: topic,
			Group: groupID,
			Body:  msg.Payload,
		})
	})

	go func() {
		router.Run(ctx)
	}()

	return nil

	//
	//o := &eventbus.ConsumerOption{}
	//for _, opt := range opts {
	//	opt(o)
	//}
	//// TODO: orderly
	//
	//consumerGroup, err := sarama.NewConsumerGroup([]string{broker}, groupID, config)
	//if err != nil {
	//	return err
	//}
	//
	//c := &consumerImpl{
	//	broker:        broker,
	//	topic:         topic,
	//	groupID:       groupID,
	//	handler:       handler,
	//	consumerGroup: consumerGroup,
	//}
	//
	//ctx := context.Background()
	//go func(ctx context.Context) {
	//	for {
	//		if err := consumerGroup.Consume(ctx, []string{topic}, c); err != nil {
	//			logger.Ctx(ctx).Err(err).Msgf("consumer group consume: %v", err)
	//			break
	//		}
	//	}
	//}(ctx)
	//
	//go func(ctx context.Context) {
	//	waitExit()
	//
	//	if err := c.consumerGroup.Close(); err != nil {
	//		logger.Ctx(ctx).Err(err).Msgf("consumer group close: %v", err)
	//	}
	//}(ctx)
	//
	//return nil
}

func (c *consumerImpl) Setup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerImpl) Cleanup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerImpl) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.Background()

	for msg := range claim.Messages() {
		m := &eventbus.Message{
			Topic: msg.Topic,
			Group: c.groupID,
			Body:  msg.Value,
		}
		if err := c.handler.HandleMessage(ctx, m); err != nil {
			continue
		}

		sess.MarkMessage(msg, "") // TODO: Consumer policies can be configured
	}
	return nil
}
