# Requirements Document

## Introduction

Frontend React para o sistema Carteira 2.0, consumindo a `carteira-api` (porta 3002) e a `scraping-api` (porta 3001). A aplicação é uma SPA (Single Page Application) com duas telas e um modal:

- **Portfolio_Screen**: tela principal onde o usuário gerencia as stocks cadastradas (adicionar, editar e remover) e pode abrir o modal de dados fundamentalistas de cada ação.
- **Stock_Details_Modal**: modal/dialog aberto a partir da Portfolio_Screen que exibe os dados fundamentalistas de uma stock obtidos via scraping-api, permitindo análise rápida sem sair da tela.
- **Analysis_Screen**: tela secundária que exibe os pesos calculados de cada ação retornados pela carteira-api.

A navegação entre telas é reativa, sem recarregamento de página.

## Glossary

- **App**: A aplicação frontend React como um todo.
- **Portfolio_Screen**: Tela principal que lista as stocks cadastradas e permite operações de CRUD.
- **Analysis_Screen**: Tela secundária que exibe os pesos calculados de cada stock retornados pelo endpoint `GET /portfolio`.
- **Stock_Form**: Formulário de adição ou edição de uma stock, contendo os campos `ticker` e `fundamentalist_grade`.
- **Stock_List**: Componente que renderiza a lista de stocks cadastradas no portfolio.
- **Stock_Item**: Componente que representa uma única entrada na Stock_List.
- **Stock_Details_Modal**: Componente de modal/dialog que exibe os dados fundamentalistas de uma stock retornados pela scraping-api.
- **API_Client**: Módulo responsável por todas as chamadas HTTP à `carteira-api` e à `scraping-api`.
- **carteira-api**: Backend REST rodando em `http://localhost:3002`.
- **scraping-api**: Serviço de scraping de dados fundamentalistas rodando em `http://localhost:3001`.
- **ticker**: Código identificador único de uma ação (ex: `WEGE3`).
- **fundamentalist_grade**: Nota fundamentalista de uma ação, valor real entre 0 (exclusivo) e 100 (inclusivo).
- **weight**: Peso calculado de uma ação no portfolio, retornado pela API como número decimal.
- **stock_fundamentals**: Conjunto de dados fundamentalistas de uma ação retornado pela scraping-api, contendo os campos `symbol`, `price`, `pe`, `pbv`, `psr`, `bvps`, `eps`, `dy`, `source` e `invalid_fields`.

---

## Requirements

### Requirement 1: Navegação Reativa entre Telas

**User Story:** Como usuário, quero navegar entre a tela de Portfolio e a tela de Análise sem recarregar a página, para que a experiência seja fluida e responsiva.

#### Acceptance Criteria

1. THE App SHALL renderizar a Portfolio_Screen como tela inicial ao ser carregado.
2. WHEN o usuário aciona o botão de navegação para a Analysis_Screen, THE App SHALL exibir a Analysis_Screen sem recarregar a página.
3. WHEN o usuário aciona o botão "Voltar" na Analysis_Screen, THE App SHALL exibir a Portfolio_Screen sem recarregar a página.
4. THE App SHALL manter apenas uma tela visível por vez.

---

### Requirement 2: Listagem de Stocks no Portfolio

**User Story:** Como usuário, quero ver todas as stocks cadastradas no portfolio na tela inicial, para que eu possa ter uma visão geral das minhas posições.

#### Acceptance Criteria

1. WHEN a Portfolio_Screen é exibida, THE API_Client SHALL enviar uma requisição `GET /portfolio` à carteira-api.
2. WHEN a carteira-api retorna a lista de stocks com sucesso, THE Stock_List SHALL renderizar um Stock_Item para cada entrada retornada.
3. WHEN a carteira-api retorna uma lista vazia, THE Portfolio_Screen SHALL exibir uma mensagem indicando que nenhuma stock está cadastrada.
4. IF a requisição `GET /portfolio` falhar, THEN THE Portfolio_Screen SHALL exibir uma mensagem de erro descritiva ao usuário.
5. WHILE a requisição `GET /portfolio` estiver em andamento, THE Portfolio_Screen SHALL exibir um indicador de carregamento.

