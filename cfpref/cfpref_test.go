package cfpref

import (
	"reflect"
	"testing"
)

func TestCopyAppValue(t *testing.T) {
	type args struct {
		key   string
		appID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default",
			args: args{key: "HomePage", appID: "com.apple.safari"},
			want: "www.apple.com",
		},
		{
			name: "unset",
			args: args{key: "FooBarBaz", appID: "com.apple.safari"},
			want: "www.apple.com",
		},
	}
	for _, tt := range tests {
		if got := CopyAppValue(tt.args.key, tt.args.appID); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CopyAppValue() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
