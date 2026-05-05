package config

// DefaultServerAddr is the address the server listens on when not specified.
const DefaultServerAddr = "0.0.0.0"

// DefaultServerPort is the default HTTP listener port.
const DefaultServerPort = 8080

// ApplyDefaults fills in zero-value fields with sensible defaults.
func ApplyDefaults(cfg *Config) {
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = DefaultServerAddr
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = DefaultServerPort
	}
	for i := range cfg.Sources {
		if cfg.Sources[i].Labels == nil {
			cfg.Sources[i].Labels = make(map[string]string)
		}
	}
	for i := range cfg.Sinks {
		if cfg.Sinks[i].Options == nil {
			cfg.Sinks[i].Options = make(map[string]string)
		}
	}
}
