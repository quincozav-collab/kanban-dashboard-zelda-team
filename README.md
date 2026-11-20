# Tablero Kanban — Zelda Team (`kanban-dashboard-zelda-team`)

Aplicación web full-stack para gestionar el ciclo de vida de proyectos mediante tableros Kanban, control estricto de límites de trabajo en progreso (**WIP**), movimiento de tarjetas entre etapas y métricas de flujo en tiempo real (**Lead Time**, **Cycle Time** y conteo por columna).

El sistema se construye bajo principios de **Domain-Driven Design (DDD)** y metodología ágil **Kanban + MoSCoW**, sirviendo tanto como herramienta funcional de gestión como evidencia técnica de la asignatura *Procesos y Modelos Arquitectónicos de Software*.

> **Regla central del dominio:** el sistema impide crear o mover una tarjeta hacia una columna que ya alcanzó su límite WIP configurado, devolviendo un error de negocio controlado y alertando en la interfaz web para evitar cuellos de botella.

## Pila tecnológica

| Capa | Tecnología | Descripción | 
| ----- | ----- | ----- | 
| **Backend / API** | **Go (`net/http`)** | API REST nativa y servidor de archivos estáticos sin dependencias pesadas | 
| **Base de datos** | **MySQL 8** | Persistencia relacional de tableros, columnas, tarjetas y marcas de tiempo | 
| **Contenedores** | **Docker & Docker Compose** | Aislamiento y despliegue reproducible del motor de base de datos (`kanban_mysql`) | 
| **Frontend** | **HTML5 + CSS3 + JavaScript (ES6)** | Interfaz web reactiva conectada vía `fetch` y servida desde el mismo servidor Go | 
| **Arquitectura** | **Domain-Driven Design (DDD)** | Separación limpia entre dominio, aplicación e infraestructura | 

## Arquitectura DDD y estructura del proyecto

El código está organizado para mantener las reglas del negocio desacopladas de la infraestructura HTTP y de MySQL:

```
kanban-dashboard-zelda-team/
├── cmd/
│   └── api/
│       └── main.go                       # Punto de entrada: conexión a MySQL, rutas REST y servidor estático
├── internal/
│   ├── domain/                           # Núcleo: entidades, objetos de valor y reglas de negocio (WIP)
│   ├── application/                      # Casos de uso y orquestación del flujo Kanban
│   └── infrastructure/
│       ├── mysql_board_repository.go     # Persistencia MySQL, conteo WIP, movimiento y cálculo de métricas
│       └── http_handler.go               # Manejadores HTTP REST (JSON, validación y códigos de estado)
├── frontend/
│   └── index.html                        # Dashboard visual: métricas, formulario de creación y tablero WIP
├── migrations/                           # Scripts SQL de inicialización y esquema de tablas
├── docker-compose.yml                    # Configuración del servicio MySQL 8 en Docker
├── go.mod                                # Dependencias del módulo Go
└── README.md                             # Documentación técnica del proyecto

```

### Responsabilidad de cada capa

* **Dominio (`internal/domain/`):** Define las entidades principales (`Board`, `Column`, `Card`) y los invariantes del negocio, como la validación de capacidad máxima por columna.

* **Aplicación (`internal/application/`):** Coordina las operaciones de consulta, creación y traslado de tarjetas entre etapas.

* **Infraestructura (`internal/infrastructure/`):** Implementa el acceso a datos en MySQL 8 (`MySQLBoardRepository`) y expone los endpoints HTTP (`BoardHandler`).

* **Presentación (`frontend/index.html`):** Renderiza las columnas del flujo, los indicadores visuales de saturación WIP, las métricas promedio y los controles para crear y mover tarjetas.

## Flujo Kanban y límites WIP

El tablero modela el ciclo de vida del trabajo con límites explícitos de **Trabajo en Progreso (WIP)** para mantener un flujo continuo y detectar cuellos de botella de inmediato:

