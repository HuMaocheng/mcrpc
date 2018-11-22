package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"os"
)

type Server struct {
	network string
	addr string
	funcs map[string]reflect.Value

	m_listener net.Listener
	running bool

	logger *log.Logger
}

func NewServer() *Server{
	m_server := new(Server)

	//init log
	logFile, err := os.OpenFile(Local_Log_addr, os.O_CREATE | os.O_WRONLY | os.O_APPEND,0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：",err)
	}
	m_server.logger = log.New(io.MultiWriter(os.Stderr, logFile),"Info:",log.Ldate | log.Ltime | log.Lshortfile)

	//init net
	m_server.network = network
	m_server.addr = Local_TCP_addr

	//init func map
	m_server.funcs = make(map[string]reflect.Value)

	return m_server
}

func (m_server *Server) Start() (error) {
	//start listener
	var err error
	m_server.m_listener, err = net.Listen(m_server.network, m_server.addr)
	if err != nil {
		m_server.logger.Println("[Error] server listener failed to start")
		os.Exit(0)
	}
	m_server.logger.Println("[Info] server listener started successfully")
	m_server.running = true

	//listener accept
	for m_server.running {
		conn, err := m_server.m_listener.Accept()
		if err != nil {
			m_server.logger.Println("[Error] server accept error")
			continue
		}
		m_server.logger.Println("[Info] server got a new accept")
		go m_server.acceptor(conn)
	}
	return nil
}

func (m_server *Server) acceptor(conn net.Conn) {
	m_conn := new(connector)
	m_conn.m_connector = conn

	for {
		data, err := m_conn.ReadMessage()
		if err != nil {
			m_server.logger.Println("[Error] acceptor read error:", err.Error())
			return
		}
		data, err = m_server.worker(data)
		if err != nil {
			m_server.logger.Println("[Error] worker error:", err.Error())
			return
		}
		err = m_conn.WriteMessage(data)
		if err != nil {
			m_server.logger.Println("[Error] acceptor write error", err.Error())
			return
		}
	}
}

func (m_server *Server) worker(data []byte) ([]byte, error) {
	funcName, args, err := decodeData(data)
	if err != nil {
		return nil, err
	}

	m_func, ok := m_server.funcs[funcName]

	if !ok {
		return nil, fmt.Errorf("[Info] rpc function %s not registered", funcName)
	}

	funcArgs := make([]reflect.Value, len(args))

	for i := 0; i < len(args); i++ {
		funcArgs[i] = reflect.ValueOf(args[i])
	}

	funcOutput := m_func.Call(funcArgs)

	m_output := make([]interface{}, len(funcOutput))
	for i := 0; i < len(funcOutput); i++ {
		m_output[i] = funcOutput[i].Interface()
	}

	return encodeData(funcName, m_output)
}
func (m_server *Server) Stop() bool {
	if m_server.m_listener == nil {
		m_server.logger.Println("[Error] server failed to stop, listener is nil")
		return false
	}
	m_server.m_listener.Close()
	m_server.running = false
	return true
}

func (m_server *Server) Register(funcName string, newfunc interface{}) (error) {
	_, ok := m_server.funcs[funcName]
	if ok {
		m_server.logger.Printf("[Info] function %s has registered\n", funcName)
		return fmt.Errorf("[Info] function %s has registered\n", funcName)
	}
	m_server.funcs[funcName] = reflect.ValueOf(newfunc)
	return nil
}

