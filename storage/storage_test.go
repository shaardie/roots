package storage

import (
	"reflect"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Storage
	}{
		{
			name: "Test",
			want: &Storage{
				Storage: make(map[string]*WorldMap),
				m:       &sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_New(t *testing.T) {
	type args struct {
		name string
	}

	s := New()

	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    *WorldMap
		wantErr bool
	}{
		{
			name:    "map1",
			s:       s,
			args:    args{"NewMap"},
			want:    NewWorldMap(),
			wantErr: false,
		},
		{
			name:    "map2",
			s:       s,
			args:    args{"AnotherMap"},
			want:    NewWorldMap(),
			wantErr: false,
		},
		{
			name:    "map3",
			s:       s,
			args:    args{"NewMap"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.New(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.New() = %v, want %v", got, tt.want)
			}
		})
	}
}
