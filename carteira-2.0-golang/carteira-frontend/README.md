# carteira-frontend

Interface web do sistema Carteira 2.0. Uma SPA (Single Page Application) construída com React + Vite que consome a `carteira-api` e a `scraping-api`.

## Visão Geral

A aplicação possui duas telas e um modal:

| Componente | Responsabilidade |
|---|---|
| `Portfolio_Screen` | Tela principal — lista stocks, formulário de CRUD, acesso ao modal de detalhes |
| `Analysis_Screen` | Tela secundária — exibe os pesos calculados de cada stock |
| `Stock_Details_Modal` | Modal sobreposto — exibe indicadores fundamentalistas com coloração semântica |

A navegação entre telas é controlada por estado React (`useState`), sem React Router.

---

## Pré-requisitos

- **Node.js 18+**
- **npm 9+**
- `carteira-api` rodando em `http://localhost:3002`
- `scraping-api` rodando em `http://localhost:3001`

---

## Instalação

```bash
cd carteira-2.0-golang/carteira-frontend
npm install
```

---

## Rodando em desenvolvimento

```bash
npm run dev
# Aplicação disponível em http://localhost:5173
```

> As duas APIs precisam estar rodando antes de abrir o frontend. Veja as instruções no README principal do projeto.

---

## Build para produção

```bash
npm run build
# Arquivos gerados em dist/
```

Para servir o build localmente:

```bash
npm run preview
```

---

## Testes

Os testes cobrem a lógica de coloração de indicadores (property-based) e o módulo de chamadas HTTP (testes de exemplo com fetch mockado).

```bash
npm run test
# ou
npx vitest --run
```

### Cobertura dos testes

| Arquivo | Tipo | O que valida |
|---|---|---|
| `src/utils/indicators.test.js` | Property-based (fast-check) | Coloração de PE, PBV, PSR, DY, Graham e campos inválidos |
| `src/api/client.test.js` | Exemplo com fetch mockado | URLs base corretas e propagação de erros HTTP |

---

## Estrutura do projeto

```
src/
├── main.jsx                    # Ponto de entrada — monta <App /> no DOM
├── App.jsx                     # Componente raiz — estado de navegação
│
├── api/
│   ├── client.js               # Todas as chamadas HTTP (carteira-api e scraping-api)
│   └── client.test.js          # Testes do API_Client
│
├── utils/
│   ├── indicators.js           # getIndicatorColor() — lógica de coloração de indicadores
│   └── indicators.test.js      # Property tests para getIndicatorColor
│
├── components/
│   ├── Portfolio_Screen.jsx    # Tela principal
│   ├── Analysis_Screen.jsx     # Tela de análise de pesos
│   ├── Stock_List.jsx          # Lista de stocks
│   ├── Stock_Item.jsx          # Item individual da lista
│   ├── Stock_Form.jsx          # Formulário de adição e edição
│   └── Stock_Details_Modal.jsx # Modal de dados fundamentalistas
│
└── styles/
    └── index.css               # Estilos globais
```

---

## Indicadores fundamentalistas

O modal de detalhes exibe os seguintes indicadores com coloração semântica:

| Indicador | Verde | Vermelho |
|---|---|---|
| P/E | `0 < pe ≤ 8` | caso contrário |
| P/BV | `0 < pbv ≤ 2` | caso contrário |
| PSR | `0 < psr < 2` | caso contrário |
| DY | `dy ≥ 4` | caso contrário |
| Graham | `price < √(22.5 × eps × bvps)` com `eps > 0` e `bvps > 0` | caso contrário |

Campos presentes em `invalid_fields` (retornados pela scraping-api) são exibidos sem coloração e com indicação visual de dado indisponível.

---

## Stack

- [React 19](https://react.dev/) — UI
- [Vite 8](https://vitejs.dev/) — bundler e dev server
- [Vitest](https://vitest.dev/) — test runner
- [fast-check](https://fast-check.io/) — property-based testing
