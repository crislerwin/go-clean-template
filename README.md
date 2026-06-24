# Go Clean Template

Template minimalista e agnóstico de tecnologia para APIs em Go, seguindo Clean Architecture, Domain-Driven Design (DDD) e Ports and Adapters (Hexagonal Architecture).

## A Teoria (Theory Building)

> Source Code is Design.

O domínio expressa apenas regras de negócio e contratos (ports). O mundo externo (HTTP, banco de dados, cache, etc.) implementa esses contratos via adapters. Trocar uma tecnologia não deve alterar o domínio nem os casos de uso.

## Estrutura

```
cmd/api/                    # entrypoint da aplicação
internal/
  domain/<entidade>/        # entidades, value objects, erros de domínio
  application/<caso>/       # casos de uso (orquestram domínio e ports de saída)
  ports/
    input/                  # interfaces que a aplicação expõe para o mundo externo
    output/                 # interfaces que a aplicação requer do mundo externo
  |infrastructure/              # adapters externos
  |    persistence/memory/      # adapter de repositório em memória
  |    http/nethttp/            # adapter HTTP com net/http padrão
  |    logger/slog/             # logger estruturado com log/slog
  |    telemetry/otlp/          # tracer OTLP com fallback no-op
  |  config/                    # composição de adapters
  ```

## Dependências

As dependências apontam sempre para o centro:

- `domain` não conhece ninguém.
- `application` conhece `domain` e `ports/output`.
- `infrastructure` conhece `ports` e `domain`.

## Regras

- **TDD**: Red-Green-Refactor.
- **Testes table-driven** sempre que possível.
- **Sem lógica de domínio** em adapters HTTP ou repositórios.
- **Repositórios em memória** para testes e bootstrap.
- **Agnóstico de framework HTTP**: usa `net/http` padrão; fácil trocar por Gin, Echo, Fiber, etc.
- **OpenTelemetry opcional**: ativa apenas com `OTEL_EXPORTER_OTLP_ENDPOINT`; caso contrário usa no-op.
- **Logs estruturados**: JSON via `log/slog`, com `trace_id`/`span_id` para correlação.
- **LGTM stack ready**: docker-compose inclui Grafana, Loki, Tempo, Mimir, Alloy e OpenTelemetry Collector.
- **Docker-ready**: roda com `docker compose up`, sem banco externo.

## Comandos

```bash
make run                  # roda a API com repositório em memória
make test                 # testes unitários e de integração leves
make test-unit            # testes unitários
make lint                 # linters (requer golangci-lint)
make tidy                 # atualiza go.mod e go.sum
make install              # instala dependências e Air (opcional)
make docker-build         # build da imagem Docker
make docker-run           # roda container Docker localmente
make docker-compose-up    # sobe a API via Docker Compose
make docker-compose-down  # derruba a stack do Docker Compose
make lgtm-up              # sobe a stack completa LGTM + API
make lgtm-down            # derruba a stack LGTM e remove volumes
make lgtm-logs            # acompanha logs da API e do Alloy
make lgtm-ps              # status dos containers LGTM
```

## Exemplo de uso

### Local

Inicie o servidor:

```bash
make run
```

