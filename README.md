# Wild-Life

[![Docker Go](https://github.com/ChicoState/Wild-Life/actions/workflows/docker-go.yml/badge.svg)](https://github.com/ChicoState/Wild-Life/actions/workflows/docker-go.yml)
[![Vue CI](https://github.com/ChicoState/Wild-Life/actions/workflows/vite.yml/badge.svg)](https://github.com/ChicoState/Wild-Life/actions/workflows/vite.yml)

Wildlife is a web app that allows users to upload an image and get analytic data about possible irritants on a plant or
other piece of nature.

<img src="docs/UI.png" alt="UI" width="200px"/>

## Install with docker-compose

Navigate yourself to the root of the repository directory where `docker-compose.yaml` is visible.

Build frontend and backend with docker-compose:

```bash
# Build the two containers
docker-compose build
```

Now, start both containers (headless or otherwise):

```bash
# Run docker-compose in the foreground
docker-compose up
# Run docker-compose in the background
docker-compose up -d
```

## Usage

### Development

Initialize a .env file in the root directory of the repository.

```bash
# Copy the example env to your local env
cp example.env .env
```

Ensure the port specified matches the configuration for the front-end.

```dotenv
PORT=5069
```

### Production

Initialize a .env file in the root directory of the repository. Docker-compose will feed the file into the container.

```bash
# Copy the example env to your local env
cp example.env .env
```

Ensure production mode is set to true

```dotenv
PRODUCTION=true
PORT=5069
```

#### Nginx

Route incoming requests on `https://wildlife.bradenn.com/*` to `http://127.0.0.1:5001/*`.

```nginx
location / {
    proxy_pass      http://127.0.0.1:5001;
}
```

Route incoming requests on `https://wildlife.bradenn.com/api/*` to `http://127.0.0.1:5069/*`.

```nginx
location /api/ {
    proxy_pass      http://127.0.0.1:5069/;
    
    # Set the file upload limit to 12MB
    client_max_body_size 12M;
}
```

Route incoming requests on `https://wildlife.bradenn.com/api/sockets/*` to `http://127.0.0.1:5069/sockets/*`. Except
this time we change the proxy header so the connection can be upgraded to a websocket.

```nginx
location /api/sockets {
    
proxy_buffering off;
proxy_set_header Host $host;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Real-IP $remote_addr;

proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
proxy_pass      http://127.0.0.1:5069/sockets;
}
```

## License

Copyright &copy; 2022 Braden Nicholson, Ryan Fong, Thomas Smale, and David Coles - All Rights Reserved. DO NOT
DISTRIBUTE.
