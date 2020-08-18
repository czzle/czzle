package token

import (
	"reflect"
	"testing"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/pkg/uuid"
)

func TestToken_Sign(t *testing.T) {
	type fields struct {
		header  header
		payload Payload
	}
	type args struct {
		secret string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "sign",
			fields: fields{
				header: header{
					Version: V1,
				},
				payload: Payload{
					ClientID:  uuid.FromString("7a7097ed-0a09-4a04-9bb9-fadad1445f4b"),
					ClientIP:  "127.0.0.1",
					ExpiresAt: 1697429688183,
					IssuedAt:  1597429688183,
					Generator: "b789d597-d5d3-4517-b324-86d5eeeae254",
					Level:     czzle.Medium,
					Solved:    false,
				},
			},
			args: args{
				secret: "secret",
			},
			want: "eyJ2ZXIiOjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
				"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
				"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
				"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNDQ1ZjR" +
				"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
				"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
				"jU0In0.Valb7dN2C0mDej-lUpMtuZ6ZCyMSKeyb22" +
				"7B0fQiODA",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := Token{
				header:  tt.fields.header,
				payload: tt.fields.payload,
			}
			err := tkn.Sign(tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := tkn.String()
			if got != tt.want {
				t.Errorf("Token.Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		wantTkn Token
		wantErr error
	}{
		{
			name:    "empty",
			str:     "",
			wantErr: ErrMalformed,
		},
		{
			name:    "only_header",
			str:     "eyJ2ZXIiOjF9..",
			wantErr: ErrMalformed,
		},
		{
			name: "no_signature",
			str: "eyJ2ZXIiOjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
				"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
				"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
				"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNDQ1ZjR" +
				"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
				"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
				"jU0In0",
			wantErr: ErrMalformed,
		},
		{
			name: "incomplete_signature",
			str: "eyJ2ZXIiOjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
				"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
				"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
				"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNDQ1ZjR" +
				"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
				"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
				"jU0In0.Valb7dN2C0mDej",
			wantErr: ErrMalformed,
		},
		{
			name: "corrupted_header",
			str: "eyJ2ZXIiOsdfjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
				"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
				"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
				"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNDQ1ZjR" +
				"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
				"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
				"jU0In0.Valb7dN2C0mDej-lUpMtuZ6ZCyMSKeyb22" +
				"7B0fQiODA",
			wantErr: ErrMalformed,
		},
		{
			name: "corrupted_payload",
			str: "eyJ2ZXIiOjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
				"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
				"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
				"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNsdfDQ1ZjR" +
				"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
				"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
				"jU0In0.Valb7dN2C0mDej-lUpMtuZ6ZCyMSKeyb22" +
				"7B0fQiODA",
			wantErr: ErrMalformed,
		},
		{
			name: "valid",
			str: "eyJ2ZXIiOjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
				"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
				"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
				"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNDQ1ZjR" +
				"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
				"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
				"jU0In0.Valb7dN2C0mDej-lUpMtuZ6ZCyMSKeyb22" +
				"7B0fQiODA",
			wantErr: nil,
			wantTkn: Token{
				raw: "eyJ2ZXIiOjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
					"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
					"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
					"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNDQ1ZjR" +
					"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
					"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
					"jU0In0.Valb7dN2C0mDej-lUpMtuZ6ZCyMSKeyb22" +
					"7B0fQiODA",
				header: header{
					Version: V1,
				},
				payload: Payload{
					ClientID:  uuid.FromString("7a7097ed-0a09-4a04-9bb9-fadad1445f4b"),
					ClientIP:  "127.0.0.1",
					ExpiresAt: 1697429688183,
					IssuedAt:  1597429688183,
					Generator: "b789d597-d5d3-4517-b324-86d5eeeae254",
					Level:     czzle.Medium,
					Solved:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTkn, err := Parse(tt.str)
			if err != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTkn, tt.wantTkn) {
				t.Errorf("Parse() = %v, want %v", gotTkn, tt.wantTkn)
			}
		})
	}
}

func TestToken_Validate(t *testing.T) {
	type fields struct {
		raw     string
		header  header
		payload Payload
	}
	type args struct {
		secret string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "unsigned",
			args: args{
				secret: "secret",
			},
			fields: fields{
				raw: "",
			},
			want: false,
		},
		{
			name: "signed_valid",
			args: args{
				secret: "secret",
			},
			fields: fields{
				raw: "eyJ2ZXIiOjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
					"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
					"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
					"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNDQ1ZjR" +
					"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
					"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
					"jU0In0.Valb7dN2C0mDej-lUpMtuZ6ZCyMSKeyb22" +
					"7B0fQiODA",
			},
			want: true,
		},
		{
			name: "signed_invalid",
			args: args{
				secret: "badsecret",
			},
			fields: fields{
				raw: "eyJ2ZXIiOjF9.eyJzb2x2ZWQiOmZhbHNlLCJpYXQi" +
					"OjE1OTc0Mjk2ODgxODMsImV4cCI6MTY5NzQyOTY4O" +
					"DE4MywibHZsIjoibWVkaXVtIiwiY2lkIjoiN2E3MD" +
					"k3ZWQtMGEwOS00YTA0LTliYjktZmFkYWQxNDQ1ZjR" +
					"iIiwiY2lwIjoiMTI3LjAuMC4xIiwiZ2VuIjoiYjc4" +
					"OWQ1OTctZDVkMy00NTE3LWIzMjQtODZkNWVlZWFlM" +
					"jU0In0.Valb7dN2C0mDej-lUpMtuZ6ZCyMSKeyb22" +
					"7B0fQiODA",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := Token{
				raw:     tt.fields.raw,
				header:  tt.fields.header,
				payload: tt.fields.payload,
			}
			if got := tkn.Validate(tt.args.secret); got != tt.want {
				t.Errorf("Token.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
