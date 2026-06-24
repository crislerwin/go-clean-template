# AGENTS Guidelines for go-clean-template

## Core Philosophy

Template minimalista para APIs Go seguindo Clean Architecture, DDD e Ports and Adapters.

## Rules of Engagement

### 1. Domain First

- O pacote `internal/domain` não importa nada de fora.
- Regras de negócio e erros de domínio vivem junto das entidades.

### 2. Ports Before Adapters

- Defina a interface (port) em `internal/ports` antes de criar o adapter.
- Adapters HTTP, repositórios, caches, loggers etc. implementam ports, nunca o contrário.
- O pacote `internal/ports/telemetry` define os contratos mínimos de tracing e logging para manter o domínio/casos de uso desacoplados de OTel e `log/slog`.

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
- OpenTelemetry é opcional: se `OTEL_EXPORTER_OTLP_ENDPOINT` não estiver configurado, o sistema usa no-op.
- TLS do OTLP é configurável via `OTEL_INSECURE` (`true` por padrão para desenvolvimento; `false` força TLS).
- Se a inicialização do exporter OTel falhar, a aplicação inicia com tracer no-op (resiliência sobre rigidez).
- O logger padrão (`internal/infrastructure/logger/slog`) emite JSON e extrai `trace_id`/`span_id` do contexto para correlação com traces.

### 6. Docker First

- A aplicação deve rodar com `docker compose up` sem necessidade de banco externo.
- Dockerfile multi-stage com usuário não-root e healthcheck apontando para `/health`.
- O endpoint `/health` não deve depender de lógica de negócio.
- A stack LGTM (Loki, Grafana, Tempo, Mimir, Alloy) sobe com `make lgtm-up` e é opcional; sem ela, a API continua funcionando.

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
| `make lgtm-up`           | Start full LGTM stack (API + Grafana + Loki + Tempo + Mimir + Alloy). |
| `make lgtm-down`         | Stop LGTM stack and remove volumes.      |
| `make lgtm-logs`         | Follow API and Alloy logs.               |
| `make lgtm-ps`           | Show LGTM container status.              |

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
    logger/slog/                # logger estruturado JSON (trocável)
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
7. Adicionar logs estruturados via `cfg.Logger` em handlers e casos de uso.
8. Adicionar spans de tracing via `cfg.Tracer` quando relevante.
9. Compor em `cmd/api/main.go`.

## Why net/http?

Manter o adapter HTTP na biblioteca padrão remove o acoplamento a frameworks. Demonstra que a fronteira hexagonal é real: você pode trocar esse adapter por Gin, Echo, Fiber ou gRPC sem alterar o domínio.

## Why OpenTelemetry as Optional?

O template precisa funcionar sem infraestrutura de observabilidade. O tracer OTLP ativa apenas quando `OTEL_EXPORTER_OTLP_ENDPOINT` está definido; caso contrário, usa no-op. Isso mantém o bootstrap simples e evita erros de conexão em desenvolvimento.

## Why Logger as Port?

O logger usa `log/slog` por baixo, mas a aplicação depende apenas de `telemetry.Logger`. Isso permite trocar por Zap, Zerolog ou outro logger estruturado sem alterar casos de uso. O adapter atual injeta `trace_id`, `span_id` e `trace_flags` automaticamente quando o contexto contém um span OTel ativo.

## Why LGTM Stack?

LGTM (Loki, Grafana, Tempo, Mimir) é a stack observability open-source da Grafana Labs. Ela centraliza logs, traces, métricas e dashboards com baixo custo local. A configuração em `docker-compose.yaml` é totalmente opcional: a API continua funcionando sem a stack, e vice-versa.

## Docker

- Dockerfile multi-stage gera imagem pequena e segura (usuário não-root).
- `docker compose up` levanta apenas a API por padrão.
- `make lgtm-up` levanta a API junto com Grafana, Loki, Tempo, Mimir, Alloy e OpenTelemetry Collector.
- Aplicação e LGTM stack são resilientes: se o OpenTelemetry Collector estiver indisponível, a API continua operando com tracer no-op.

## Pull Requests

Sempre abrir Pull Request para `main`. Não fazer push direto na `main`.
