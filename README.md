# 💰 Gestor de Gastos Full-Stack

Aplicación web full-stack diseñada para el control, categorización y análisis de finanzas personales, construida con un stack de alto rendimiento y una arquitectura modular y escalable.

![Go](https://img.shields.io/badge/Backend-Go%20(Gin)-00ADD8?style=flat-square&logo=go&logoColor=white)
![React](https://img.shields.io/badge/Frontend-React%2019-61DAFB?style=flat-square&logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-Ready-3178C6?style=flat-square&logo=typescript&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=flat-square&logo=docker&logoColor=white)
![Turso](https://img.shields.io/badge/Database-Turso%20(SQLite)-41FFFF?style=flat-square&logo=turso&logoColor=black)

---

## 🚀 Características Principales

* **Autenticación de Usuarios:** Sistema de registro e inicio de sesión seguro con persistencia de sesión.
* **CRUD Completo de Gastos:** Permite registrar, editar y eliminar movimientos financieros de forma dinámica.
* **Resumen Financiero en Tiempo Real:** Cálculo automático de totales y distribución de gastos agrupados por categorías.
* **Filtros Avanzados:** Filtrado de movimientos por categoría y rangos de fechas específicos.
* **Interfaz Moderna y Responsiva:** UI diseñada con Tailwind CSS, optimizada para dispositivos móviles y de escritorio.
* **Arquitectura Escalable:** Backend estructurado en Go bajo una separación estricta de capas (handlers, repositorios y modelos).
* **Despliegue Contenedorizado:** Imagen Docker multi-etapa optimizada que compila el frontend con Bun y ejecuta la API de Go de manera eficiente.

---

## 🛠️ Stack Tecnológico

* **Frontend:** React, TypeScript, Tailwind CSS, Lucide Icons, Vite / Bun.
* **Backend:** Go (Golang), Gin Framework.
* **Base de Datos:** SQLite distribuida mediante **Turso**.
* **DevOps / Herramientas:** Docker, Docker Compose.

---

## 📁 Estructura del Proyecto

```text
GestorGastos/
├── backend/
│   ├── cmd/api/        # Punto de entrada de la aplicación Go
│   ├── internal/       # Lógica de negocio (auth, gastos, base de datos)
│   ├── go.mod          # Dependencias de Go
│   └── go.sum
├── frontend/
│   ├── src/            # Componentes React, vistas, tipos y cliente API
│   ├── package.json    # Dependencias del frontend
│   └── bun.lock
├── Dockerfile          # Configuración de compilación multi-etapa (Bun + Go)
└── docker-compose.yml  # Orquestación local de contenedores
```

⚙️ Cómo Ejecutar el Proyecto Localmente
Puedes clonar este repositorio y ejecutar la aplicación de dos maneras: mediante Docker (recomendado para simular producción) o de forma manual para desarrollo.

## Clona el repositorio:


```bash
git clone [https://github.com/RedelRodrigo/GestorGastos.git](https://github.com/RedelRodrigo/GestorGastos.git)
cd GestorGastos
```

### Configura tus variables de entorno creando un archivo .env en la raíz (o en el directorio del backend) con tus credenciales de Turso:

Fragmento de código
```env
TURSO_DATABASE_URL=tu_url_de_turso
TURSO_AUTH_TOKEN=tu_token_de_turso
```
### Opción A: Usando Docker Compose (Recomendado):
```bash
docker compose build --no-cache
docker compose up -d
```
Abre tu navegador y accede a: http://localhost:8080

### Opción B: Ejecución Manual (Modo Desarrollo)
1. Levantar el Backend (Go):

```Bash
cd backend
go mod download
go run cmd/api/main.go
```
(El servidor de Go correrá por defecto en el puerto 8080).

2. Levantar el Frontend (Vite / Bun):
Abre otra pestaña de la terminal:

```Bash
cd frontend
bun install
bun run dev
```
(El cliente de desarrollo de Vite se abrirá en http://localhost:5173, conectándose a tu API).

---
### 📄 Licencia
Este proyecto es de código abierto y está disponible bajo los términos de la licencia MIT.

***Desarrollado por Rodrigo Redel.***