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
  infrastructure/
    persistence/memory/     # adapter de repositório em memória
    http/nethttp/           # adapter HTTP com net/http padrão
  config/                   # composição de adapters
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

### OpenTelemetry

Para enviar traces, descomente o serviço `otel-collector` em `docker-compose.yaml` e exporte:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
make run
```

Sem essa variável, o tracer opera em modo no-op.

## Ports

### Input (driven by infrastructure)

- `CreateUserUseCase`
- `FindUserByIDUseCase`
- `ListUsersUseCase`

### Output (driven by application)

- `UserRepository`
- `Tracer` (`internal/ports/telemetry`)

## Extensão futura

- Trocar `infrastructure/http/nethttp` por `gin`, `echo` ou `fiber` sem tocar em `domain` ou `application`.
- Trocar `infrastructure/persistence/memory` por `postgres`, `mysql`, `mongo` — basta implementar `output.UserRepository`.
- Trocar `infrastructure/telemetry/otlp` por Jaeger, Zipkin ou stdout — basta implementar `telemetry.Tracer`.

## Design Decisions

- **`net/http` ao invés de framework**: reduz dependências e demonstra que adapters são trocáveis.
- **Injeção manual de dependências**: sem magic frameworks, fácil de testar e entender.
- **Erros de domínio exportados**: adapters podem reagir de forma diferente (`404` para `ErrUserNotFound`, por exemplo) sem vazar lógica de negócio.
- **OpenTelemetry como port**: o domínio/caso de uso dependem apenas de `telemetry.Tracer`, não do SDK OTel.
- **Dockerfile multi-stage**: imagem pequena, sem ferramentas de build em runtime, rodando com usuário não-root.
- **Healthcheck no Dockerfile**: permite orquestradores verificarem a saúde do container.
