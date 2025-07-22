# API de Noticias do G1 em Go

![Status](https://img.shields.io/badge/status-ativo-success.svg)
![Go](https://img.shields.io/badge/go-1.21%2B-blue.svg)
![Licença](https://img.shields.io/badge/licen%C3%A7a-MIT-green.svg)

Uma API robusta e performática construída em Go para extrair e exibir notícias mais recentes do portal G1. Este projeto foi desenvolvido com foco em boas práticas, arquitetura limpa e concorrência para alta eficiência.

## ✨ Features Principais

-   **Extração de Múltiplas Notícias**: Busca os comentários das 10 notícias mais recentes por categoria.
-   **Processamento Concorrente**: Utiliza Goroutines e Channels para processar as notícias em paralelo, garantindo alta velocidade.
-   **Filtragem por Categoria**: Permite buscar notícias de categorias específicas (ex: `economia`, `tecnologia`, `carros`).
-   **Arquitetura Hexagonal**: Desacoplamento total entre a lógica de negócio e as dependências externas (framework web, scrapers).
-   **Observabilidade**: Logs estruturados com `slog` para fácil monitoramento e depuração.
-   **Cache de Feeds**: Carrega e armazena os feeds RSS em memória para respostas mais rápidas.

## 🏛️ Arquitetura

Este projeto utiliza **Arquitetura Hexagonal (Portas e Adaptadores)**. Essa abordagem isola o "coração" da aplicação (lógica de negócio) de detalhes de implementação, como o framework web ou as fontes de dados. Isso torna o código mais testável, flexível e fácil de manter.

A estrutura de diretórios reflete essa separação:

```
g1-comments-api/
├── cmd/api/                  # Ponto de entrada da aplicação (main)
├── internal/
│   ├── core/
│   │   ├── domain/           # Modelos de dados puros (o coração)
│   │   └── ports/            # Interfaces (as Portas)
│   │   └── services/         # Lógica de negócio (implementação das portas)
│   ├── adapters/
│   │   ├── g1_scraper/       # Adaptador que raspa dados do G1 (Adaptador Secundário)
│   │   └── feed_reader/      # Adaptador que lê os feeds RSS (Adaptador Secundário)
│   └── handlers/
│       └── http/             # Adaptador que lida com requisições HTTP (Adaptador Primário)
└── go.mod
```

## 🚀 Começando

Siga os passos abaixo para executar a API localmente.

### Pré-requisitos

-   [Go](https://go.dev/doc/install) (versão 1.21 ou superior)
-   [Git](https://git-scm.com/)

### Instalação

1.  **Clone o repositório:**
    ```bash
    git clone [https://github.com/seu-usuario/g1-comments-api.git](https://github.com/seu-usuario/g1-comments-api.git)
    cd g1-comments-api
    ```

2.  **Instale as dependências:**
    O Go Modules cuidará disso automaticamente, mas para garantir que tudo está sincronizado, execute:
    ```bash
    go mod tidy
    ```

3.  **Execute a aplicação:**
    ```bash
    go run ./cmd/api/main.go
    ```

    O servidor será iniciado na porta `8080`. Você verá uma mensagem de log confirmando:
    ```json
    {"time":"...","level":"INFO","msg":"servidor iniciando na porta :8080"}
    ```

## 📡 Endpoints da API

A API expõe os seguintes endpoints para consulta. A resposta para todos é um **array de objetos JSON**, onde cada objeto representa uma notícia processada.

| Método | Endpoint                    | Descrição                                                              | Exemplo                                                |
| :----- | :-------------------------- | :--------------------------------------------------------------------- | :----------------------------------------------------- |
| `GET`  | `/`                         | Busca comentários **recentes** das últimas 10 notícias da categoria "todas". | `http://localhost:8080/`                               |
| `GET`  | `/:categoria`               | Busca comentários **recentes** das últimas 10 notícias da categoria especificada. | `http://localhost:8080/tecnologia`                     |
| `GET`  | `/populares`                | Busca comentários **populares** das últimas 10 notícias da categoria "todas". | `http://localhost:8080/populares`                      |
| `GET`  | `/:categoria/populares`     | Busca comentários **populares** das últimas 10 notícias da categoria especificada. | `http://localhost:8080/economia/populares`             |

### Exemplo de Resposta (`GET /carros`)

```json
[

  {
  "id": "multi-content/eefbd706-0717-4e81-ace6-c92dbfa7c4ca",
  "url": "https://g1.globo.com/inovacao/noticia/2025/07/21/quem-e-o-milionario-de-47-anos-que-tenta-reverter-o-tempo-e-voltar-a-ter-18-com-experimentos-no-corpo.ghtml",
  "title": "Quem é o milionário de 47 anos que tenta &#39;reverter o tempo&#39; e voltar a ter 18 com experimentos no corpo",
  "sub_title": "",
  "description": ""
  }
]
```

## 📦 Tecnologias Utilizadas

-   **[Go](https://go.dev/)**: Linguagem de programação principal.
-   **[Gin](https://github.com/gin-gonic/gin)**: Framework web de alta performance.
-   **[gofeed](https://github.com/mmcdole/gofeed)**: Parser robusto para feeds RSS e Atom.
-   **[slog](https://pkg.go.dev/log/slog)**: Biblioteca nativa para logging estruturado.

## 🧪 Testes

Para rodar os testes unitários e de integração, execute o seguinte comando na raiz do projeto:

```bash
go test ./... -v
```

## 🤝 Contribuição

Contribuições são bem-vindas! Se você tiver alguma ideia ou encontrar um bug, sinta-se à vontade para abrir uma *Issue* ou um *Pull Request*.

1.  Faça um *Fork* do projeto.
2.  Crie uma nova *Branch* (`git checkout -b feature/sua-feature`).
3.  Faça o *Commit* das suas alterações (`git commit -m 'Adiciona sua-feature'`).
4.  Faça o *Push* para a *Branch* (`git push origin feature/sua-feature`).
5.  Abra um *Pull Request*.

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

---
Feito com ❤️ e Goroutines por Programador Golang.