| Columna ID | Etapa del flujo | Límite WIP | Propósito | 
| ----- | ----- | ----- | ----- | 
| `col-1` | **Backlog (Por hacer)** | `3` | Tareas priorizadas listas para tomarse | 
| `col-2` | **En progreso / Desarrollo** | `2` | Tareas activas en implementación | 
| `col-3` | **En revisión y pruebas** | `2` | Validación técnica y cierre de etapa | 
| `col-4` | **Hecho (Done)** | Sin límite | Trabajo terminado y verificado | 

## Cómo ejecutar el proyecto

### 1. Levantar la base de datos MySQL 8 en Docker

Desde la raíz del proyecto, inicia el contenedor de base de datos:

```
docker-compose up -d

```

Verifica que el contenedor `kanban_mysql` esté activo:

```
docker ps

```

*(Opcional)* Si deseas conectarte directamente a la consola de MySQL para inspeccionar las tablas:

```
docker exec -it kanban_mysql mysql -u kanban_user -pkanban_password kanban_db

```

### 2. Iniciar el servidor backend y el frontend en Go

Ejecuta la aplicación desde la carpeta raíz del repositorio:

```
go run cmd/api/main.go

```

Verás en consola la confirmación de conexión a MySQL 8 y el servidor escuchando en el puerto `8080`.

### 3. Abrir el Dashboard en el navegador

Ingresa a:

* **Interfaz Web (Dashboard Kanban):** `http://localhost:8080/`

* **API REST Base:** `http://localhost:8080`

## Endpoints de la API REST

| Método | Ruta | Descripción | Respuesta esperada | 
| ----- | ----- | ----- | ----- | 
| `GET` | `/boards` | Obtiene los tableros y columnas configuradas | `200 OK` (JSON) | 
| `POST` | `/boards` | Crea un nuevo tablero en el sistema | `201 Created` | 
| `GET` | `/cards` | Lista todas las tarjetas registradas en el tablero | `200 OK` (JSON) | 
| `POST` | `/cards` | Crea una nueva tarjeta validando el límite WIP de la columna destino | `201 Created` / `409 Conflict` | 
| `PATCH` | `/cards/move` | Mueve una tarjeta a otra columna, valida WIP y registra `started_at` / `completed_at` | `200 OK` / `409 Conflict` | 
| `GET` | `/metrics` | Devuelve tarjetas completadas, Lead Time promedio, Cycle Time promedio y conteo por columna | `200 OK` (JSON) | 

### Ejemplo de respuesta de `GET /metrics`

```
{
  "avg_lead_time_hours": 4.98,
  "avg_cycle_time_hours": 1.98,
  "completed_cards": 1,
  "cards_by_column": {
    "col-1": 3,
    "col-3": 1
  }
}

```

## Cálculo de métricas de flujo

Para medir el desempeño real del equipo sin valores simulados, la tabla `cards` registra marcas de tiempo automáticas durante el recorrido de cada tarjeta:

* **`created_at`:** Marca el momento en que la tarjeta entra al tablero.

* **`started_at`:** Se registra automáticamente la primera vez que la tarjeta avanza a una etapa activa (`col-2`, `col-3` o `col-4`).

* **`completed_at`:** Se registra cuando la tarjeta alcanza la etapa de finalización (`col-3` o `col-4`).

* **Lead Time promedio:** Tiempo total transcurrido desde `created_at` hasta `completed_at` (en horas).

* **Cycle Time promedio:** Tiempo efectivo de ejecución desde `started_at` hasta `completed_at` (en horas).

## Equipo y contexto académico

* **Proyecto:** Kanban Dashboard (`kanban-dashboard-zelda-team`)

* **Asignatura:** Procesos y Modelos Arquitectónicos de Software (`IOC-2405`)

* **Institución:** Instituto Tecnológico Superior de Uruapan (TecNM)