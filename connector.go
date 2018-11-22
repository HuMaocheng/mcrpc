package main

import (
	"net"
)

type connector struct {
	m_connector net.Conn
}

func newConn(network, addr string) (*connector, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	m_conn := new(connector)
	m_conn.m_connector = conn
	return m_conn, nil
}

func (conn *connector) Close() (error){
	return conn.m_connector.Close()
}

func (conn *connector) WriteMessage(data []byte) (error) {

	_, err := conn.m_connector.Write(data)
	if err != nil {
		conn.Close()
		return err
	}
	return nil
}

func (conn *connector) ReadMessage() ([]byte, error) {
	data := make([]byte, 4)

	_, err := conn.m_connector.Read(data)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return data, nil
}
