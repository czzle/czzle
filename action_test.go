package czzle

import (
	"reflect"
	"sort"
	"testing"
)

func TestAction_UnmarshalJSON(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantRes *Action
	}{
		{
			name: "flip",
			args: args{
				src: `{
					"type": "flip",
					"time": 5555555,
					"data": {
						"pos": {
							"x": 5,
							"y": 6
						}
					}
				}`,
			},
			wantRes: &Action{
				Type: FlipAction,
				Time: 5555555,
				Data: &FlipActionData{
					Pos: &Pos{
						X: 5,
						Y: 6,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "begin",
			args: args{
				src: `{
					"type": "begin",
					"time": 5555555
				}`,
			},
			wantRes: &Action{
				Type: BeginAction,
				Time: 5555555,
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "confirm",
			args: args{
				src: `{
					"type": "confirm",
					"time": 5555555
				}`,
			},
			wantRes: &Action{
				Type: ConfirmAction,
				Time: 5555555,
				Data: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := new(Action)
			if err := a.UnmarshalJSON([]byte(tt.args.src)); (err != nil) != tt.wantErr {
				t.Errorf("Action.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && !reflect.DeepEqual(a, tt.wantRes) {
				t.Errorf("Action.UnmarshalJSON() got = %#v, wantRes %#v", a, tt.wantRes)
			}

		})
	}
}

func TestActionList_Sort(t *testing.T) {
	tests := []struct {
		name string
		list ActionList
		want ActionList
	}{
		{
			name: "sort_by_time",
			list: ActionList{
				&Action{
					Time: 5,
				},
				&Action{
					Time: 6,
				},
				&Action{
					Time: 3,
				},
				&Action{
					Time: 2,
				},
				&Action{
					Time: 4,
				},
				&Action{
					Time: 1,
				},
			},
			want: ActionList{
				&Action{
					Time: 1,
				},
				&Action{
					Time: 2,
				},
				&Action{
					Time: 3,
				},
				&Action{
					Time: 4,
				},
				&Action{
					Time: 5,
				},
				&Action{
					Time: 6,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.list)
			for i := range tt.list {
				if tt.list[i].GetTime() != tt.want[i].GetTime() {
					t.Errorf("ActionList.Sort() [%d] = %d, want = %d ", i, tt.list[i].GetTime(), tt.want[i].GetTime())
				}
			}
		})
	}
}
