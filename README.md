# Run the code
run
```
./dev.sh
```

if db schema changes run

```
make gen
```


## env

### prod
dans le fichier `.env`

```
clientId=google_oauth_client_id
clientSecret=google_oauth_secret
DB_USER=postgres_username
DB_PASSWORD=postgres_password
DB_NAME=postgres_dbname
DB_PORT=postgres_port
```

### dev
si vous utilisez le script `dev.sh` -> mettre les mêmes variables dans le fichier `dev.env`
sinon les loader autrement


