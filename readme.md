# Netcup DDNS docker

Simple docker container to update Netcup DNS A records using your Public IP.
Only supports IPv4.

## Usage

### Docker

```bash
docker run -e "DOMAIN=example.com" -e "SUBDOMAINS=abc,xyz" -e "INTERVAL=600" -e "NETCUP_CUSTOMER_NUMBER=123456" -e "NETCUP_API_KEY=adfsasf" -e "NETCUP_API_PASSWORD=adfsasf" --name netcup-ddns netcup-ddns
```

### Docker-Compose

`.env`

```ini
DOMAIN=example.com
SUBDOMAINS=abc,xyz
INTERVAL=600
NETCUP_CUSTOMER_NUMBER=123456
NETCUP_API_KEY=adfsasf
NETCUP_API_PASSWORD=adfsasf
```

```yaml
version: "3.7"

services:
  netcup-ddns:
    container_name: netcup-ddns
    hostname: netcup-ddns
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      DOMAIN: $DOMAIN
      SUBDOMAINS: $SUBDOMAINS
      INTERVAL: $INTERVAL
      NETCUP_CUSTOMER_NUMBER: $NETCUP_CUSTOMER_NUMBER
      NETCUP_API_KEY: $NETCUP_API_KEY
      NETCUP_API_PASSWORD: $NETCUP_API_PASSWORD
    restart: unless-stopped
```
