# Alejandrinas Web

Sitio web en Go que expone una API y vistas HTML usando [Echo](https://echo.labstack.com/). Utiliza PostgreSQL como base de datos, migraciones SQL planas y componentes generados con `templ`.

## Requisitos

- Go 1.25+
- PostgreSQL 14+ (local o remoto)P
- Herramienta `migrate` para ejecutar migraciones (`brew install golang-migrate` o binario oficial)
- Make (opcional, solo para atajos)

## Configuración local

1. Clona el repositorio y descarga dependencias:
   ```bash
   git clone git@github.com:tikimcrzx723/alejandrinasweb.git
   cd alejandrinasweb
   go mod download
   ```
2. Crea un archivo `.envrc` o exporta las variables necesarias:
   ```bash
   export SERVER_HOST=0.0.0.0
   export SERVER_PORT=8080
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_NAME=alejandrinasweb
   export DB_SSL_MODE=disable
   export DB_MIN_CONN=3
   export DB_MAX_CONN=20
   ```
3. Ejecuta migraciones:
   ```bash
   make migrate-up
   # o directamente
   migrate -path=./migrations -database=$DB_ADDR up
   ```
4. Corre la app:
   ```bash
   go run ./...
   # o compila
   go build -o bin/alejandrinasweb .
   ./bin/alejandrinasweb
   ```

## Despliegue en Ubuntu Server

1. **Preparar servidor**
   ```bash
   sudo apt update && sudo apt upgrade -y
   sudo apt install -y git curl make unzip postgresql-client nginx
   ```
   Instala Go 1.25 descargando el tarball oficial y añadiendo `/usr/local/go/bin` al `PATH`.

2. **Usuario y estructura**
   ```bash
   sudo adduser --disabled-password alejandrinas
   sudo mkdir -p /opt/alejandrinasweb/{src,bin} /var/log/alejandrinasweb
   sudo chown -R alejandrinas:alejandrinas /opt/alejandrinasweb /var/log/alejandrinasweb
   ```

3. **Código y build**
   ```bash
   sudo -u alejandrinas git clone git@github.com:tikimcrzx723/alejandrinasweb.git /opt/alejandrinasweb/src
   cd /opt/alejandrinasweb/src
   go build -o /opt/alejandrinasweb/bin/alejandrinasweb .
   ```

4. **Variables de entorno**
   ```bash
   sudo mkdir -p /etc/alejandrinasweb
   sudo tee /etc/alejandrinasweb/env >/dev/null <<'EOF'
   SERVER_HOST=0.0.0.0
   SERVER_PORT=8080
   DB_USER=alejandrinas
   DB_PASSWORD=********
   DB_HOST=127.0.0.1
   DB_PORT=5432
   DB_NAME=alejandrinasweb
   DB_SSL_MODE=disable
   DB_MIN_CONN=3
   DB_MAX_CONN=20
   EOF
   sudo chown root:alejandrinas /etc/alejandrinasweb/env
   sudo chmod 640 /etc/alejandrinasweb/env
   ```

5. **Servicio systemd**
   `/etc/systemd/system/alejandrinasweb.service`
   ```
   [Unit]
   Description=Alejandrinas Web
   After=network.target

   [Service]
   User=alejandrinas
   WorkingDirectory=/opt/alejandrinasweb/src
   EnvironmentFile=/etc/alejandrinasweb/env
   ExecStart=/opt/alejandrinasweb/bin/alejandrinasweb
   Restart=on-failure
   StandardOutput=append:/var/log/alejandrinasweb/app.log
   StandardError=append:/var/log/alejandrinasweb/error.log

   [Install]
   WantedBy=multi-user.target
   ```
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now alejandrinasweb
   sudo systemctl status alejandrinasweb
   ```

6. **Nginx + SSL**
   `/etc/nginx/sites-available/alejandrinasweb`:
   ```
   server {
       listen 80;
       server_name alejandrina.shop www.alejandrina.shop;

       location / {
           proxy_pass http://127.0.0.1:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```
   ```bash
   sudo ln -s /etc/nginx/sites-available/alejandrinasweb /etc/nginx/sites-enabled/
   sudo nginx -t && sudo systemctl reload nginx
   sudo certbot --nginx -d alejandrina.shop -d www.alejandrina.shop
   ```
   Certbot añadirá el bloque HTTPS y gestionará las renovaciones automáticas.

## Migraciones

- Crear una nueva migración:
  ```bash
  make migration nombre_migracion
  ```
- Ejecutar hacia arriba: `make migrate-up`
- Ejecutar hacia abajo: `make migrate-down NUM`

## Depuración y mantenimiento

- Ver estado del servicio: `sudo systemctl status alejandrinasweb`
- Logs de la app: `journalctl -u alejandrinasweb -f`
- Verificar puertos: `ss -tulpn | grep 8080`
- Probar sitio: `curl -I https://alejandrina.shop`

## Deploy manual rápido

```bash
ssh alejandrinas@your-server '
  cd /opt/alejandrinasweb/src &&
  git fetch --all &&
  git reset --hard origin/main &&
  go build -o /opt/alejandrinasweb/bin/alejandrinasweb . &&
  make migrate-up &&
  sudo systemctl restart alejandrinasweb
'
```


## Script de deploy

Se incluye `scripts/deploy.sh` para automatizar los pasos anteriores. Copia ese archivo al servidor (por ejemplo a `/opt/alejandrinasweb/deploy.sh`), hazlo ejecutable y lánzalo cada vez que quieras actualizar:

```bash
sudo install -m 755 scripts/deploy.sh /opt/alejandrinasweb/deploy.sh
sudo /opt/alejandrinasweb/deploy.sh
```

El script sincroniza `main`, compila el binario, intenta correr migraciones (si `migrate` está disponible y `DB_ADDR` existe) y reinicia el servicio systemd.

Adapta los valores (usuario, dominio, credenciales) a tu entorno antes de ejecutar los comandos.
