# PPToIMG

Extrait les images embarquées d'un fichier PowerPoint (`.pptx`) ou PDF (`.pdf`), sans avoir besoin d'ouvrir le fichier dans un logiciel.

## Utilisation

- **macOS** : glissez un `.pptx` ou `.pdf` sur l'icône de `PPToIMG.app` (Dock ou Finder). Un dossier `NomDuFichier-images` apparaît sur le Bureau.
- **Windows** : double-cliquez sur `PPToIMG.exe`, sélectionnez le fichier dans la fenêtre qui s'ouvre.

Voir les README dans `Mac Os/` et `Windows/` pour le détail.

## Mise à jour

Les deux applications vérifient automatiquement s'il existe une nouvelle version publiée dans les [Releases](../../releases) de ce repo, et proposent le téléchargement le cas échéant.

## Structure du repo

```
_icon/              Icône source (png)
Mac Os/              Application macOS compilée (.app) + README
Windows/             Exécutable Windows compilé (.exe) + README
src/macos/           Code source AppleScript
src/windows/         Code source Go (compilé en .exe via cross-compilation)
```

## Compiler depuis les sources

**Windows (.exe)** — depuis macOS ou Linux, avec Go installé :

```bash
cd src/windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-H windowsgui -s -w" -o PPToIMG.exe .
```

**macOS (.app)** — avec `osacompile` (inclus dans macOS) :

```bash
osacompile -o PPToIMG.app src/macos/PPToIMG.applescript
```

## Publier une nouvelle version

1. Mettre à jour la constante de version dans `src/windows/main.go` (`VERSION`) et `src/macos/PPToIMG.applescript` (`current_version`)
2. Compiler les deux binaires
3. Créer une Release GitHub taguée `vX.Y.Z` avec les binaires attachés — c'est ce tag que les apps comparent à leur version locale
