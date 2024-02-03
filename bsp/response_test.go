package bsp

import (
	"reflect"
	"testing"
)

func TestNewInfo(t *testing.T) {
	type args struct {
		i    ReplyType
		info [][]byte
	}
	tests := []struct {
		name string
		args args
		want *InfoReply
	}{
		{
			name: "test1",
			args: args{
				i:    ReplyInfo,
				info: [][]byte{[]byte("test")},
			},
			want: &InfoReply{
				i:    ReplyInfo,
				info: []byte("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInfo(tt.args.i, tt.args.info...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfoReply_Bytes(t *testing.T) {
	type fields struct {
		i    ReplyType
		info []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "test1",
			fields: fields{
				i:    ReplyInfo,
				info: []byte("test"),
			},
			want: []byte{byte(ReplyInfo), 0x74, 0x65, 0x73, 0x74, 0x0a},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := InfoReply{
				i:    tt.fields.i,
				info: tt.fields.info,
			}
			if got := i.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumResp_Bytes(t *testing.T) {
	type fields struct {
		n   ReplyType
		num []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "test1",
			fields: fields{
				n:   ReplyNumber,
				num: []byte("test"),
			},
			want: []byte{byte(ReplyNumber), 0x74, 0x65, 0x73, 0x74, 0x0a},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NumResp{
				n:   tt.fields.n,
				num: tt.fields.num,
			}
			if got := n.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNum(t *testing.T) {
	type args struct {
		num any
	}
	tests := []struct {
		name string
		args args
		want *NumResp
	}{
		{
			name: "test1",
			args: args{
				num: []byte("test"),
			},
			want: &NumResp{
				n:   ReplyNumber,
				num: []byte("test"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNum(tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewList(t *testing.T) {
	type args struct {
		list [][]byte
	}
	tests := []struct {
		name string
		args args
		want *ListResp
	}{
		{
			name: "test1",
			args: args{
				list: [][]byte{[]byte("test1"), []byte("test2"), []byte("test3")},
			},
			want: &ListResp{
				l:    ReplyList,
				list: [][]byte{[]byte("test1"), []byte("test2"), []byte("test3")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewList(tt.args.list...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}
		})
	}
}
