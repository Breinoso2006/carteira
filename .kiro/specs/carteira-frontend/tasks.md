# Implementation Plan: carteira-frontend

## Overview

Implementação de uma SPA React com Vite que consome a `carteira-api` (porta 3002) e a `scraping-api` (porta 3001). A aplicação possui duas telas (`Portfolio_Screen` e `Analysis_Screen`) e um modal (`Stock_Details_Modal`), com navegação controlada por estado React sem React Router.

## Tasks

- [x] 1. Inicializar projeto e estrutura de arquivos
  - Criar projeto Vite + React com `npm create vite@latest carteira-frontend -- --template react` dentro de `carteira-2.0-golang/`
  - Criar a estrutura de diretórios: `src/api/`, `src/utils/`, `src/components/`, `src/styles/`
  - Criar arquivos vazios (stubs) para todos os módulos: `client.js`, `indicators.js`, `App.jsx`, `Portfolio_Screen.jsx`, `Analysis_Screen.jsx`, `Stock_List.jsx`, `Stock_Item.jsx`, `Stock_Form.jsx`, `Stock_Details_Modal.jsx`
  - Instalar dependência de teste: `npm install --save-dev vitest fast-check`
  - Configurar Vitest no `vite.config.js` (adicionar bloco `test: { environment: 'jsdom' }`)
  - _Requirements: 8.1, 8.2_

- [x] 2. Implementar o módulo API_Client
  - [x] 2.1 Implementar `src/api/client.js` com as funções `getPortfolio`, `addStock`, `updateStock`, `removeStock` e `getStockFundamentals`
    - Definir constantes `CARTEIRA_API = 'http://localhost:3002'` e `SCRAPING_API = 'http://localhost:3001'`
    - Implementar `handleResponse` que lança `Error` com a mensagem do corpo da resposta em caso de status 4xx/5xx
    - Cada função deve ser `async` e usar `fetch` nativo com headers `Content-Type: application/json` onde aplicável
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6_

  - [ ]* 2.2 Escrever testes de exemplo para o API_Client (`src/api/client.test.js`)
    - **Property 14: API_Client usa URLs base corretas**
    - **Validates: Requirements 8.1, 8.2**
    - **Property 15: API_Client propaga erros HTTP com mensagem da API**
    - **Validates: Requirements 8.4, 8.5**
    - Mockar `global.fetch` com `vi.fn()` para cada teste
    - Verificar que `getPortfolio` chama URL contendo `http://localhost:3002`
    - Verificar que `getStockFundamentals` chama URL contendo `http://localhost:3001`
    - Verificar que resposta 4xx lança `Error` com a mensagem retornada pela API

- [x] 3. Implementar a lógica de coloração de indicadores
  - [x] 3.1 Implementar `src/utils/indicators.js` com a função pura `getIndicatorColor(indicator, fundamentals)`
    - Cobrir os indicadores: `'pe'`, `'pbv'`, `'psr'`, `'dy'`, `'graham'`
    - Retornar `'neutral'` quando o campo relevante estiver em `invalid_fields` ou for `null`
    - Implementar cálculo Graham: `Math.sqrt(22.5 * eps * bvps)`, retornar `'red'` se `eps <= 0 || bvps <= 0`
    - _Requirements: 3.8, 3.9, 3.10, 3.11, 3.12, 3.13_

  - [ ]* 3.2 Escrever property tests para `getIndicatorColor` (`src/utils/indicators.test.js`)
    - **Property 7: Coloração correta dos indicadores PE, PBV, PSR, DY**
    - **Validates: Requirements 3.8, 3.9, 3.10, 3.11**
    - **Property 8: Coloração correta do indicador Graham**
    - **Validates: Requirements 3.12**
    - **Property 9: Campos inválidos omitem coloração**
    - **Validates: Requirements 3.13**
    - Usar `fc.assert` com `numRuns: 100` para cada propriedade
    - Incluir tag de referência `// Feature: carteira-frontend, Property N: ...` em cada teste

- [x] 4. Checkpoint — Garantir que todos os testes passam
  - Executar `npx vitest --run` e confirmar que todos os testes passam. Perguntar ao usuário se houver dúvidas.

- [x] 5. Implementar o componente raiz `App`
  - Implementar `src/App.jsx` com estado `currentScreen` (`'portfolio' | 'analysis'`) via `useState`
  - Renderizar condicionalmente `<Portfolio_Screen>` ou `<Analysis_Screen>` com base em `currentScreen`
  - Passar `onNavigate={setCurrentScreen}` como prop para ambas as telas
  - Garantir que apenas uma tela é renderizada por vez (nunca zero, nunca duas)
  - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [x] 6. Implementar `Stock_Item` e `Stock_List`
  - [x] 6.1 Implementar `src/components/Stock_Item.jsx`
    - Receber props: `stock`, `onEdit`, `onDelete`, `onViewDetails`
    - Exibir `stock.ticker` e `stock.fundamentalist_grade`
    - Renderizar botão "Editar" → chama `onEdit(stock)`
    - Renderizar botão "Remover" → chama `onDelete(stock.ticker)`
    - Renderizar botão "Ver Detalhes" → chama `onViewDetails(stock.ticker)`
    - _Requirements: 3.1, 5.1, 6.1_

  - [x] 6.2 Implementar `src/components/Stock_List.jsx`
    - Receber props: `stocks`, `onEdit`, `onDelete`, `onViewDetails`
    - Renderizar um `<Stock_Item>` para cada entrada em `stocks`, usando `stock.ticker` como `key`
    - _Requirements: 2.2_

