package rabbitMQ

import (
	"fmt"
	"testing"

	"github.com/streadway/amqp"
)

func TestRabbitMQ_SendMsg(t *testing.T) {
	topic := NewRabbitMQTopic("exKutengTopic", "#")
	defer topic.channel.Close()
	type fields struct {
		conn      *amqp.Connection
		channel   *amqp.Channel
		QueueName string
		Exchange  string
		Key       string
		Mqurl     string
	}
	type args struct {
		topic string
		msg   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "1",
			fields: fields{
				conn:      topic.conn,
				channel:   topic.channel,
				QueueName: topic.QueueName,
				Exchange:  topic.Exchange,
				Key:       topic.Key,
				Mqurl:     topic.Mqurl,
			},
			args: args{
				topic: "topic",
				msg:   "222222222",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RabbitMQ{
				conn:      tt.fields.conn,
				channel:   tt.fields.channel,
				QueueName: tt.fields.QueueName,
				Exchange:  tt.fields.Exchange,
				Key:       tt.fields.Key,
				Mqurl:     tt.fields.Mqurl,
			}
			r.SendMsg(tt.args.topic, tt.args.msg)
		})
	}
}

func TestRabbitMQ_SubsribeMsg(t *testing.T) {
	topic := NewRabbitMQTopic("exKutengTopic", "122")
	defer topic.channel.Close()
	type fields struct {
		conn      *amqp.Connection
		channel   *amqp.Channel
		QueueName string
		Exchange  string
		Key       string
		Mqurl     string
	}
	type args struct {
		topic   string
		handler func(msg string)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "1",
			fields: fields{
				conn:      topic.conn,
				channel:   topic.channel,
				QueueName: topic.QueueName,
				Exchange:  topic.Exchange,
				Key:       topic.Key,
				Mqurl:     topic.Mqurl,
			},
			args: args{
				topic: "topic",
				handler: func(msg string) {
					fmt.Println(msg)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RabbitMQ{
				conn:      tt.fields.conn,
				channel:   tt.fields.channel,
				QueueName: tt.fields.QueueName,
				Exchange:  tt.fields.Exchange,
				Key:       tt.fields.Key,
				Mqurl:     tt.fields.Mqurl,
			}
			r.SubsribeMsg(tt.args.topic, tt.args.handler)
		})
	}
}
