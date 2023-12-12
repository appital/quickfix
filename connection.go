// Copyright (c) quickfixengine.org  All rights reserved.
//
// This file may be distributed under the terms of the quickfixengine.org
// license as defined by quickfixengine.org and appearing in the file
// LICENSE included in the packaging of this file.
//
// This file is provided AS IS with NO WARRANTY OF ANY KIND, INCLUDING
// THE WARRANTY OF DESIGN, MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE.
//
// See http://www.quickfixengine.org/LICENSE for licensing information.
//
// Contact ask@quickfixengine.org if any conditions of this licensing
// are not clear to you.

package quickfix

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

const maxBurst = 1

func writeLoop(connection io.Writer, messageOut chan []byte, log Log, maxMessagesPerSecond int) {
	limiter := rate.NewLimiter(rate.Limit(maxMessagesPerSecond), maxBurst)
	ctx := context.Background()

	for {
		msg, ok := <-messageOut
		if !ok {
			return
		}

		if maxMessagesPerSecond > 0 {
			err := limiter.Wait(ctx)
			// if we get an error from the limiter we log, but continue on to
			// attempt to write the message to the remote party
			if err != nil {
				log.OnEvent(err.Error())
			}
		}

		if _, err := connection.Write(msg); err != nil {
			log.OnEvent(err.Error())
		}
	}
}

func readLoop(parser *parser, msgIn chan fixIn) {
	defer close(msgIn)

	for {
		msg, err := parser.ReadMessage()
		if err != nil {
			return
		}
		msgIn <- fixIn{msg, parser.lastRead}
	}
}
