# API Template Gin

Este proyecto es un template básico para desarrollar APIs utilizando Go y el framework Gin. Está diseñado para facilitar la creación de servicios RESTful con una estructura limpia y organizada.

## Tabla de Contenidos

- [Características](#características)
- [Tecnologías](#tecnologías)
- [Requisitos Previos](#requisitos-previos)
- [Instalación](#instalación)
- [Configuración](#configuración)
- [Ejecutar el Proyecto](#ejecutar-el-proyecto)
- [Ejecutar en Docker](#ejecutar-en-docker)
- [Uso](#uso)
- [Contribuciones](#contribuciones)
- [Licencia](#licencia)

## Características

- Estructura modular para facilitar la expansión y mantenimiento.
- Configuración opcional a través de variables de entorno.
- Registro y manejo de errores.
- Ejemplo de conexión a una base de datos.

## Tecnologías

- [Go](https://golang.org) - Lenguaje de programación
- [Gin](https://gin-gonic.com) - Framework HTTP para Go
- [Docker](https://www.docker.com) - Contenedorización
- [Pflag](https://github.com/spf13/pflag) - Manejo de banderas de línea de comandos
- [Godotenv](https://github.com/joho/godotenv) - Carga de variables de entorno desde un archivo `.env`

## Requisitos Previos

- Tener instalado [Go](https://golang.org/doc/install) (1.23 o superior).
- Tener instalado [Docker](https://docs.docker.com/get-docker/).
- Tener instalado [Docker Compose](https://docs.docker.com/compose/install/).

## Instalación

1. Clona el repositorio:

   ```bash
   git clone https://github.com/oswaldom-code/api-template-gin.git
   cd api-template-gin
   ```

2. (Opcional) Crea un archivo `.env` para tus variables de entorno:

   ```bash
   cp .env.example .env
   ```

3. Modifica el archivo `.env` según tus configuraciones.

## Configuración

Asegúrate de que el archivo `.env` contenga las siguientes variables:

```env
DB_USER=oswaldom-code
DB_PASSWORD=tu_contraseña
DB_HOST=localhost
DB_PORT=5432
SERVER_HOST=localhost
SERVER_PORT=9000
SERVER_MODE=debug
AUTH_SECRET=tu_secreto
```

## Ejecutar el Proyecto

### Ejecutar en Docker

Para construir y ejecutar la aplicación en un contenedor Docker, sigue estos pasos:

1. Asegúrate de estar en el directorio del proyecto.
2. Construye la imagen de Docker:

   ```bash
   sudo docker build -t api-template-app .
   ```

3. Ejecuta el contenedor, asegurándote de mapear los puertos y configuraciones necesarias:

   ```bash
   sudo docker run -d --name api-template-container -p 9000:9000 --env-file .env api-template-app
   ```

   En este comando:
   - `-d` ejecuta el contenedor en segundo plano.
   - `--name` da un nombre al contenedor para facilitar su referencia.
   - `-p 9000:9000` mapea el puerto 9000 del contenedor al puerto 9000 de tu máquina local.
   - `--env-file .env` carga las variables de entorno desde el archivo `.env`.

### Sin Docker

Para ejecutar la aplicación directamente en tu máquina:

1. Instala las dependencias:

   ```bash
   go mod tidy
   ```

2. Ejecuta la aplicación (modo desarrollo):

   ```bash
   go run main.go server
   ```

## Uso

Una vez que la aplicación esté en ejecución, puedes hacer solicitudes a la API utilizando herramientas como [Postman](https://www.postman.com/) o [curl](https://curl.se/).

Ejemplo de solicitud GET:

```bash
curl http://localhost:9000/ping
```

## Contribuciones

Las contribuciones son bienvenidas. Por favor, sigue estos pasos:

1. Haz un fork del proyecto.
2. Crea una rama para tu feature (`git checkout -b feature/nueva-feature`).
3. Realiza tus cambios y haz un commit (`git commit -m 'Agregada nueva feature'`).
4. Haz un push a la rama (`git push origin feature/nueva-feature`).
5. Abre un Pull Request.

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

