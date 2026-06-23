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
- O pacote `internal/ports/telemetry` define o contrato mínimo de tracing para manter o domínio/casos de uso desacoplados de OTel.

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
- OpenTelemetry é opcional: se `OTEL_EXPORTER_OTLP_ENDPOINT` não estiver configurado, o sistema usa um tracer no-op.

### 6. Docker First

- A aplicação deve rodar com `docker compose up` sem necessidade de banco externo.
- Dockerfile multi-stage com usuário não-root e healthcheck.

## Commands

| Command                  | Purpose                                  |
| ------------------------ | ---------------------------------------- |
| `make run`               | Run API with in-memory repository.       |
| `make test`              | Run all tests.                           |
| `make test-unit`         | Run unit tests only.                     |
| `make lint`              | Run linters.                             |
| `make tidy`              | Update `go.mod` and `go.sum`.            |
| `make install`           | Install Go dependencies.                 |
| `make docker-build`      | Build Docker image.                      |
| `make docker-run`        | Run Docker container locally.            |
| `make docker-compose-up` | Start API via Docker Compose.            |
| `make docker-compose-down` | Stop Docker Compose stack.             |

## Project Structure

```
cmd/api/                        # entrypoint
internal/
  domain/<entity>/              # entidades e regras
  application/<usecase>/        # casos de uso
  ports/
    input/                      # contratos expostos para infra
    output/                     # contratos requeridos da infra
    telemetry/                  # contrato mínimo de tracing
  infrastructure/
    persistence/memory/         # repositório em memória (trocável)
    http/nethttp/               # adapter HTTP padrão (trocável)
    telemetry/
      otlp/                     # adapter OpenTelemetry OTLP/HTTP
      decorator/                # decorators de tracing para repositórios
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
8. Adicionar spans de tracing via `cfg.Tracer` quando relevante.

## Why net/http?

Manter o adapter HTTP na biblioteca padrão remove o acoplamento a frameworks. Demonstra que a fronteira hexagonal é real: você pode trocar esse adapter por Gin, Echo, Fiber ou gRPC sem alterar o domínio.

## Why OpenTelemetry as Optional?

O template precisa funcionar sem infraestrutura de observabilidade. O tracer OTLP ativa apenas quando `OTEL_EXPORTER_OTLP_ENDPOINT` está definido; caso contrário, usa no-op. Isso mantém o bootstrap simples e evita erros de conexão em desenvolvimento.

## Docker

- Dockerfile multi-stage gera imagem pequena e segura (usuário não-root).
- `docker compose up` levanta apenas a API por padrão.
- Descomente o serviço `otel-collector` em `docker-compose.yaml` para ativar coleta de traces.

## Pull Requests

Sempre abrir Pull Request para `main`. Não fazer push direto na `main`.
