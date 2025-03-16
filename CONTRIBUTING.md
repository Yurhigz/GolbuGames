# Guide de contribution à GolbuGames

## Gestion des branches

1. **Créer une nouvelle branche**
```bash
# Depuis main
git checkout -b feature/nom-de-la-feature
# ou
git checkout -b fix/nom-du-fix
```

2. **Nomenclature des branches**
- `feature/*` : Nouvelles fonctionnalités
- `fix/*` : Corrections de bugs
- `docs/*` : Documentation
- `test/*` : Ajout ou modification de tests
- `refactor/*` : Refactoring de code existant

## Commits

1. **Format des messages**
```
type(scope): description courte

Description détaillée si nécessaire
```

2. **Types de commit**
- `feat`: Nouvelle fonctionnalité
- `fix`: Correction de bug
- `docs`: Documentation
- `test`: Tests
- `refactor`: Refactoring
- `style`: Formatage, point-virgules manquants, etc.
- `chore`: Maintenance générale

## Processus de Merge Request

1. **Avant de créer une MR**
- Assurez-vous d'effectuer des tests préalables, qu'ils soient positifs et de les incorporer dans la MR
- Mettez à jour votre branche avec dev afin d'éviter de MR des changements passés
```bash
git checkout dev
git pull
git checkout votre-branche
git merge dev
```

2. **Créer la Merge Request**
- Titre clair et descriptif
- Description détaillée des changements (Si nécessaire mais pas obligatoire)

3. **Template de MR**
```markdown
## Description
[Description des changements]

## Type de changement
- [ ] Nouvelle fonctionnalité
- [ ] Correction de bug
- [ ] Documentation
- [ ] Autre (préciser)

## Tests
- [ ] Tests unitaires ajoutés/modifiés
- [ ] Tests manuels effectués

## Checklist
- [ ] Code commenté
- [ ] Documentation mise à jour
- [ ] Tests passent
```

## Tests

- Tous les tests doivent passer avant de soumettre une MR
- Ajouter des tests pour les nouvelles fonctionnalités
```bash
go test ./...
```

## Documentation

- Mettre à jour la documentation si nécessaire
- Commenter le code complexe
- Ajouter des exemples d'utilisation (Dans la mesure où c'est un cas de figure complexe sinon skip)

## Points importants

1. **Ne jamais commit directement sur main**
2. **Garder les commits cohérents**
3. **Suivre les conventions de code Go (bonnes pratiques dans la mesure du possible)**
4. **Tester localement avant de push (fichiers test.go)**

## Code Review

- Au moins un reviewer doit approuver
- De préférence revoir ensemble le code mais sinon un commentaire sur la MR peut faire l'affaire
- Corriger en fonction du commentaire