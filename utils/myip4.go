package utils

import "net"

func MyIp4() (string, error) {

	addr, err := net.ResolveUDPAddr("udp4", "8.8.8.8:53")
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	host, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return "", err
	}

	return host, nil
}

func ListenUDP4() (*net.UDPConn, error) {
	localhost, err := MyIp4()
	if err != nil {
		return nil, err
	}

	localUDPAddr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(localhost, "0"))
	if err != nil {
		return nil, err
	}

	return net.ListenUDP("udp4", localUDPAddr)
}