package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	Grpc     GrpcConfig
	Temporal TemporalConfig
}

type TemporalConfig struct {
	HostPort  string `env:"TEMPORAL_HOST_PORT" envDefault:"localhost:7233"`
	Namespace string `env:"TEMPORAL_NAMESPACE" envDefault:"default"`
}

type GrpcConfig struct {
	ProductSvcAddr string `env:"PRODUCT_SERVICE_ADDRESS" envDefault:"localhost:50053"`
	OrderSvcAddr   string `env:"ORDER_SERVICE_ADDRESS" envDefault:"localhost:50054"`
	PaymentSvcAddr string `env:"PAYMENT_SERVICE_ADDRESS" envDefault:"localhost:50055"`
}

func Load() (*Config, error) {
	godotenv.Load()
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
