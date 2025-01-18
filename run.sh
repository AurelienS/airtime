#!/bin/sh

# Étape 1: Télécharger les dépendances Go
go mod download

# Étape 2: Installer templ (si ce n'est pas encore fait)
go install github.com/a-h/templ/cmd/templ@latest

# Étape 3: Générer les templates avec templ
templ generate

# Étape 4: Vérifier si Node.js et npm sont installés, sinon les installer
if ! command -v npm &> /dev/null
then
    echo "npm (et npx) n'est pas installé. Installation de Node.js et npm..."
    curl -sL https://deb.nodesource.com/setup_16.x | bash -  # Installer Node.js (version LTS)
    apt-get install -y nodejs  # Installer npm et npx
fi

# Étape 5: Installer tailwindcss avec npm (en mode développement)
npm install -D tailwindcss

# Étape 6: Générer le fichier CSS avec Tailwind
npx tailwindcss -i ./web/view/styles.css -o ./web/static/styles.css

# Étape 7: Compiler le projet Go
go build -tags netgo -ldflags '-s -w' -o main cmd/server/main.go
