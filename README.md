# Plataforma Distribuída de RPA (Robotic Process Automation)

Este documento descreve a arquitetura, o fluxo de dados e os padrões de desenvolvimento da nova plataforma distribuída de automação. A arquitetura foi reprojetada para eliminar o acoplamento com plataformas terceiras, garantir resiliência contra falhas de infraestrutura, permitir paralelismo horizontal seguro e fornecer controle total (Agendamento, Início, Parada e Reinício) sobre o ciclo de vida das tarefas.

---

## 1. Visão Geral da Arquitetura

A plataforma adota o padrão **Arquitetura Orientada a Eventos (EDA)** e separa rigorosamente o **Orquestrador/Cérebro** (escrito em Go) dos **Workers/Músculos** (escritos em Python + Playwright).

```plaintext
                  +---------------------------------------+
                  |           DASHBOARD / FRONTEND        |
                  +-------------------+-------------------+
                                      | HTTP
                                      v
                  +-------------------+-------------------+
                  |        ORQUESTRADOR (Go API)          |<--------+
                  +----+--------------+---------------+---+         |
                       |              |               |             |
         +-------------+              |               +--------+    | HTTP
         | Postgres                   | RabbitMQ               |    | (Finish/Log)
         v                            v                        v    |
+--------+--------+          +--------+--------+      +--------+----+--+
|   PostgreSQL    |          |   RabbitMQ      |      |     Redis      |
|                 |          |                 |      |                |
| - Histórico     |          | Fila:           |      | Canal de       |
| - Agendamento   |          | `rpa_tasks`     |      | Controle (Stop)|
| - Configurações |          +--------+--------+      +--------+-------+
+-----------------+                   |                        ^
                                      | Mensagem               | Espia
                                      v                        | Cancelamento
                             +--------+--------+               |
                             |  WORKER PYTHON  +---------------+
                             |    (Nó 1)       |
                             +--------+--------+
                                      |
                                      | Salva Prints/Logs
                                      v
                        +-------------+-------------+
                        |   VOLUME DOCKER COMPART.  |
                        |     (/var/rpa/storage)    |
                        +---------------------------+
```

### Componentes Principais

1. **Orquestrador (Go)**: Um binário único e estático encarregado de servir a API REST, gerenciar agendamentos (Cron), persistir dados no banco e coordenar o ciclo de vida das tarefas. Ele não executa automações web diretamente.
2. **Banco de Dados (PostgreSQL)**: Detém o estado absoluto do sistema. Armazena as tabelas de tarefas (`tasks`), agendamentos (`schedules`), logs persistentes e variáveis de ambiente criptografadas.
3. **Mensageria (RabbitMQ)**: Garante a resiliência da fila. Utiliza o mecanismo de *Acknowledgment* (ACK). Se um worker cair durante a execução, o RabbitMQ devolve a tarefa para a fila automaticamente.
4. **Canal de Controle (Redis)**: Banco em memória usado para guardar estados rápidos de sinalização de cancelamento de tarefas em tempo real (comandos `STOP`).
5. **Workers (Python + Playwright)**: Containers idênticos, síncronos e independentes. Eles entram em estado de espera (*listen*) consumindo a fila do RabbitMQ. Cada container executa uma única tarefa por vez.
6. **Armazenamento Local Baseado em Volumes (Docker Volumes)**: Solução de custo zero para persistência de artefatos (screenshots de erro e arquivos baixados). Compartilhado fisicamente entre o container Go e os containers Python no mesmo host.

---

## 2. Estrutura de Diretórios do Projeto

O ecossistema é organizado em duas pastas principais em um modelo de monorepo:

