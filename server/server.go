package server

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"

	"github.com/ethaniccc/estatz/packet"
	"github.com/rs/zerolog"
)

// PacketHandler is a function that is able to process a packet from the given client.
type PacketHandler func(sender *net.UDPAddr, header *packet.PacketHeader, pk packet.Packet) error

// The server is responsible for interpreting and storing all information given from the clients.
// It should only recieve data, and never should have to send a response back to any sender.
type Server struct {
	// logger is the logger used for the server.
	logger *zerolog.Logger
	// conn is the underlying connection from where the Server reads incoming packets from.
	conn *net.UDPConn
	// msgQueue is the queue for messages sent by clients to the server. This queue is handled by various
	// workers.
	msgQueue chan *packet.Message

	packetHandlers []PacketHandler
}

func New(cfg Config) *Server {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: cfg.Port})
	if err != nil {
		panic(err)
	}

	return &Server{
		conn:           conn,
		msgQueue:       make(chan *packet.Message, 65535),
		packetHandlers: make([]PacketHandler, 0),
	}
}

func (srv *Server) worker(id int) {
	defer func() {
		if v := recover(); v != nil {
			srv.logger.Err(fmt.Errorf("%v", v)).Int("workerID", id).Msg("worker crashed")
			go srv.worker(id)
		}
	}()
	srv.logger.Debug().Int("worker", id).Msg("server worker started")

	for msg := range srv.msgQueue {
		header, pk, ok := msg.Decode()
		if !ok {
			srv.logger.Warn().
				Str("addr", msg.Sender().String()).
				Uint64("packetID", header.PacketID).
				Str("JWT", base64.StdEncoding.EncodeToString(header.JWT)).
				Msg("unable to find packet")
			msg.Dispose()
			return
		}

		for _, handler := range srv.packetHandlers {
			handler(msg.Sender(), header, pk)
		}
		msg.Dispose()
	}
}

func (srv *Server) listen() {
	msgBuffer := make([]byte, 1492)
	for {
		size, senderAddr, err := srv.conn.ReadFromUDP(msgBuffer)
		if err != nil {
			srv.logger.Err(err).Str("addr", senderAddr.String()).Msg("failed to read message")
			continue
		}

		select {
		case srv.msgQueue <- packet.NewMessage(msgBuffer[:size], senderAddr):
			// OK
		case <-time.After(time.Second):
			srv.logger.Warn().Str("addr", senderAddr.String()).Msg("failed to push message into queue")
		}
	}
}