- [x] 7. Implementar `Stock_Form`
  - Implementar `src/components/Stock_Form.jsx` com estado controlado para `ticker` e `fundamentalistGrade`
  - No modo `'add'`: ambos os campos são editáveis
  - No modo `'edit'`: campo `ticker` é somente leitura; preencher com `initialValues.ticker`
  - Validação antes de chamar `onSubmit`:
    - `ticker` não pode ser vazio → exibir mensagem de erro inline
    - `fundamentalist_grade` deve ser número real em `(0, 100]` → exibir mensagem de erro inline
  - Não chamar `onSubmit` se houver erros de validação
  - _Requirements: 4.1, 4.4, 4.5, 5.2_

- [x] 8. Implementar `Portfolio_Screen`
  - Implementar `src/components/Portfolio_Screen.jsx` com estado: `stocks`, `loading`, `error`, `formMode`, `editingStock`, `modalTicker`
  - Chamar `getPortfolio()` no mount via `useEffect` e após cada mutação bem-sucedida (add/update/delete)
  - Exibir indicador de carregamento enquanto `loading === true`
  - Exibir mensagem de erro quando `error !== null`
  - Exibir mensagem de lista vazia quando `stocks.length === 0` e não há loading/error
  - Renderizar `<Stock_List>` passando callbacks de edição, remoção e visualização
  - Renderizar `<Stock_Form>` no modo `'add'` por padrão; alternar para `'edit'` ao acionar edição de um item
  - Ao submeter o formulário no modo `'add'`: chamar `addStock`, depois recarregar lista
  - Ao submeter o formulário no modo `'edit'`: chamar `updateStock`, depois recarregar lista
  - Ao acionar remoção: chamar `removeStock`, depois recarregar lista
  - Ao acionar "Ver Detalhes": definir `modalTicker` com o ticker correspondente
  - Renderizar `<Stock_Details_Modal>` quando `modalTicker !== null`; fechar ao receber `onClose`
  - Limpar `error` antes de cada nova requisição
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 4.2, 4.3, 4.6, 5.3, 5.4, 5.5, 6.2, 6.3, 6.4_

- [x] 9. Implementar `Stock_Details_Modal`
  - Implementar `src/components/Stock_Details_Modal.jsx` com estado: `fundamentals`, `loading`, `error`
  - Chamar `getStockFundamentals(ticker)` no mount via `useEffect`
  - Exibir indicador de carregamento enquanto `loading === true`
  - Exibir mensagem de erro quando `error !== null`
  - Quando `fundamentals` disponível, exibir os campos: `symbol`, `price`, `pe`, `pbv`, `psr`, `bvps`, `eps`, `dy`, `source`
  - Para cada indicador (`pe`, `pbv`, `psr`, `dy`, `graham`): chamar `getIndicatorColor(indicator, fundamentals)` e aplicar classe CSS correspondente (`green`, `red` ou sem coloração para `neutral`)
  - Indicar visualmente campos presentes em `invalid_fields` (ex: texto riscado ou ícone de aviso)
  - Fechar ao clicar no botão de fechar (chama `onClose`) ou ao clicar no overlay externo
  - _Requirements: 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 3.10, 3.11, 3.12, 3.13_

- [x] 10. Implementar `Analysis_Screen`
  - Implementar `src/components/Analysis_Screen.jsx` com estado: `stocks`, `loading`, `error`
  - Chamar `getPortfolio()` no mount via `useEffect`
  - Exibir indicador de carregamento enquanto `loading === true`
  - Exibir mensagem de erro quando `error !== null`
  - Exibir mensagem de lista vazia quando `stocks.length === 0` e não há loading/error
  - Para cada stock, exibir: `ticker`, `fundamentalist_grade` e `weight` formatado como `(weight * 100).toFixed(2) + '%'`
  - Renderizar botão "Voltar" que chama `onNavigate('portfolio')`
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_

- [x] 11. Adicionar estilos e integrar tudo em `main.jsx`
  - Atualizar `src/main.jsx` para montar `<App />` no elemento `#root` do `index.html`
  - Criar `src/styles/index.css` com estilos básicos: layout, cores para indicadores (`.green`, `.red`), estilos do modal (overlay, container), indicador de loading, mensagens de erro e lista vazia
  - Importar `index.css` em `main.jsx`
  - _Requirements: 1.1, 3.8, 3.9, 3.10, 3.11, 3.12_

- [x] 12. Checkpoint final — Garantir que todos os testes passam
  - Executar `npx vitest --run` e confirmar que todos os testes passam. Perguntar ao usuário se houver dúvidas.

## Notes

- Tasks marcadas com `*` são opcionais e podem ser puladas para um MVP mais rápido
- Cada task referencia os requisitos específicos para rastreabilidade
- Os checkpoints garantem validação incremental
- Os property tests (fast-check) validam as propriedades de corretude universais da função `getIndicatorColor`
- Os testes de exemplo validam o comportamento do API_Client com fetch mockado
- A navegação entre telas é feita por estado React (`useState`) sem React Router — adequado para apenas duas telas
