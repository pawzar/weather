package net

import "net"

func Must(p int, e error) int {
	if e != nil {
		panic(e)
	}

	return p
}

func GetFreePort() (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	err = ln.Close()
	if err != nil {
		return 0, err
	}

	return ln.Addr().(*net.TCPAddr).Port, nil
}
