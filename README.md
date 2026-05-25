```plaintext

grpa-core-go/
├── cmd/
│   └── grpa/
│       └── main.go         # Ponto de entrada (CLI). Trata os comandos: start, test, config
├── internal/
│   ├── config/
│   │   └── config.go       # Carrega variáveis de ambiente e validações do .env no startup
│   ├── orchestrator/
│   │   ├── worker.go       # Loop contínuo (Goroutine) que escuta a fila da API do GOEVO
│   │   └── dispatcher.go   # Gerencia a concorrência e dispara os scripts Python (Subprocessos)
│   ├── client/
│   │   └── goevo_api.go    # Cliente HTTP nativo e otimizado para se comunicar com o GOEVO
│   ├── storage/
│   │   ├── sqlite.go       # Inicializa e gerencia a conexão com o banco SQLite local
│   │   ├── migration.go    # Cria as tabelas locais (tarefas, logs, controle de versão) se não existirem
│   │   └── repository.go   # Operações de banco (Insert log, Update status do executor)
│   ├── updater/
│   │   └── updater.go      # Baixa os novos pacotes zip de automação e faz o hot-reload em disco
│   └── platform/
│       ├── logger.go       # Logger centralizado do sistema (usa a biblioteca nativa 'slog' do Go)
│       └── os_utils.go     # Captura sinais do S.O. (Ctrl+C, SIGTERM) para fechar o robô com segurança
├── automations/            # PASTA COMPARTILHADA (Onde o Go baixa e o Python lê)
│   ├── meu_sistema/
│   │   └── use_cases/
│   │       └── emitir_cnd.py
│   └── shared_venv/        # Ambiente virtual Python (.venv) compartilhado para os scripts
├── pkg/
│   └── sys_evidences/      # Pacote utilitário reutilizável para manipulação de arquivos de evidência (Base64)
├── data/
│   └── grpa.db             # O banco de dados SQLite real gerado localmente
├── go.mod                  # Declaração do módulo Go e versão
└── go.sum                  # Hash das dependências para garantir imutabilidade do Core
```

run:
```bash
 go run ./cmd/grpa/main.go .
```