```plaintext
/
├── orchestrator-go/              # Código Fonte do Orquestrador (Go)
│   ├── cmd/api/main.go           # Ponto de entrada da API
│   ├── internal/scheduler/       # Agendador nativo (Cron/gocron)
│   ├── internal/queue/           # Publicador de mensagens do RabbitMQ
│   └── internal/database/        # Modelos e migrações do Postgres
│
├── worker-python/                # Código Fonte do Motor de Execução (Python)
│   ├── core/
│   │   ├── rabbit_consumer.py    # Consumer ativo conectado ao RabbitMQ
│   │   ├── browser_automation.py # Gerenciador do Playwright (Conplaintext Manager)
│   │   ├── worker.py             # Roteador (Dispatcher) e Importador Dinâmico
│   │   └── logger.py             # Emissor de logs estruturados
│   │
│   ├── automations/              # Os "Cartuchos" (Regras de Negócio)
│   │   └── [nome_do_sistema]/    # Ex: senior, receitafederal, solutionerp
│   │       ├── __init__.py
│   │       ├── common/           # Módulos reutilizáveis (login, preencher_nf)
│   │       └── use_cases/        # Scripts de automação finais
│   │           ├── __init__.py
│   │           └── fechar_folha.py
│   │
│   ├── requirements.txt
│   └── Dockerfile
│
└── docker-compose.yml            # Orquestração local da infraestrutura
```

---

## 3. Ciclo de Vida de uma Tarefa

```plaintext
[Cliente] --> (POST /api/tasks) --> [Go API] --> Grava 'pending' no Postgres
                                       |
                                       v
[Worker Python] <-- Puxa da fila <-- [RabbitMQ `rpa_tasks`]
       |
       +--> Notifica Go (PATCH /tasks/{id} -> status 'running')
       |
       +--> Executa `importlib.import_module()` e injeta o Browser Conplaintext
       |
       +--> Termina com sucesso ou captura erro (Gera Screenshot no Volume)
       |
       +--> Envia resultado/logs para o Go e dispara ACK para o RabbitMQ
```

### Mecanismo de Cancelamento (Sinal de STOP)

Para interromper uma tarefa de raspagem longa ou travada de forma segura:
1. O usuário clica em "Parar" no dashboard -> O Go altera o status da tarefa no Postgres para `canceling` e cria uma chave temporária no Redis: `SET task:[GUID]:status "CANCEL" EX 3600`.
2. Durante loops ou pontos críticos da automação, o script Python executa uma checagem rápida no Redis.
3. Ao detectar a string `"CANCEL"`, o Python lança uma exceção customizada `AutomationCancelledException`.
4. O Conplaintext Manager (`__exit__`) do `BrowserAutomation` intercepta a exceção, fecha o navegador limpando os processos do Chromium da memória, e notifica o Go da interrupção.

---

## 4. Padrão de Desenvolvimento de Automações (Python)

Para que o roteamento dinâmico funcione, todas as automações dentro da pasta `use_cases/` devem expor uma função obrigatória chamada `run`.

### Exemplo de um Caso de Uso (`worker-python/automations/senior/use_cases/fechar_folha.py`):

```python
import time
from core.exceptions import AutomationCancelledException

def run(page, log, data):
    """
    Contrato Padrão de Execução.
    :param page: Instância do Playwright Page injetada pronta pelo Core.
    :param log: Logger configurado e atrelado ao GUID da tarefa.
    :param data: Dicionário contendo o payload de entrada da tarefa.
    :return: Tupla (sucesso: bool, dados_de_retorno: dict)
    """
    log.info("Iniciando o processo de fechamento de folha...")
    
    # 1. Uso de módulos comuns reutilizáveis
    # de automations.senior.common.login import fazer_login
    # fazer_login(page, data['username'], data['password'])
    
    page.goto("[https://sistema.senior.com.br/folha](https://sistema.senior.com.br/folha)")
    
    # Simulação de um loop longo (Exemplo de checagem de cancelamento)
    for step in range(1, 5):
        # Lógica hipotética de verificação de cancelamento injetada no core
        # if core.redis.check_status(data['guid']) == "CANCEL":
        #     raise AutomationCancelledException("Tarefa abortada via painel.")
            
        log.info(f"Processando etapa {step}/4")
        time.sleep(1)
        
    page.click("#btn-finalizar")
    
    # Retorna o resultado que o Go irá armazenar
    return True, {"protocolo": "2026_FECH_9982"}
```

