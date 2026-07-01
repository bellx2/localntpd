package ntp

import (
	"encoding/binary"
	"testing"
	"time"
)

func TestBuildResponse(t *testing.T) {
	req := make([]byte, 48)
	// クライアントのtransmit timestampを設定
	binary.BigEndian.PutUint32(req[40:44], 100)
	binary.BigEndian.PutUint32(req[44:48], 200)

	recv := time.Unix(1700000000, 500000000)
	xmit := time.Unix(1700000000, 600000000)

	resp := buildResponse(req, recv, xmit, 2)

	if resp[0] != 0x24 {
		t.Errorf("mode/version: got 0x%02x, want 0x24", resp[0])
	}
	if resp[1] != 2 {
		t.Errorf("stratum: got %d, want 2", resp[1])
	}
	if got := binary.BigEndian.Uint32(resp[24:28]); got != 100 {
		t.Errorf("origin sec: got %d, want 100", got)
	}
	if got := binary.BigEndian.Uint32(resp[24+4 : 28+4]); got != 200 {
		t.Errorf("origin frac: got %d, want 200", got)
	}

	sec, frac := readTimestamp(resp[40:48])
	wantSec := uint32(xmit.Unix()) + epochOffset
	wantFrac := uint32((uint64(xmit.Nanosecond()) << 32) / 1e9)
	if sec != wantSec || frac != wantFrac {
		t.Errorf("transmit: got (%d,%d), want (%d,%d)", sec, frac, wantSec, wantFrac)
	}
}

func readTimestamp(b []byte) (uint32, uint32) {
	return binary.BigEndian.Uint32(b[0:4]), binary.BigEndian.Uint32(b[4:8])
}
