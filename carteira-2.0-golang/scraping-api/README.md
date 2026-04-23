# scraping-api

Serviço de scraping de dados fundamentalistas com camada de cache em memória. Parte do sistema [Carteira 2.0](../README.md).

Faz scraping de fundamentos de ações em tempo real a partir de múltiplas fontes financeiras brasileiras (Investidor10, Auvp, Fundamentus) e armazena os resultados em cache para reduzir requisições redundantes.

Roda na **porta 3001**.

---

## Pré-requisitos

- Go 1.21+
- Acesso à internet (o serviço faz scraping de sites externos)

> Não requer CGO nem SQLite — usa cache em memória puro Go.

---

## Rodando

```bash
go run ./cmd/main.go
```

O serviço irá:
1. Carregar configuração das variáveis de ambiente
2. Inicializar o cache em memória com o TTL configurado
3. Iniciar o servidor HTTP na porta `:3001`

---

## Configuração

| Variável | Padrão | Descrição |
|---|---|---|
| `CACHE_TTL_HOURS` | `24` | Por quanto tempo (em horas) os dados de scraping são considerados válidos. Deve ser um inteiro positivo; valores inválidos usam `24`. |
| `CACHE_ENABLED` | `true` | Defina como `false` para ignorar o cache e sempre fazer scraping fresco. Qualquer valor diferente de `false` é tratado como `true`. |

### Exemplos

```bash
# Desenvolvimento — TTL curto para dados mais frescos
CACHE_TTL_HOURS=1 go run ./cmd/main.go

# Produção — cache de 24 horas (padrão)
CACHE_TTL_HOURS=24 CACHE_ENABLED=true go run ./cmd/main.go

# Cache desabilitado (sempre scraping fresco)
CACHE_ENABLED=false go run ./cmd/main.go
```

---

## Comportamento do Cache

- **Cache hit**: Se existirem dados válidos (não expirados) para um ticker, são retornados imediatamente sem scraping.
- **Cache miss**: O serviço tenta cada scraper configurado em ordem. O primeiro que retornar dados completos vence, e o resultado é armazenado em cache.
- **Dados parciais**: Se um scraper retornar dados com campos nulos ou inválidos, o resultado **não é cacheado** e o próximo scraper é tentado.
- **Todos os scrapers falham**: O erro é retornado ao chamador. O cache não é modificado.
- **Cache desabilitado** (`CACHE_ENABLED=false`): Toda requisição dispara um scraping fresco.
- **Entrada expirada**: Tratada como cache miss; um scraping fresco é disparado.

---

## API

### GET /:ticker

Retorna os fundamentos de uma ação pelo código do ticker.

```
GET /WEGE3
```

**Resposta 200**

```json
{
  "Symbol": "WEGE3",
  "Price": 35.50,
  "PE": 28.4,
  "PBV": 8.1,
  "PSR": 4.2,
  "BVps": 4.38,
  "EPS": 1.25,
  "DY": 1.8,
  "Source": "Investidor10"
}
```

**Resposta 404** — ticker não encontrado ou todos os scrapers falharam

```json
{ "error": "failed to get stock data for XXXX: ..." }
```

### Campos da resposta

| Campo | Tipo | Descrição |
|---|---|---|
| `Symbol` | string | Código do ticker |
| `Price` | float | Preço atual |
| `PE` | float | Preço/Lucro |
| `PBV` | float | Preço/Valor Patrimonial |
| `PSR` | float | Preço/Receita |
| `BVps` | float | Valor Patrimonial por Ação |
| `EPS` | float | Lucro por Ação |
| `DY` | float | Dividend Yield (%) |
| `Source` | string | Fonte que forneceu os dados |

> Campos com valor `0` indicam que o scraper não conseguiu obter aquele dado para o ticker consultado.

O formato da resposta é idêntico independente de os dados virem do cache ou de um scraping fresco.

---

### DELETE /cache

Remove todas as entradas do cache em memória. Útil para forçar scraping fresco em todas as próximas requisições.

```
DELETE /cache
```

**Resposta 200**

```json
{ "message": "cache cleared" }
```

> Este endpoint não tem efeito quando `CACHE_ENABLED=false`, pois o cache já está desabilitado.

---

## Rodando os Testes

```bash
go test ./...
```

---

## Estrutura do Projeto

```
scraping-api/
├── cmd/
│   └── main.go                  # Ponto de entrada — configura cache e rotas
└── internal/
    ├── cache/
    │   └── cache_repository.go  # Cache em memória com TTL (go-cache)
    ├── config/
    │   └── config.go            # Carrega CACHE_TTL_HOURS e CACHE_ENABLED
    ├── http/
    │   └── http_client.go       # Cliente HTTP compartilhado pelos scrapers
    ├── models/
    │   └── stock_model.go       # Modelo StockData
    ├── repository/
    │   └── stock_repository.go  # Orquestra cache + scraper
    └── scraping/
        ├── scraper.go           # Interface do scraper
        ├── scraper_manager.go   # Tenta scrapers em ordem, retorna o primeiro resultado completo
        ├── scraper_rescraper.go # Lógica de re-scraping para resultados parciais
        ├── scrapers_configs.go  # Seletores de campos por fonte
        ├── sources_config.go    # URLs e prioridades das fontes
        └── helpers.go           # Utilitários de parsing
```
