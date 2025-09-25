# TP3 Perso — Pipeline Concurrent & JSON (Go)

Contexte et objectif  
Ce petit projet est un exercice personnel réalisé à partir du Cours 3 (concurrence, channels, mutex) : il montre un pipeline simple de producteurs -> canal -> consommateur, la sauvegarde d'un historique en JSON et une démonstration de conditions de course (race) avec et sans mutex. C'est un TP d'initiation pour pratiquer les concepts vus en cours.

Résumé des fonctionnalités

- Lancer un pipeline concurrent de producteurs qui émettent des événements (Event).
- Collecter les événements depuis un channel tamponné et les sauvegarder dans un fichier JSON (`events.json`).
- Option "show" pour relire et afficher l'historique JSON.
- Option "race" pour démontrer l'effet d'un mutex vs accès non protégé (comportement non déterministe).
- Timeout gérant la durée maximale d'écoute du canal (sécurité en cas de blocage).

Structure (fichiers principaux)

- `main.go` : implémentation entière (producteurs, consommateur, JSON I/O, démonstration race).
- `events.json` : fichier généré après `run` (historique des événements).

Comment utiliser (PowerShell - Windows)

- Se placer dans le dossier du projet :

- Exécuter le pipeline (génère et sauvegarde `events.json`) :

```powershell
go run main.go run
```

(Équivalent si vous avez compilé :)

```powershell
go build -o tp3-perso.exe .
.\tp3-perso.exe run
```

- Afficher l'historique (relit `events.json`) :

```powershell
go run main.go show
# ou
.\tp3-perso.exe show
```

- Démo Race (affiche compteur avec et sans mutex) :

```powershell
go run main.go race
# ou
.\tp3-perso.exe race
```

Explication rapide du comportement (ce que fait le programme)

- `run`
  - Lance deux producteurs (`alpha`, `beta`) chacun produisant N événements.
  - Les producteurs écrivent sur un channel tamponné.
  - Le consommateur lit le channel jusqu'à fermeture ou timeout, collecte les événements et les écrit dans `events.json`.
  - Timeout (3s) protège contre blocage si quelque chose se passe mal.
- `show`
  - Lit `events.json`, désérialise et affiche une liste lisible des événements.
- `race`
  - Deux démonstrations : version protégée (avec `sync.Mutex`) et version non protégée (sans mutex) pour comparer les résultats d'un compteur partagé par 1000 goroutines.

Points techniques et concepts utilisés

- Channels tamponnés pour découpler producteurs/consommateur.
- sync.WaitGroup pour attendre la fin des producteurs.
- select + time.After pour implémenter un timeout.
- JSON (encoding/json) pour persister l'historique.
- erreurs gérées via `errors.Is`/comparaisons simples.
- Mutex (`sync.Mutex`) pour éviter les races ; tester le détecteur de race avec `go run -race main.go race`.

Exemples d'usage et sortie attendue

- Pendant `run` vous verrez des lignes du type :
  - "⬇️ Reçu: alpha | message #1 de alpha"
  - "💾 Historique écrit dans events.json"
- `show` affiche les événements avec timestamp et source.

Conseils / améliorations possibles

