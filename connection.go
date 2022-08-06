package quickfix

import (
	"context"
	"io"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)

var (
	outgoingMessageCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "quickfix",
		Name:      "outgoing_message_count",
		Help:      "Count of messages sent inside quickfix writeLoop function",
	}, []string{"sessionID"})

	outgoingMessageBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "quickfix",
		Name:      "ougoing_message_bytes",
		Help:      "Count of bytes sent inside quickfix writeLoop function",
	}, []string{"sessionID"})

	incomingMessageCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "quickfix",
		Name:      "incoming_message_count",
		Help:      "Count of messages received inside quickfix writeLoop function",
	}, []string{"sessionID"})

	incomingMessageBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "quickfix",
		Name:      "incoming_message_bytes",
		Help:      "Count of bytes received inside quickfix readLoop function",
	}, []string{"sessionID"})
)

func init() {
	prometheus.MustRegister(outgoingMessageCount)
	prometheus.MustRegister(outgoingMessageBytes)
	prometheus.MustRegister(incomingMessageCount)
	prometheus.MustRegister(incomingMessageBytes)
}

func writeLoop(connection io.Writer, messageOut chan []byte, log Log, maxMessagesPerSecond int, sessionID string) {
	limiter := rate.NewLimiter(rate.Limit(maxMessagesPerSecond), 1)

	for {
		msg, ok := <-messageOut
		if !ok {
			return
		}

		if maxMessagesPerSecond > 0 {
			err := limiter.Wait(context.Background())
			if err != nil {
				log.OnEvent(err.Error())
				continue
			}
		}

		outgoingMessageCount.WithLabelValues(sessionID).Inc()
		outgoingMessageBytes.WithLabelValues(sessionID).Add(float64(len(msg)))

		if _, err := connection.Write(msg); err != nil {
			log.OnEvent(err.Error())
		}
	}
}

func readLoop(parser *parser, msgIn chan fixIn, sessionID string) {
	defer close(msgIn)

	for {

		msg, err := parser.ReadMessage()
		if err != nil {
			return
		}

		incomingMessageCount.WithLabelValues(sessionID).Inc()
		incomingMessageBytes.WithLabelValues(sessionID).Add(float64(msg.Len()))

		msgIn <- fixIn{msg, parser.lastRead}
	}
}
