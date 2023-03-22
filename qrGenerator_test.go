package main

import (
	"image/color"
	"testing"
)

func TestBuildQRCode(t *testing.T) {
	type args struct {
		content string
		dir     string
		name    string
		backCo  color.Color
		forgeCo color.Color
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				content: "我很健康",
				dir:     "./",
				name:    "健康.png",
				backCo:  color.White,
				forgeCo: HexToRGBA("0fe285"),
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				content: "我是疑似",
				dir:     "./",
				name:    "疑似.png",
				backCo:  color.White,
				forgeCo: HexToRGBA("ffd038"),
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				content: "我是隔离",
				dir:     "./",
				name:    "隔离.png",
				backCo:  color.White,
				forgeCo: HexToRGBA("f8931d"),
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "",
			args: args{
				content: "我是确诊",
				dir:     "./",
				name:    "确诊.png",
				backCo:  color.White,
				forgeCo: HexToRGBA("f53e77"),
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildQRCode(tt.args.content, tt.args.dir, tt.args.name, tt.args.backCo, tt.args.forgeCo)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildQRCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildQRCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
