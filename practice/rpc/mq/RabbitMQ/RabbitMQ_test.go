package RabbitMQ

import (
	"testing"

	"github.com/streadway/amqp"
)

func TestRabbitMQ_SendMsg(t *testing.T) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
