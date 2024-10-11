# Web-Tp

## Description
Une application web pour gérer une classe.

## Installation
1. Cloner le dépôt:
    ```sh
    git clone <URL_DU_DEPOT>
    ```
2. Aller dans le dossier:
    ```sh
    cd Web-Tp
    ```
3. Installer les dépendances:
    ```sh
    go get -u ./...
    ```

## Utilisation
1. Lancer le serveur:
    ```sh
    go run main.go
    ```
2. Ouvrir `http://localhost:8080` dans le navigateur.

## Routes
- `/`: Accueil
- `/user/form`: Formulaire pour ajouter un étudiant
- `/user/display`: Affiche les infos de l'étudiant
- `/promo`: Affiche la promotion
- `/change`: Affiche le compteur de vues