---

### Requirement 3: Visualização de Dados Fundamentalistas de uma Stock

**User Story:** Como usuário, quero clicar em uma stock na tela de Portfolio e ver seus dados fundamentalistas em um modal, para que eu possa analisar rapidamente os indicadores da ação sem sair da tela.

#### Acceptance Criteria

1. THE Stock_Item SHALL exibir um controle de visualização de dados fundamentalistas para cada stock listada.
2. WHEN o usuário aciona o controle de visualização de um Stock_Item, THE API_Client SHALL enviar uma requisição `GET /:ticker` à scraping-api com o ticker correspondente.
3. WHILE a requisição `GET /:ticker` estiver em andamento, THE Stock_Details_Modal SHALL ser exibido com um indicador de carregamento.
4. WHEN a scraping-api retorna os dados fundamentalistas com sucesso, THE Stock_Details_Modal SHALL exibir os campos `symbol`, `price`, `pe`, `pbv`, `psr`, `bvps`, `eps`, `dy` e `source`.
5. WHEN a scraping-api retorna dados com `invalid_fields` não vazio, THE Stock_Details_Modal SHALL indicar visualmente quais campos não puderam ser obtidos.
6. IF a requisição `GET /:ticker` falhar, THEN THE Stock_Details_Modal SHALL exibir uma mensagem de erro descritiva ao usuário.
7. WHEN o usuário aciona o controle de fechamento do Stock_Details_Modal, THE Stock_Details_Modal SHALL ser removido da tela.
8. WHEN o Stock_Details_Modal exibe o indicador `pe`, THE Stock_Details_Modal SHALL colorir o indicador em verde se `pe > 0 AND pe <= 8`, e em vermelho caso contrário.
9. WHEN o Stock_Details_Modal exibe o indicador `pbv`, THE Stock_Details_Modal SHALL colorir o indicador em verde se `pbv > 0 AND pbv <= 2`, e em vermelho caso contrário.
10. WHEN o Stock_Details_Modal exibe o indicador `psr`, THE Stock_Details_Modal SHALL colorir o indicador em verde se `psr > 0 AND psr < 2`, e em vermelho caso contrário.
11. WHEN o Stock_Details_Modal exibe o indicador `dy`, THE Stock_Details_Modal SHALL colorir o indicador em verde se `dy >= 4`, e em vermelho caso contrário.
12. WHEN o Stock_Details_Modal exibe o indicador Graham, THE Stock_Details_Modal SHALL calcular o valor Graham como `sqrt(22.5 * eps * bvps)` e colorir o indicador em verde se `price < sqrt(22.5 * eps * bvps) AND eps > 0 AND bvps > 0`, e em vermelho caso contrário.
13. WHEN um campo utilizado no cálculo de coloração de um indicador estiver presente em `invalid_fields`, THE Stock_Details_Modal SHALL omitir a coloração desse indicador.

---

### Requirement 4: Adição de Nova Stock

**User Story:** Como usuário, quero adicionar uma nova stock ao portfolio informando o ticker e a nota fundamentalista, para que eu possa expandir minha carteira.

#### Acceptance Criteria

1. THE Portfolio_Screen SHALL exibir um Stock_Form para adição de nova stock.
2. WHEN o usuário submete o Stock_Form com `ticker` e `fundamentalist_grade` válidos, THE API_Client SHALL enviar uma requisição `POST /portfolio` com o body `{ ticker, fundamentalist_grade }`.
3. WHEN a carteira-api confirma a adição com sucesso, THE Stock_List SHALL ser atualizada para refletir a nova stock sem recarregar a página.
4. IF o campo `ticker` estiver vazio no momento da submissão, THEN THE Stock_Form SHALL exibir uma mensagem de validação indicando que o campo é obrigatório.
5. IF o campo `fundamentalist_grade` estiver fora do intervalo (0, 100] no momento da submissão, THEN THE Stock_Form SHALL exibir uma mensagem de validação indicando o intervalo permitido.
6. IF a requisição `POST /portfolio` falhar, THEN THE Portfolio_Screen SHALL exibir uma mensagem de erro descritiva ao usuário.

---

