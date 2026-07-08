package main

import (
	"testing"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name    string
		stratum uint
		wantErr bool
	}{
		{name: "最小値", stratum: 1, wantErr: false},
		{name: "デフォルト", stratum: 2, wantErr: false},
		{name: "最大値", stratum: 15, wantErr: false},
		{name: "下限未満", stratum: 0, wantErr: true},
		{name: "上限超過", stratum: 16, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*stratum = tt.stratum
			err := validateFlags()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceArguments(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		stratum uint
		want    []string
	}{
		{name: "デフォルト", addr: ":123", stratum: 2, want: []string{"run"}},
		{name: "addrのみ", addr: ":12345", stratum: 2, want: []string{"run", "-addr", ":12345"}},
		{name: "stratumのみ", addr: ":123", stratum: 1, want: []string{"run", "-stratum", "1"}},
		{name: "両方", addr: "0.0.0.0:123", stratum: 3, want: []string{"run", "-addr", "0.0.0.0:123", "-stratum", "3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*addr = tt.addr
			*stratum = tt.stratum
			got := serviceArguments()
			if len(got) != len(tt.want) {
				t.Fatalf("serviceArguments() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("serviceArguments() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
