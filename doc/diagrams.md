# Architecture Diagrams

Visual documentation for the API Template Gin project. All diagrams use [Mermaid](https://mermaid.js.org/) syntax and render natively on GitHub.

---

## Table of Contents

- [1. Clean Architecture Overview](#1-clean-architecture-overview)
- [2. HTTP Request Flow](#2-http-request-flow)
- [3. Authentication Middleware Flow](#3-authentication-middleware-flow)
- [4. Server Startup Sequence](#4-server-startup-sequence)
- [5. Ports & Adapters (Class Diagram)](#5-ports--adapters-class-diagram)

---

## 1. Clean Architecture Overview

Hexagonal architecture with three layers. Dependencies always point inward — adapters depend on application, application depends on domain. External systems are accessed only through port interfaces.

```mermaid
flowchart TB
    subgraph External["External Systems"]
        Client([HTTP Client])
        Terminal([Terminal / CLI])
        PG[(PostgreSQL)]
    end

    subgraph Adapters["Adapters Layer"]
        direction TB

        subgraph HTTP["HTTP Adapter"]
            Handlers["handlers/\nGin HTTP Handlers"]
            DTO["dto/\nSuccessResponse & ErrorResponse"]
            Infra["infrastructure/\nGin Engine, Routes, Middleware"]
        end

        subgraph CLI["CLI Adapter"]
            CobraCLI["cli/\nCobra Subcommand"]
        end

        subgraph Repo["Repository Adapter"]
            GormRepo["repository/\nGORM + sync.Once Singleton"]
        end
    end

    subgraph Application["Application Layer"]
        Services["system_services/\nHealth Service"]
        Ports["ports/\nStore Interface"]
    end

    subgraph Domain["Domain Layer"]
        Entities["entities & rules\n(ready for expansion)"]
    end

    subgraph Pkg["Shared Packages"]
        Config["config/\ngodotenv + pflag"]
        Logger["log/\nLogrus wrapper"]
    end

    Client -->|HTTP| Infra
    Infra --> Handlers
    Handlers --> DTO
    Handlers --> Services
    Terminal -->|CLI| CobraCLI
    CobraCLI --> Services
    Services --> Ports
    Ports -.-|implements| GormRepo
    GormRepo -->|SQL| PG
    Services --> Entities
    Config -.->|used by| Infra
    Config -.->|used by| GormRepo
    Logger -.->|used by| Handlers
    Logger -.->|used by| GormRepo

    classDef external fill:#64748b,stroke:#475569,color:#fff
    classDef adapters fill:#3b82f6,stroke:#2563eb,color:#fff
    classDef application fill:#8b5cf6,stroke:#7c3aed,color:#fff
    classDef domain fill:#10b981,stroke:#059669,color:#fff
    classDef shared fill:#f59e0b,stroke:#d97706,color:#fff

    class Client,Terminal,PG external
    class Handlers,DTO,Infra,CobraCLI,GormRepo adapters
    class Services,Ports application
    class Entities domain
    class Config,Logger shared
```

---

## 2. HTTP Request Flow

End-to-end flow of a `GET /ping` request through all architectural layers.

```mermaid
sequenceDiagram
    actor Client
    participant Gin as Gin Engine
    participant Router as RegisterHandlersWithOptions
    participant Handler as Handler.Ping()
    participant DTO as dto.OK()

    Client->>+Gin: GET /ping
    Gin->>Gin: Logger & Recovery middleware
    Gin->>+Router: Match route (public group)
    Router->>+Handler: Invoke Ping(c *gin.Context)
    Handler->>DTO: dto.OK(c, gin.H{"ping": "pong"})
    DTO-->>Handler: SuccessResponse{Data, Meta} written
    Handler-->>-Router: Response sent
    Router-->>-Gin: Response written
    Gin-->>-Client: 200 {"data": {"ping": "pong"}, "meta": {"timestamp": "..."}}
```

---

## 3. Authentication Middleware Flow

How the Basic Auth middleware validates requests on protected routes. Public routes bypass this entirely.

```mermaid
sequenceDiagram
    actor Client
    participant Gin as Gin Engine
    participant Router as Route Groups
    participant Auth as basicAuthMiddleware
    participant Handler as Protected Handler
    participant Config as config.GetAuthenticationKey()

    Client->>+Gin: Request to protected route

    alt Public Route (/ping, /metrics)
        Gin->>+Router: Match public group
        Router->>Handler: Direct handler call (no auth)
        Handler-->>Router: Response
        Router-->>-Gin: 200 OK
    else Protected Route
        Gin->>+Router: Match protected group
        Router->>+Auth: Execute middleware chain
        Auth->>Config: Get AUTH_SECRET
        Config-->>Auth: Secret value

        alt Missing or invalid Authorization header
            Auth-->>Router: dto.Unauthorized(c, "Invalid or missing auth token")
            Router-->>Gin: 401 {"error": {"code": "UNAUTHORIZED", "message": "..."}}
        else Valid "Basic <base64(secret)>" header
            Auth->>Auth: Compare token == "Basic " + base64(secret)
            Auth-->>-Router: c.Next()
            Router->>+Handler: Invoke handler
            Handler-->>-Router: Response
            Router-->>Gin: 200 OK
        end
        Router-->>-Gin: Response sent
    end

    Gin-->>-Client: HTTP Response
```

---

## 4. Server Startup Sequence

Complete initialization flow from `main.go` to a running HTTP server with graceful shutdown.

```mermaid
flowchart TD
    Start([main.go]) --> LogStart["log.Info('Starting application...')"]
    LogStart --> Execute["app.Execute()"]
    Execute --> Init["Initialize()"]

    Init --> LoadEnv["config.LoadConfiguration()"]
    LoadEnv --> DotEnv["godotenv.Load(.env)"]
    DotEnv --> PFlags["pflag.Parse()"]
    PFlags --> MapEnv["loadEnvVariables()\nmap ENV → pflag keys"]
    MapEnv --> SetLog["log.SetLogLevel()"]
    SetLog --> SetEnvName["config.SetEnvironment()"]

    SetEnvName --> Cobra["rootCmd.Execute()"]
    Cobra --> CmdChoice{Subcommand?}

    CmdChoice -->|server| StartServer["StartServer()"]
    CmdChoice -->|cli -f test| CLIRun["cli.RunCliCmd()"]
    CmdChoice -->|none| Help["Show help"]

    StartServer --> NewServer["infrastructure.NewServer()"]
    NewServer --> Validate["serverConfig.Validate()"]
    Validate --> GinMode{"Mode?"}
    GinMode -->|debug| DebugMode["gin.SetMode(debug)\nconsole output"]
    GinMode -->|release| ReleaseMode["gin.SetMode(release)\nlog to file"]
    DebugMode --> CreateRouter["gin.Default()"]
    ReleaseMode --> CreateRouter
    CreateRouter --> Metrics["setMetrics(/metrics)"]
    Metrics --> Register["RegisterHandlersWithOptions()\npublic + protected groups"]
    Register --> LoadHandlers["loadHandlers()\nNewRestHandler()"]
    LoadHandlers --> ConfigLogger["log.ConfigureLogger()"]
    ConfigLogger --> ListenServe["srv.ListenAndServe()\n(goroutine)"]

    ListenServe --> Running([Server Running])
    Running --> WaitSignal["signal.Notify(SIGINT)"]
    WaitSignal --> Shutdown["srv.Shutdown(ctx)\n10s timeout"]
    Shutdown --> Stopped([Server Stopped])

    CLIRun --> HealthSvc["HealthService()"]
    HealthSvc --> RepoInit["NewRepository()\nsync.Once singleton"]
    RepoInit --> TestDB["healthService.TestDb()"]
    TestDB --> DBResult{Success?}
    DBResult -->|Yes| CLISuccess(["Print success"])
    DBResult -->|No| CLIError(["Print error"])

    classDef startEnd fill:#64748b,stroke:#475569,color:#fff
    classDef config fill:#f59e0b,stroke:#d97706,color:#fff
    classDef server fill:#3b82f6,stroke:#2563eb,color:#fff
    classDef cli fill:#8b5cf6,stroke:#7c3aed,color:#fff
    classDef decision fill:#f97316,stroke:#ea580c,color:#fff

    class Start,Running,Stopped,CLISuccess,CLIError,Help startEnd
    class LoadEnv,DotEnv,PFlags,MapEnv,SetLog,SetEnvName config
    class NewServer,Validate,DebugMode,ReleaseMode,CreateRouter,Metrics,Register,LoadHandlers,ConfigLogger,ListenServe,WaitSignal,Shutdown server
    class CLIRun,HealthSvc,RepoInit,TestDB cli
    class CmdChoice,GinMode,DBResult decision
```

---

## 5. Ports & Adapters (Class Diagram)

Interfaces (ports) and their concrete implementations (adapters). Shows the Dependency Inversion principle — application layer defines the contracts, adapter layer implements them.

```mermaid
classDiagram
    direction LR

    namespace Ports {
        class Store {
            <<interface>>
            +TestDb() error
        }

        class Health {
            <<interface>>
            +TestDb() error
        }

        class ServerInterface {
            <<interface>>
            +Ping(c *gin.Context)
        }
    }

    namespace Adapters {
        class repository {
            -db *gorm.DB
            +TestDb() error
        }

        class Handler {
            +Ping(c *gin.Context)
        }
    }

    namespace Application {
        class healthImp {
            -r Store
            +TestDb() error
        }

        class HealthService {
            <<factory>>
            +HealthService() (Health, error)
        }
    }

    namespace Infrastructure {
        class NewRepository {
            <<factory>>
            -once sync.Once
            -instance *repository
            +NewRepository() (Store, error)
            +NewConnection(dsn DBConfig) (Store, error)
        }

        class GinServer {
            +NewGinServer(handler ServerInterface) *gin.Engine
            +NewServer() *gin.Engine
            +RegisterHandlersWithOptions()
        }

        class GinServerOptions {
            +BaseURL string
            +Middlewares []gin.HandlerFunc
        }
    }

    namespace DTO {
        class SuccessResponse {
            +Data any
            +Meta Meta
        }

        class ErrorResponse {
            +Error ErrorDetail
        }

        class ErrorDetail {
            +Code ErrorCode
            +Message string
        }

        class Meta {
            +Timestamp string
        }
    }

    Store <|.. repository : implements
    Health <|.. healthImp : implements
    ServerInterface <|.. Handler : implements
    healthImp --> Store : depends on
    HealthService --> healthImp : creates
    HealthService --> NewRepository : uses
    Handler --> SuccessResponse : returns
    ErrorResponse --> ErrorDetail : contains
    GinServer --> ServerInterface : receives
    GinServer --> GinServerOptions : configures
    NewRepository --> repository : creates (singleton)
```

---

## Rendering

These diagrams render natively on:
- **GitHub** — Markdown preview
- **GitLab** — Markdown preview
- **VS Code** — With [Markdown Preview Mermaid](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-mermaid) extension

For local editing, use the [Mermaid Live Editor](https://mermaid.live/).
