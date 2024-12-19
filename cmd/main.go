package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"time"
)

/**
*
* main
* <p>
* main file
*
* Copyright (c) 2024 All rights reserved.
*
* This source code is shared under a collaborative license.
* Contributions, suggestions, and improvements are welcome!
* Feel free to fork, modify, and submit pull requests under the terms of the repository's license.
* Please ensure proper attribution to the original author(s) and maintain this notice in derivative works.
*
* @author christian
* @author dbacilio88@outlook.es
* @since 18/12/2024
*
 */

var URI = "amqp://guest:guest@localhost:5672/"

const QU = "QU-CORE-CONSUMER-TRANSACTION-RESPONSE"
const RK = "SERVICE.CONNECTOR.MD-CORE.TRANSACTION.RESPONSE"

const RKP = "SERVICE.MD-CORE.CONNECTOR.TRANSACTION.REQUEST"
const EX = "topic.exchange.transaction"
const MESSAGE = "HELLO WORLD FROM THE CORE"

func main() {

	amqpConfig := amqp.Config{
		Connection: amqp.ConnectionConfig{
			AmqpURI:   URI,
			Reconnect: amqp.DefaultReconnectConfig(),
		},
		Marshaler: amqp.DefaultMarshaler{},
		Exchange: amqp.ExchangeConfig{
			GenerateName: func(topic string) string {
				fmt.Println("topic exchange ", topic)
				return EX
			},
			Type: "topic",
			//Durable:     true,
			//AutoDeleted: true,
			//NoWait:      true,
		},
		Queue: amqp.QueueConfig{
			//GenerateName: amqp.GenerateQueueNameTopicName,
			GenerateName: func(topic string) string {
				fmt.Println("topic queue ", topic)
				return QU
			},
			Durable: true,
		},
		QueueBind: amqp.QueueBindConfig{
			GenerateRoutingKey: func(topic string) string {
				fmt.Println("topic binding ", topic)
				return RK
			},
		},
		Publish: amqp.PublishConfig{
			GenerateRoutingKey: func(topic string) string {
				fmt.Println("topic pub ", topic)
				return topic
			},
		},
		Consume: amqp.ConsumeConfig{
			Qos: amqp.QosConfig{
				PrefetchCount: 1,
			},
		},
		TopologyBuilder: &amqp.DefaultTopologyBuilder{},
	}
	//amqp.NewDurableQueueConfig()
	//amqpConfig = amqp.NewDurableQueueConfig(amqpURI)

	subscriber, err := amqp.NewSubscriber(
		// This config is based on this example: https://www.rabbitmq.com/tutorials/tutorial-two-go.html
		// It works as a simple queue.
		//
		// If you want to implement a Pub/Sub style service instead, check
		// https://watermill.io/pubsubs/amqp/#amqp-consumer-groups
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("subscribing to messages")
	messages, err := subscriber.Subscribe(context.Background(), QU)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Subscribed to messages")
	go process(messages)

	publisher, err := amqp.NewPublisher(amqpConfig, watermill.NewStdLogger(true, true))
	if err != nil {
		log.Fatal(err)
	}

	publishMessages(publisher)
}

func publishMessages(publisher message.Publisher) {
	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte(MESSAGE))

		if err := publisher.Publish(RKP, msg); err != nil {
			log.Fatal(err)
		}
		log.Println("Published message to topic", watermill.NewUUID())
		time.Sleep(time.Second * 10)
	}
}

func process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
