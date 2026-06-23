package config

import "github.com/crislerwin/go-clean-template/internal/ports/output"

// AppConfig agrega os adapters escolhidos para esta execução.
// A "theory": a aplicação é montada por composição de adapters
// que satisfazem as ports; trocar um adapter não exige mudança no domínio.
type AppConfig struct {
	UserRepository output.UserRepository
}
