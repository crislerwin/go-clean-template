# AGENTS Guidelines for go-clean-template

## Core Philosophy

Template minimalista para APIs Go seguindo Clean Architecture, DDD e Ports and Adapters.

## Rules of Engagement

### 1. Domain First

- O pacote `internal/domain` não importa nada de fora.
- Regras de negócio e erros de domínio vivem junto das entidades.

### 2. Ports Before Adapters

- Defina a interface (port) em `internal/ports` antes de criar o adapter.
- Adapters HTTP, repositórios, caches etc. implementam ports, nunca o contrário.

### 3. TDD

- Red: escreva o teste que falha.
- Green: implemente o mínimo.
- Refactor: melhore preservando comportamento.

### 4. Table-Driven Tests

- Prefira table-driven tests para casos de borda e múltiplos cenários.

### 5. Agnostic by Default

- Escolha a implementação mais simples e agnóstica possível.
- Só adicione dependência externa (PostgreSQL, Redis, etc.) quando houver requisito real.
- O adapter HTTP padrão usa `net/http` para demonstrar que frameworks são substituíveis.

## Commands

| Command           | Purpose                                  |
| ----------------- | ---------------------------------------- |
| `make run`        | Run API with in-memory repository.       |
| `make test`       | Run all tests.                           |
| `make test-unit`  | Run unit tests only.                     |
| `make lint`       | Run linters.                             |
| `make tidy`       | Update `go.mod` and `go.sum`.            |
| `make install`    | Install Go dependencies.                 |

## Project Structure

```
cmd/api/                        # entrypoint
internal/
  domain/<entity>/              # entidades e regras
  application/<usecase>/      # casos de uso
  ports/
    input/                      # contratos expostos para infra
    output/                     # contratos requeridos da infra
  infrastructure/
    persistence/memory/         # repositório em memória (trocável)
    http/nethttp/               # adapter HTTP padrão (trocável)
  config/                       # composição de adapters
```

## Adding a New Entity

1. Criar entidade em `internal/domain/<entity>` com testes.
2. Criar port de saída em `internal/ports/output/<entity>_repository.go`.
3. Criar ports de entrada em `internal/ports/input/<entity>_usecases.go`.
4. Implementar caso(s) de uso em `internal/application/<entity>` com testes.
5. Implementar repositório em memória em `internal/infrastructure/persistence/memory`.
6. Implementar handler em `internal/infrastructure/http/nethttp` com testes de tradução HTTP.
7. Compor em `cmd/api/main.go`.

## Why net/http?

Manter o adapter HTTP na biblioteca padrão remove o acoplamento a frameworks. Demonstra que a fronteira hexagonal é real: você pode trocar esse adapter por Gin, Echo, Fiber ou gRPC sem alterar o domínio.
