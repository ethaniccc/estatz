package server

import (
	"encoding/base64"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/ethaniccc/estatz/packet"
	"github.com/rs/zerolog"
)

const queueTimeout = time.Second * 10

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
	// running is an atomic boolean that is set to true once the server starts listening for connections.
	running atomic.Bool

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
		closeChan:      make(chan struct{}, 1),
		packetHandlers: make([]PacketHandler, 0),
	}
}

// AddHandler adds a packet handler to the list of the server's packet handlers.
func (srv *Server) AddHandler(handler PacketHandler) {
	if srv.running.Load() {
		panic(fmt.Errorf("cannot add packet handler when server is running"))
	}
	srv.packetHandlers = append(srv.packetHandlers, handler)
}

// Start starts the server and listens for connections on it's given connection.
func (srv *Server) Start() {
	if srv.running.Load() {
		panic("server is already started")
	}
	srv.running.Store(true)
	srv.listen()
}

// Stop stops the server and it's processing of packets.
func (srv *Server) Stop() {
	srv.running.Store(false)
	close(srv.msgQueue)
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
		srv.handleMessage(msg)
	}
}

func (srv *Server) handleMessage(msg *packet.Message) {
	defer func() {
		if v := recover(); v != nil {
			srv.logger.Err(fmt.Errorf("%v", v)).
				Str("sender", msg.Sender().String()).
				Msg("error occured when attempting to process message")
			// TODO: Sentry logging.
		}
	}()

	header, pk, ok := msg.Decode()
	if !ok {
		srv.logger.Warn().
			Str("addr", msg.Sender().String()).
			Uint64("packetID", header.PacketID).
			Str("JWT", base64.StdEncoding.EncodeToString(header.JWT)).
			Msg("unable to find packet with ID")
		msg.Dispose()
		return
	}

	for _, handler := range srv.packetHandlers {
		handler(msg.Sender(), header, pk)
	}
	msg.Dispose()
}

func (srv *Server) listen() {
	msgBuffer := make([]byte, 1492)
	for srv.running.Load() {
		size, senderAddr, err := srv.conn.ReadFromUDP(msgBuffer)
		if err != nil {
			srv.logger.Err(err).Str("addr", senderAddr.String()).Msg("failed to read message")
			continue
		}

		select {
		case srv.msgQueue <- packet.NewMessage(msgBuffer[:size], senderAddr):
			// OK
		case <-time.After(queueTimeout):
			srv.logger.Warn().Str("addr", senderAddr.String()).Msg("failed to push message into queue")
		}
	}
}
