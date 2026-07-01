╔══════════════════════════════════════════╗
║              PPToIMG — macOS             ║
╚══════════════════════════════════════════╝

Extrait toutes les images d'un fichier PowerPoint (.pptx)
ou PDF (.pdf) et les dépose dans un dossier sur votre Bureau.


── PREMIÈRE OUVERTURE (obligatoire une seule fois) ────────

  L'application n'étant pas signée par un compte développeur
  Apple, macOS affiche un avertissement de sécurité au tout
  premier lancement : "PPToIMG non ouvert".

  Pour débloquer l'app (une seule fois) :

  1. Ouvrez l'app Terminal (Spotlight → tapez "Terminal")
  2. Tapez la commande suivante puis Entrée
     (adaptez le chemin si PPToIMG.app n'est pas dans
     le dossier Téléchargements) :

       xattr -cr ~/Downloads/PPToIMG.app

  3. Relancez PPToIMG.app normalement — l'avertissement
     ne réapparaîtra plus.


── UTILISATION ────────────────────────────────────────────

  1. Glissez un fichier .pptx ou .pdf directement sur l'icône
     de PPToIMG.app (dans le Finder ou le Dock)
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

  1. Glissez PPToIMG.app dans votre Dock
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


──────────────────────────────────────────────────────────
