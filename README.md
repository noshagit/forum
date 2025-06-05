# Forum

Une application de forum web simple avec authentification des utilisateurs, publication, commentaires et fonctionnalités de "like". Le projet est divisé en un backend Go et un frontend statique.

## Fonctionnalités

- Inscription et authentification des utilisateurs
- Création, visualisation et suppression de publications
- Commenter les publications
- Aimer les publications
- Gestion du profil utilisateur

## Structure du projet

```
forum/
├── back/           # Backend Go (API, gestionnaires, base de données)
│   ├── api/        # Points de terminaison API (JS)
│   ├── database/   # Base SQLite et scripts d'installation
│   ├── handlers/   # Gestionnaires HTTP Go
│   └── main/       # Point d'entrée principal Go
├── front/          # Frontend (HTML, CSS, JS)
│   ├── comments/   # Interface des commentaires
│   ├── images/     # Images statiques
│   ├── login/      # Interface de connexion
│   ├── password/   # Interface de réinitialisation du mot de passe
│   ├── post-list/  # Interface de liste des publications
│   ├── pp/         # Photos de profil
│   ├── profil/     # Interface de profil
│   └── register/   # Interface d'inscription
└── README.md       # Documentation du projet
```

## Pour commencer

### Prérequis

- Go (1.23.0)
- SQLite3
- CGO (compatibilité golang et sqlite)
- Git

## Utilisation

- Inscrivez-vous ou connectez-vous.
- Créez, consultez et interagissez avec les publications.
- Gérez votre profil et votre mot de passe.

### Accès au site

Deux choix s'offrent à vous :

1.  Lancer le serveur en local :

    - Accédez au répertoire principal :
      ```bash
      cd back/main
      ```
    - Lancez le serveur :
      ```bash
      go run main.go
      ```
    - Sur votre navigateur, entrez :
      "_http://localhost:8080_"

2.  Accéder au site hébergé sur Azur :

    - Sur votre navigateur, entrez :
      "*http://20.117.108.238:8080*"
