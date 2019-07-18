package KangDB

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_newShardmap(t *testing.T) {
	type args struct {
		shardNum int
	}
	tests := []struct {
		name string
		args args
		want shardmap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newShardmap(tt.args.shardNum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newShardmap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want DBInstance
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBInstance_GetShard(t *testing.T) {
	type fields struct {
		mu               sync.RWMutex
		bucket           shardmap
		shardNum         int
		IsActiveEviction bool
		activeeviction   chan struct{}
		EvictionInterval time.Duration
		hf               Hashfunc
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mapwithmutex
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DBInstance{
				mu:               tt.fields.mu,
				bucket:           tt.fields.bucket,
				shardNum:         tt.fields.shardNum,
				IsActiveEviction: tt.fields.IsActiveEviction,
				activeeviction:   tt.fields.activeeviction,
				EvictionInterval: tt.fields.EvictionInterval,
				hf:               tt.fields.hf,
			}
			if got := db.GetShard(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBInstance.GetShard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBInstance_Get(t *testing.T) {
	type fields struct {
		mu               sync.RWMutex
		bucket           shardmap
		shardNum         int
		IsActiveEviction bool
		activeeviction   chan struct{}
		EvictionInterval time.Duration
		hf               Hashfunc
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Item
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DBInstance{
				mu:               tt.fields.mu,
				bucket:           tt.fields.bucket,
				shardNum:         tt.fields.shardNum,
				IsActiveEviction: tt.fields.IsActiveEviction,
				activeeviction:   tt.fields.activeeviction,
				EvictionInterval: tt.fields.EvictionInterval,
				hf:               tt.fields.hf,
			}
			got, err := db.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBInstance.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBInstance.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBInstance_Set(t *testing.T) {
	type fields struct {
		mu               sync.RWMutex
		bucket           shardmap
		shardNum         int
		IsActiveEviction bool
		activeeviction   chan struct{}
		EvictionInterval time.Duration
		hf               Hashfunc
	}
	type args struct {
		key   string
		value interface{}
		ttl   int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DBInstance{
				mu:               tt.fields.mu,
				bucket:           tt.fields.bucket,
				shardNum:         tt.fields.shardNum,
				IsActiveEviction: tt.fields.IsActiveEviction,
				activeeviction:   tt.fields.activeeviction,
				EvictionInterval: tt.fields.EvictionInterval,
				hf:               tt.fields.hf,
			}
			db.Set(tt.args.key, tt.args.value, tt.args.ttl)
		})
	}
}

func TestDBInstance_Delete(t *testing.T) {
	type fields struct {
		mu               sync.RWMutex
		bucket           shardmap
		shardNum         int
		IsActiveEviction bool
		activeeviction   chan struct{}
		EvictionInterval time.Duration
		hf               Hashfunc
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DBInstance{
				mu:               tt.fields.mu,
				bucket:           tt.fields.bucket,
				shardNum:         tt.fields.shardNum,
				IsActiveEviction: tt.fields.IsActiveEviction,
				activeeviction:   tt.fields.activeeviction,
				EvictionInterval: tt.fields.EvictionInterval,
				hf:               tt.fields.hf,
			}
			if err := db.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("DBInstance.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBInstance_IsExists(t *testing.T) {
	type fields struct {
		mu               sync.RWMutex
		bucket           shardmap
		shardNum         int
		IsActiveEviction bool
		activeeviction   chan struct{}
		EvictionInterval time.Duration
		hf               Hashfunc
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DBInstance{
				mu:               tt.fields.mu,
				bucket:           tt.fields.bucket,
				shardNum:         tt.fields.shardNum,
				IsActiveEviction: tt.fields.IsActiveEviction,
				activeeviction:   tt.fields.activeeviction,
				EvictionInterval: tt.fields.EvictionInterval,
				hf:               tt.fields.hf,
			}
			if got := db.IsExists(tt.args.key); got != tt.want {
				t.Errorf("DBInstance.IsExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_activeEviction(t *testing.T) {
	type args struct {
		db *DBInstance
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activeEviction(tt.args.db)
		})
	}
}

func TestDBInstance_GracefulCloseDB(t *testing.T) {
	type fields struct {
		mu               sync.RWMutex
		bucket           shardmap
		shardNum         int
		IsActiveEviction bool
		activeeviction   chan struct{}
		EvictionInterval time.Duration
		hf               Hashfunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DBInstance{
				mu:               tt.fields.mu,
				bucket:           tt.fields.bucket,
				shardNum:         tt.fields.shardNum,
				IsActiveEviction: tt.fields.IsActiveEviction,
				activeeviction:   tt.fields.activeeviction,
				EvictionInterval: tt.fields.EvictionInterval,
				hf:               tt.fields.hf,
			}
			db.GracefulCloseDB()
		})
	}
}
