package config_test

import (
	"testing"

	"github.com/nem0z/hashplace/backend/internal/config"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		env     map[string]string
		want    config.Config
		wantErr bool
	}{
		{name: "default", env: map[string]string{}, want: config.Config{Addr: ":8080"}},
		{name: "custom addr", env: map[string]string{"HASHPLACE_ADDR": "127.0.0.1:9000"}, want: config.Config{Addr: "127.0.0.1:9000"}},
		{name: "port zero", env: map[string]string{"HASHPLACE_ADDR": ":0"}, want: config.Config{Addr: ":0"}},
		{name: "missing port", env: map[string]string{"HASHPLACE_ADDR": "localhost"}, wantErr: true},
		{name: "non numeric port", env: map[string]string{"HASHPLACE_ADDR": ":http"}, wantErr: true},
		{name: "port out of range", env: map[string]string{"HASHPLACE_ADDR": ":65536"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := config.Load(func(key string) string { return tt.env[key] })
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Load() = %+v, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
