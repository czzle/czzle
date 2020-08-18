package uuid

import (
	"reflect"
	"testing"
)

func TestUUID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		u       UUID
		want    []byte
		wantErr bool
	}{
		{
			name:    "success_normal",
			u:       FromString("a3f8be07-17d1-42e6-a6c8-234767c4b88c"),
			want:    []byte("\"a3f8be07-17d1-42e6-a6c8-234767c4b88c\""),
			wantErr: false,
		},
		{
			name:    "success_null",
			u:       Null(),
			want:    []byte("null"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UUID.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UUID.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		args    args
		want    UUID
		wantErr bool
	}{
		{
			name: "success_normal",
			args: args{
				src: []byte("\"a3f8be07-17d1-42e6-a6c8-234767c4b88c\""),
			},
			want:    FromString("a3f8be07-17d1-42e6-a6c8-234767c4b88c"),
			wantErr: false,
		},
		{
			name: "success_null",
			args: args{
				src: []byte("null"),
			},
			want:    Null(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u UUID
			if err := u.UnmarshalJSON(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("UUID.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if !u.Equals(tt.want) {
				t.Errorf("UUID.UnmarshalJSON() got = %s, want %s", u, tt.want)
			}

		})
	}
}
