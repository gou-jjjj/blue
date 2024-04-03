package network

import "testing"

func TestParseAddr(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{"127.0.0.1:8080"}, true},
		{"test2", args{"39.101.169.250:7893"}, true},
		{"test3", args{"127.0.0.1:80801"}, false},
		{"test4", args{"3123.312.2.2:8801"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseAddr(tt.args.addr); got != tt.want {
				t.Errorf("ParseAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