- Numérotation globale des events (ID unique) au lieu de repartir à 1 par producteur.
- Persistance append (au lieu d'écraser) et rotation du fichier.
- Exposer paramètres (nombre d'événements, délai, timeout) via flags CLI.
- Ajouter tests unitaires et utiliser `go test` pour functions isolées.
- Rendre le pipeline reentrant/reconnectable (fichiers de config).

Notes pratiques

- `events.json` est créé dans le répertoire courant. Lance toujours les commandes depuis le dossier du projet pour retrouver facilement le fichier.
- Pour détecter race conditions sur la version non protégée :

```powershell
go run -race main.go race
```

Auteur  
Projet réalisé seul à des fins pédagogiques (Cours 3).  
Paix à tous


# TP3 Perso — Pipeline Concurrent & JSON (Go)

Contexte et objectif
Ce petit projet est un exercice personnel réalisé à partir du Cours 3 (concurrence, channels, mutex) : il montre un pipeline simple de producteurs -> canal -> consommateur, la sauvegarde d'un historique en JSON et une démonstration de conditions de course (race) avec et sans mutex. C'est un TP d'initiation pour pratiquer les concepts vus en cours.

Résumé des fonctionnalités
- Lancer un pipeline concurrent de producteurs qui émettent des événements (Event).
- Collecter les événements depuis un channel tamponné et les sauvegarder dans un fichier JSON (`events.json`).
- Option "show" pour relire et afficher l'historique JSON.
- Option "race" pour démontrer l'effet d'un mutex vs accès non protégé (comportement non déterministe).
- Timeout gérant la durée maximale d'écoute du canal (sécurité en cas de blocage).

Structure (fichiers principaux)
- `main.go` : implémentation entière (producteurs, consommateur, JSON I/O, démonstration race).
- `events.json` : fichier généré après `run` (historique des événements).

Comment utiliser (PowerShell - Windows)
- Se placer dans le dossier du projet :

- Exécuter le pipeline (génère et sauvegarde `events.json`) :

```powershell
go run main.go run
```

(Équivalent si vous avez compilé :)

```powershell
go build -o tp3-perso.exe .
.\tp3-perso.exe run
```

- Afficher l'historique (relit `events.json`) :

```powershell
go run main.go show
# ou
.\tp3-perso.exe show
```

- Démo Race (affiche compteur avec et sans mutex) :

```powershell
go run main.go race
# ou
.\tp3-perso.exe race
```

Explication rapide du comportement (ce que fait le programme)

- `run`
  - Lance deux producteurs (`alpha`, `beta`) chacun produisant N événements.
  - Les producteurs écrivent sur un channel tamponné.
  - Le consommateur lit le channel jusqu'à fermeture ou timeout, collecte les événements et les écrit dans `events.json`.
  - Timeout (3s) protège contre blocage si quelque chose se passe mal.
- `show`
  - Lit `events.json`, désérialise et affiche une liste lisible des événements.
- `race`
  - Deux démonstrations : version protégée (avec `sync.Mutex`) et version non protégée (sans mutex) pour comparer les résultats d'un compteur partagé par 1000 goroutines.

Points techniques et concepts utilisés

- Channels tamponnés pour découpler producteurs/consommateur.
- sync.WaitGroup pour attendre la fin des producteurs.
- select + time.After pour implémenter un timeout.
- JSON (encoding/json) pour persister l'historique.
- erreurs gérées via `errors.Is`/comparaisons simples.
- Mutex (`sync.Mutex`) pour éviter les races ; tester le détecteur de race avec `go run -race main.go race`.

Exemples d'usage et sortie attendue

- Pendant `run` vous verrez des lignes du type :
  - "⬇️ Reçu: alpha | message #1 de alpha"
  - "💾 Historique écrit dans events.json"
- `show` affiche les événements avec timestamp et source.

Conseils / améliorations possibles

- Numérotation globale des events (ID unique) au lieu de repartir à 1 par producteur.
- Persistance append (au lieu d'écraser) et rotation du fichier.
- Exposer paramètres (nombre d'événements, délai, timeout) via flags CLI.
- Ajouter tests unitaires et utiliser `go test` pour functions isolées.
- Rendre le pipeline reentrant/reconnectable (fichiers de config).

Notes pratiques

- `events.json` est créé dans le répertoire courant. Lance toujours les commandes depuis le dossier du projet pour retrouver facilement le fichier.
- Pour détecter race conditions sur la version non protégée :

```powershell
go run -race main.go race
```

Auteur  
Projet réalisé seul à des fins pédagogiques (Cours 3).  
A propos de l'auteur  
TP3 réalisé dans le cadre de l'initiation à Go. Réalisé par Fairytale-Dev(Farid-Efrei) (étudiant — Alternant à l'Efrei).

🦋 Paix à tous 🦋