---

## 5. Estrutura de Dados Mapeada (PostgreSQL)

### Tabela de Tarefas (`tasks`)
Armazena a fila de execução persistente e o histórico.

```sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system VARCHAR(50) NOT NULL,          -- ex: 'senior'
    use_case VARCHAR(50) NOT NULL,        -- ex: 'fechar_folha'
    status VARCHAR(20) NOT NULL,          -- pending, running, success, failed, canceled
    payload JSONB NOT NULL,               -- parâmetros recebidos pela automação
    result_content JSONB,                 -- dados retornados pelo sucesso da automação
    error_message plaintext,                   -- rastro de erro em caso de falha
    screenshot_path VARCHAR(255),         -- caminho físico do print dentro do volume
    created_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    finished_at TIMESTAMP
);
```

### Tabela de Agendamentos (`schedules`)
Controla as tarefas recorrentes baseadas em tempo gerenciadas pelo Go.

```sql
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system VARCHAR(50) NOT NULL,
    use_case VARCHAR(50) NOT NULL,
    cron_expression VARCHAR(50) NOT NULL, -- ex: '0 9 * * *' (Todo dia às 09:00)
    payload JSONB NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 6. Configuração e Deploy Local (Docker Compose)

Abaixo está o arquivo de configuração para erguer toda a infraestrutura localmente compartilhando os volumes para salvar os artefatos sem custos com provedores S3.

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: rpa_postgres
    environment:
      POSTGRES_USER: rpa_user
      POSTGRES_PASSWORD: rpa_password
      POSTGRES_DB: rpa_platform
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  rabbitmq:
    image: rabbitmq:3-management-alpine
    container_name: rpa_rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672" # Painel administrativo do RabbitMQ

  redis:
    image: redis:7-alpine
    container_name: rpa_redis
    ports:
      - "6379:6379"

  orchestrator-go:
    build:
      conplaintext: ./orchestrator-go
    container_name: rpa_orchestrator
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - RABBITMQ_HOST=rabbitmq
      - REDIS_HOST=redis
      - STORAGE_PATH=/app/storage
    depends_on:
      - postgres
      - rabbitmq
      - redis
    volumes:
      - shared_artifacts:/app/storage

  worker-python:
    build:
      conplaintext: ./worker-python
    deploy:
      replicas: 3 # ESCALABILIDADE HORIZONTAL: Sobe 3 workers simultâneos nativamente
    environment:
      - ENV=docker
      - RABBITMQ_HOST=rabbitmq
      - REDIS_HOST=redis
      - STORAGE_PATH=/app/storage
    depends_on:
      - rabbitmq
      - redis
    volumes:
      - shared_artifacts:/app/storage

volumes:
  postgres_data:
  shared_artifacts: # Mapeado no host salvando fotos e logs com segurança
```

---

## 7. Próximos Passos & Funcionalidades Futuras

1. **Painel de Controle (Dashboard UI)**: Construir uma interface SPA simples plugada na API em Go para visualização do tamanho da fila do RabbitMQ e consumo de memória dos workers em tempo real.
2. **Políticas de Retentativa (Retry Policies)**: Configurar no Go a inteligência de re-enfileirar até 3 vezes automações que falharam por motivos de `SystemError` (como instabilidade ou quedas de conexão nos sites governamentais/ERPs).
3. **Gestão Centralizada de Proxies**: Adicionar uma tabela de Proxies no Postgres para que o Go injete credenciais de IP rotativas dinamicamente a cada payload enviado para burlar bloqueios de Cloudflare e Captchas.