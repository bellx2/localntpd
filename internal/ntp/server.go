package ntp

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

// NTPエポック (1900-01-01) とUnixエポック (1970-01-01) の差（秒）
const epochOffset = 2208988800

// Server はローカル時刻を返すNTPサーバー
type Server struct {
	addr    string
	stratum byte
}

// NewServer はNTPサーバーを生成する
func NewServer(addr string, stratum byte) *Server {
	return &Server{addr: addr, stratum: stratum}
}

// Run はUDPでNTPリクエストを受け付ける
func (s *Server) Run(ctx context.Context) error {
	conn, err := net.ListenPacket("udp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to start UDP listener: %w", err)
	}
	defer conn.Close()

	log.Printf("NTP server started: %s (stratum=%d)", s.addr, s.stratum)

	buf := make([]byte, 48)
	for {
		if ctx.Err() != nil {
			return nil
		}

		_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				continue
			}
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("read error: %v", err)
			continue
		}

		recvTime := time.Now()
		resp := buildResponse(buf[:n], recvTime, time.Now(), s.stratum)

		if _, err := conn.WriteTo(resp, addr); err != nil {
			log.Printf("write error (%s): %v", addr, err)
		}
	}
}

// buildResponse はNTP応答パケット（48バイト）を組み立てる
func buildResponse(req []byte, recvTime, xmitTime time.Time, stratum byte) []byte {
	resp := make([]byte, 48)

	// LI=0, VN=4, Mode=4 (server)
	resp[0] = 0x24
	resp[1] = stratum
	resp[2] = 6  // poll interval (2^6 = 64秒)
	resp[3] = 0xEC // precision

	// root dispersion: 約15.6ms
	binary.BigEndian.PutUint32(resp[8:12], 1024)

	// reference ID: ローカルクロック
	copy(resp[12:16], []byte("LOCL"))

	setTimestamp(resp[16:24], xmitTime)

	// origin timestamp: クライアントのtransmit timestampをエコー
	if len(req) >= 48 {
		copy(resp[24:32], req[40:48])
	}

	setTimestamp(resp[32:40], recvTime)
	setTimestamp(resp[40:48], xmitTime)

	return resp
}

// setTimestamp はNTPタイムスタンプ形式で8バイト書き込む
func setTimestamp(b []byte, t time.Time) {
	sec := uint32(t.Unix()) + epochOffset
	frac := uint32((uint64(t.Nanosecond()) << 32) / 1e9)
	binary.BigEndian.PutUint32(b[0:4], sec)
	binary.BigEndian.PutUint32(b[4:8], frac)
}
