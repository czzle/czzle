package multierr

import (
	"errors"
	"reflect"
	"testing"
)

func TestMultiErr_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		merr    *MultiErr
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			merr: Group("test").Code(1).New("not found").With(
				errors.New("test err 1"),
			).With(
				errors.New("test err 2"),
			),
			want: []byte(
				`{"code":0,"group":"","message":"test err 2","prev":` +
					`{"code":0,"group":"","message":"test err 1","prev":` +
					`{"code":1,"group":"test","message":"not found","prev":null}}}`,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := tt.merr
			got, err := mr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MultiErr.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MultiErr.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMultiErr_UnmarshalJSON(t *testing.T) {
	want := Group("test").Code(1).New("not found")
	Register(want)
	tests := []struct {
		name    string
		src     []byte
		wantErr bool
	}{
		{
			name: "success",
			src: []byte(
				`{"code":0,"group":"","message":"test err 2","prev":` +
					`{"code":0,"group":"","message":"test err 1","prev":` +
					`{"code":1,"group":"test","message":"not found","prev":null}}}`,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := new(MultiErr)
			if err := mr.UnmarshalJSON(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("MultiErr.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !mr.Has(want) {
				t.Errorf("MultiErr.UnmarshalJSON() does not have %v", want)
			}
		})
	}
}
