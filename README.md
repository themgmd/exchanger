# Curreency Service

## 1. Initialize

Create .env file from .env.example and fill them  
Change API key in configs/{dev/stable/testing}.yml file  
[You can get a new API Key here](freecurrencyapi.com)

Make migration with [golang-migrate](https://github.com/golang-migrate/migrate)

```shell
# run next scripts 
make dev
# or
make prod 
```