### Requirement 5: Atualização de Stock Existente

**User Story:** Como usuário, quero atualizar a nota fundamentalista de uma stock já cadastrada, para que eu possa manter os dados do portfolio atualizados.

#### Acceptance Criteria

1. THE Stock_Item SHALL exibir um controle de edição para cada stock listada.
2. WHEN o usuário aciona o controle de edição de um Stock_Item, THE Portfolio_Screen SHALL exibir o Stock_Form preenchido com os dados atuais da stock selecionada.
3. WHEN o usuário submete o Stock_Form de edição com `fundamentalist_grade` válido, THE API_Client SHALL enviar uma requisição `PUT /portfolio` com o body `{ ticker, fundamentalist_grade }`.
4. WHEN a carteira-api confirma a atualização com sucesso, THE Stock_List SHALL ser atualizada para refletir o novo valor sem recarregar a página.
5. IF a requisição `PUT /portfolio` falhar, THEN THE Portfolio_Screen SHALL exibir uma mensagem de erro descritiva ao usuário.

---

### Requirement 6: Remoção de Stock

**User Story:** Como usuário, quero remover uma stock do portfolio, para que eu possa manter a carteira com apenas as posições desejadas.

#### Acceptance Criteria

1. THE Stock_Item SHALL exibir um controle de remoção para cada stock listada.
2. WHEN o usuário aciona o controle de remoção de um Stock_Item, THE API_Client SHALL enviar uma requisição `DELETE /portfolio/:ticker` com o ticker correspondente.
3. WHEN a carteira-api confirma a remoção com sucesso, THE Stock_List SHALL ser atualizada removendo o Stock_Item correspondente sem recarregar a página.
4. IF a requisição `DELETE /portfolio/:ticker` falhar, THEN THE Portfolio_Screen SHALL exibir uma mensagem de erro descritiva ao usuário.

---

### Requirement 7: Exibição dos Pesos do Portfolio

**User Story:** Como usuário, quero ver os pesos calculados de cada ação na tela de Análise, para que eu possa entender a distribuição da minha carteira.

#### Acceptance Criteria

1. WHEN a Analysis_Screen é exibida, THE API_Client SHALL enviar uma requisição `GET /portfolio` à carteira-api.
2. WHEN a carteira-api retorna a lista com sucesso, THE Analysis_Screen SHALL exibir para cada stock: o `ticker`, o `fundamentalist_grade` e o `weight` formatado como percentual.
3. WHEN a carteira-api retorna uma lista vazia, THE Analysis_Screen SHALL exibir uma mensagem indicando que não há dados de portfolio para analisar.
4. IF a requisição `GET /portfolio` falhar, THEN THE Analysis_Screen SHALL exibir uma mensagem de erro descritiva ao usuário.
5. WHILE a requisição `GET /portfolio` estiver em andamento, THE Analysis_Screen SHALL exibir um indicador de carregamento.
6. THE Analysis_Screen SHALL exibir um botão que, quando acionado, navega de volta para a Portfolio_Screen.

---

### Requirement 8: Comunicação com as APIs

**User Story:** Como desenvolvedor, quero que todas as chamadas HTTP sejam centralizadas em um módulo dedicado, para que a manutenção e os testes sejam facilitados.

#### Acceptance Criteria

1. THE API_Client SHALL centralizar todas as requisições HTTP à carteira-api no endereço base `http://localhost:3002`.
2. THE API_Client SHALL centralizar todas as requisições HTTP à scraping-api no endereço base `http://localhost:3001`.
3. THE API_Client SHALL serializar e desserializar os corpos das requisições e respostas no formato JSON.
4. WHEN a carteira-api retorna um status HTTP de erro (4xx ou 5xx), THE API_Client SHALL propagar uma exceção com a mensagem de erro retornada pela API.
5. WHEN a scraping-api retorna um status HTTP de erro (4xx ou 5xx), THE API_Client SHALL propagar uma exceção com a mensagem de erro retornada pela scraping-api.
6. THE API_Client SHALL expor funções distintas para cada operação: listar portfolio, adicionar stock, atualizar stock, remover stock e buscar dados fundamentalistas de uma stock.
