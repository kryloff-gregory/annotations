package helpers

import "testing"

func TestIsValidURL(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Valid url",
			args{"https://youtu.be/A0e"},
			true,
		},
		{
			"Invalid scheme",
			args{"http://youtu.be/A0e"},
			false,
		},
		{
			"Invalid host",
			args{"https://youtube.com/A0e"},
			false,
		},
		{
			"Invalid param",
			args{"https://youtube.com/A0e$"},
			false,
		},
		{
			"No scheme",
			args{"youtube.com/A0e$"},
			false,
		},
		{
			"Extra param",
			args{"youtube.com/A0e$?good=true"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidURL(tt.args.link); got != tt.want {
				t.Errorf("IsValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