Crie um usuário:

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com"}'
```

Liste usuários:

```bash
curl http://localhost:8080/api/v1/users
```

Busque por ID:

```bash
curl http://localhost:8080/api/v1/users/<ID>
```

### Docker

Subir a API:

```bash
make docker-compose-up
```

Derrubar:

```bash
make docker-compose-down
```

### Healthcheck

A API expõe um endpoint de saúde independente da lógica de negócio:

```bash
curl http://localhost:8080/health
```

Resposta:

```json
{"status":"healthy"}
```

Esse endpoint é usado pelo Dockerfile e pelo docker-compose para verificar readiness/liveness.

### LGTM Stack (Grafana, Loki, Tempo, Mimir)

A stack completa sobe com um comando:

```bash
make lgtm-up
```

Serviços disponíveis:

| Serviço | URL | Propósito |
|---------|-----|-----------|
| API | http://localhost:8080 | API Go |
| Grafana | http://localhost:3000 | Dashboards (admin / ***) | metrics |
| Tempo | http://localhost:3200 | Distributed traces |
| Loki | http://localhost:3100 | Logs estruturados |
| Mimir | http://localhost:9009 | Métricas |

Para visualizar:

1. Acesse http://localhost:3000 (login: `admin` / ***).
2. O dashboard `go-clean-template LGTM` já vem provisionado.
3. Logs da API aparecem automaticamente no painel de logs.
4. Clique em um `trace_id` no log para saltar para o trace no painel do Tempo.

Para derrubar e limpar volumes:

```bash
make lgtm-down
```

### OpenTelemetry

Para enviar traces para o Tempo via OpenTelemetry Collector, a stack LGTM já configura automaticamente `OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318` e `OTEL_INSECURE=true`.

Para forçar TLS em produção, defina `OTEL_INSECURE=false`:

```bash
export OTEL_INSECURE=false
export OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector.example.com:4318
make run
```

Sem a variável `OTEL_EXPORTER_OTLP_ENDPOINT`, o tracer opera em modo no-op. Se o exporter falhar ao inicializar, a aplicação loga um aviso e continua com no-op.

### Logs estruturados

O logger padrão emite JSON com campos `trace_id`, `span_id` e `trace_flags` sempre que o contexto contém um span ativo. Isso permite correlacionar logs e traces no Grafana.

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com"}'
```

A saída no console (e no Loki) será similar a:

```json
{
  "time": "2026-06-23T03:00:00Z",
  "level": "INFO",
  "msg": "user created via HTTP",
  "trace_id": "abc123...",
  "span_id": "def456...",
  "user_id": "..."
}
```

## Ports

### Input (driven by infrastructure)

- `CreateUserUseCase`
- `FindUserByIDUseCase`
- `ListUsersUseCase`

### Output (driven by application)

- `UserRepository`
- `Tracer` (`internal/ports/telemetry`)
- `Logger` (`internal/ports/telemetry`)

## Extensão futura

- Trocar `infrastructure/http/nethttp` por `gin`, `echo` ou `fiber` sem tocar em `domain` ou `application`.
- Trocar `infrastructure/persistence/memory` por `postgres`, `mysql`, `mongo` — basta implementar `output.UserRepository`.
- Trocar `infrastructure/telemetry/otlp` por Jaeger, Zipkin ou stdout — basta implementar `telemetry.Tracer`.

## Design Decisions

- **`net/http` ao invés de framework**: reduz dependências e demonstra que adapters são trocáveis.
- **Injeção manual de dependências**: sem magic frameworks, fácil de testar e entender.
- **Erros de domínio exportados**: adapters podem reagir de forma diferente (`404` para `ErrUserNotFound`, por exemplo) sem vazar lógica de negócio.
- **OpenTelemetry como port**: o domínio/caso de uso dependem apenas de `telemetry.Tracer`, não do SDK OTel.
- **Logger como port**: `telemetry.Logger` permite trocar `slog` por outro logger estruturado sem tocar no domínio.
- **OpenTelemetry resiliente**: falha no exporter cai para no-op; TLS configurável via `OTEL_INSECURE`.
- **Endpoint `/health` dedicado**: healthcheck desacoplado da lógica de negócio.
- **Dockerfile multi-stage**: imagem pequena, sem ferramentas de build em runtime, rodando com usuário não-root.
- **Healthcheck no Dockerfile e Compose**: aponta para `/health`.
- **LGTM stack em docker-compose**: observabilidade local com um comando (`make lgtm-up`).
- **Alloy para log shipping**: coleta logs dos containers Docker e envia para Loki sem alterar a aplicação.
- **Derived fields no Grafana**: `trace_id` nos logs vira link para o trace no Tempo.
