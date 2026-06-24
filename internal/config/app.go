package config

import (
	"github.com/crislerwin/go-clean-template/internal/ports/output"
	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

// AppConfig agrega os adapters escolhidos para esta execução.
// A "theory": a aplicação é montada por composição de adapters
// que satisfazem as ports; trocar um adapter não exige mudança no domínio.
type AppConfig struct {
	UserRepository output.UserRepository
	Tracer         telemetry.Tracer
	Logger         telemetry.Logger
}
