╔══════════════════════════════════════════╗
║              PPToIMG — macOS             ║
╚══════════════════════════════════════════╝

Extrait toutes les images d'un fichier PowerPoint (.pptx)
ou PDF (.pdf) et les dépose dans un dossier sur votre Bureau.

Le nom de l'application (PPToIMG-vX.Y.Z.app) contient toujours
son numéro de version, pour vérifier d'un coup d'œil que vous
avez la dernière copie.


── PREMIÈRE OUVERTURE (obligatoire une seule fois) ────────

  L'application n'étant pas signée par un compte développeur
  Apple, macOS peut afficher un avertissement de sécurité au
  premier lancement : "PPToIMG non ouvert".

  Cet avertissement n'apparaît QUE si le fichier a été
  téléchargé via un navigateur (macOS marque alors le
  fichier comme "provenant d'internet"). Deux façons de
  l'éviter ou de le lever :

  ▸ MÉTHODE A — Éviter le blocage dès le départ
    Récupérez PPToIMG-vX.Y.Z.app par copie de fichier plutôt que
    par téléchargement navigateur : disque réseau partagé,
    clé USB, ou dossier synchronisé par un client Dropbox /
    Google Drive / OneDrive installé sur le Mac (pas via
    leur site web). Dans ce cas, aucun avertissement
    n'apparaît et l'app s'ouvre directement.

    ⚠️ AirDrop ne fonctionne PAS pour éviter le blocage :
    macOS marque aussi les fichiers reçus par AirDrop.

  ▸ MÉTHODE B — Si l'avertissement apparaît quand même
    1. Ouvrez l'app Terminal (Spotlight → tapez "Terminal")
    2. Tapez la commande suivante puis Entrée
       (adaptez le chemin si PPToIMG-vX.Y.Z.app n'est pas dans
       le dossier Téléchargements) :

         xattr -cr ~/Downloads/PPToIMG-vX.Y.Z.app

    3. Relancez PPToIMG-vX.Y.Z.app normalement — l'avertissement
       ne réapparaîtra plus.


── UTILISATION ────────────────────────────────────────────

  1. Glissez un fichier .pptx ou .pdf directement sur l'icône
     de PPToIMG-vX.Y.Z.app (dans le Finder ou le Dock)
  2. Le dossier "NomDuFichier-images" apparaît sur votre
     Bureau et s'ouvre automatiquement dans le Finder

  Astuce : cliquez simplement sur l'icône (sans glisser de
  fichier) pour vérifier la version installée et si une
  mise à jour est disponible.


── RÉSULTAT ───────────────────────────────────────────────

  Un dossier est créé sur votre Bureau :

    Bureau/
    └── NomDuFichier-images/
        ├── image1.png
        ├── image2.jpg
        └── ...

  Le dossier contient toutes les images et médias
  intégrés dans la présentation (photos, logos, icônes…).


── ASTUCE — UTILISATION DEPUIS LE DOCK ────────────────────

  Pour un accès encore plus rapide :

  1. Glissez PPToIMG-vX.Y.Z.app dans votre Dock
  2. Depuis n'importe où sur votre Mac, glissez un fichier
     .pptx sur l'icône PPToIMG dans le Dock
  3. Relâchez — c'est fait !

  Plus besoin d'ouvrir de fenêtre Finder, ça fonctionne
  même si le fichier est sur un autre écran ou dans
  une autre application.


── REMARQUES ──────────────────────────────────────────────

  - Aucune installation requise, aucun logiciel tiers nécessaire.
  - PowerPoint n'a pas besoin d'être installé sur le Mac.
  - Fonctionne avec les fichiers .pptx (PowerPoint 2007+) et .pdf.
  - L'app vérifie automatiquement les mises à jour au lancement.
    En cas de nouvelle version, cliquer sur "Télécharger"
    demande où enregistrer la nouvelle application (nommée
    avec son numéro de version), la télécharge et la prépare
    automatiquement (aucun avertissement de sécurité à ce
    stade, contrairement à un téléchargement via navigateur).
    Le Finder affiche ensuite la nouvelle application, prête
    à remplacer l'ancienne — il suffit de fermer PPToIMG puis
    de faire l'échange.


──────────────────────────────────────────────────────────
