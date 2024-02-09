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
GOOGLE_CLIENT_ID=id
GOOGLE_SECRET=secret
GOOGLE_CALLBACK_URL=callback
DB_USER=username
DB_PASSWORD=password
DB_NAME=dbname
DB_PORT=port
LOG_PATH=logpath
ENV=production
```

### dev
si vous utilisez le script `dev.sh` -> mettre les mêmes variables dans le fichier `dev.env` sauf :
`ENV=dev`
sinon les loader autrement


