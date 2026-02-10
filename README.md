# 🎭 Bashcord Installer - Parce que Discord Vanilla, c'est Chiant

> *"Discord mais en mieux, ou comment transformer ton client de chat en œuvre d'art"* 💅

## 🤔 Qu'est-ce que c'est que ce bordel ?

Bashcord Installer est l'outil ultime pour installer [Bashcord](https://github.com/BashOnZsh/Bashcord), le mod Discord le plus stylé de l'univers. Parce que franchement, utiliser Discord sans mod en 2024, c'est comme manger des pâtes sans sel... techniquement possible, mais pourquoi se faire du mal ?

**🐟 Et surtout, on kiffe grave le skin poiscaille !** Parce que quoi de mieux qu'un Discord qui ressemble à un poisson stylé ? C'est la classe ultime, et ça rend ton chat plus fun qu'un aquarium ! 🐠✨

### 🎯 Pourquoi Bashcord ?
- ✨ **Interface sarcastique** - Parce qu'on a tous besoin d'un peu d'humour dans nos vies
- 🐟 **Skin poiscaille stylé** - Le plus beau skin Discord qui existe, point final !
- 🚀 **Installation rapide** - Plus rapide que le temps qu'il faut pour expliquer à ta grand-mère ce qu'est Discord
- 🛡️ **Sécurisé** - Plus sûr que tes mots de passe "123456"
- 🎨 **Personnalisable** - Parce que chacun mérite son propre style, même les développeurs

## 📥 Téléchargements - Choisis ton Poison

### 🪟 Windows (Pour les Masochistes)
- [🎨 Version GUI](https://github.com/BashOnZsh/Bashotl/releases/latest/download/Bashcord.exe) - *Pour ceux qui aiment cliquer*
- [⌨️ Version CLI](https://github.com/BashOnZsh/Bashotl/releases/latest/download/Bashcord-cli.exe) - *Pour les vrais hackers*

### 🍎 MacOS (Pour les Hipsters)
- [🎨 Version GUI](https://github.com/BashOnZsh/Bashotl/releases/latest/download/Bashcord.MacOS.zip) - *Élégant comme un MacBook à 3000€*

### 🐧 Linux (Pour les Illuminés)
- [🎨 Version GUI X11](https://github.com/BashOnZsh/Bashotl/releases/latest/download/Bashcord-x11) - *Old school mais efficace*
- [⌨️ Version CLI](https://github.com/BashOnZsh/Bashotl/releases/latest/download/Bashcord-Linux) - *Parce que les vrais utilisent le terminal*

### 🤷 Pas sûr de ta version ?

*Pas de panique, même les pros se trompent parfois !*

- **Windows** : Si tu vois des dossiers comme `C:\` et que tu pleures un peu chaque jour, c'est pour toi
- **MacOS** : Si tu as payé plus de 2000€ pour ton ordinateur et que tu te sens supérieur, c'est toi
- **Linux** : Si tu passes plus de temps à configurer ton système qu'à l'utiliser, bienvenue au club

**GUI vs CLI ?**
- **GUI** : Interface graphique pour les humains normaux
- **CLI** : Terminal pour ceux qui codent en binaire dans leur tête

### 🆘 Besoin d'aide ?

Tu n'arrives pas à télécharger ? Ton ordi fait des bruits bizarres ? Le fichier refuse de s'ouvrir ?

➡️ **Contacte `miserablepersonne` sur Discord** - *Il ne mord pas (enfin, pas tout le temps)*

> **Conseil de pro** : Préfère toujours les téléchargements depuis GitHub Releases. Les autres sites, c'est comme les bonbons d'un inconnu - ça peut mal finir.

## 🛠️ Compilation - Pour les Courageux

*"Ah, tu veux compiler toi-même ? Respect, voici comment ne pas tout casser..."*

### 🔧 Prérequis (Ou Comment Préparer ton Calvaire)

Tu auras besoin de :
- [Go](https://go.dev/doc/install) - *Le langage, pas le jeu de société*
- GCC - *MinGW sur Windows, parce que Microsoft aime compliquer*

<details>
<summary>🐧 Dépendances Linux (Clique si tu es assez fou pour utiliser Linux)</summary>

#### Dépendances de base (Le minimum syndical)
```bash
# Ubuntu/Debian (Pour les débutants)
apt install -y pkg-config libsdl2-dev libglx-dev libgl1-mesa-dev

# Fedora/RHEL (Pour les rebelles)
dnf install pkg-config libGL-devel libXxf86vm-devel
```

#### Dépendances X11 (L'ancêtre qui refuse de mourir)
```bash
# Ubuntu/Debian
apt install -y xorg-dev

# Fedora/RHEL
dnf install libXcursor-devel libXi-devel libXinerama-devel libXrandr-devel
```

#### Dépendances Wayland (Le futur, paraît-il)
```bash
# Ubuntu/Debian
apt install -y libwayland-dev libxkbcommon-dev wayland-protocols extra-cmake-modules

# Fedora/RHEL
dnf install wayland-devel libxkbcommon-devel wayland-protocols-devel extra-cmake-modules
```

</details>

### 🏗️ Construction (Ou l'Art de Transformer du Code en Miracle)

#### 1. Installer les dépendances
```bash
go mod tidy
# Prie pour que ça marche du premier coup
```

#### 2. Compiler la version GUI
```bash
# Windows/Mac/Linux X11 (Le trio classique)
go build

# Linux Wayland (Pour les avant-gardistes)
go build --tags wayland
```

#### 3. Compiler la version CLI
```bash
go build --tags cli
# Parce que parfois, moins c'est plus
```

> 💡 **Astuce de Pro** : Regarde [notre workflow GitHub](https://github.com/BashOnZsh/Bashotl/blob/main/.github/workflows/release.yml) pour les flags de compilation optimaux. Ou pas, fais comme tu veux, c'est ta vie après tout.

## 🎭 Fonctionnalités Exclusives

- 🎪 **Messages sarcastiques** - Parce que l'installation doit être divertissante
- 🎯 **Interface française** - Oui, on parle la langue de Molière ici
- 🛡️ **Détection automatique** - Trouve Discord même s'il se cache
- 🎨 **Personnalisation** - Ton Discord, tes règles
- 🚀 **Mise à jour automatique** - Parce qu'on n'a pas que ça à faire

## 🤝 Contribution - Rejoins la Révolution

Tu veux contribuer ? Fantastique ! Voici comment :

1. 🍴 Fork le projet (comme si tu volais une recette)
2. 🌿 Crée une branche (`git checkout -b feature/ma-super-feature`)
3. 💾 Commit tes changements (`git commit -m 'Ajout de ma super feature'`)
4. 📤 Push ta branche (`git push origin feature/ma-super-feature`)
5. 🎯 Ouvre une Pull Request et prie

## 📜 Licence

Ce projet est sous licence [MIT](LICENSE) - Fais-en ce que tu veux, mais ne nous blame pas si ça explose.

## 🙏 Remerciements

- À l'équipe [Equicord](https://github.com/Equicord/Equicord) pour le mod génial
- À Discord pour avoir créé quelque chose d'assez bien pour qu'on veuille le modifier
- À toi, utilisateur courageux, qui lit jusqu'au bout

---

<div align="center">

**Fait avec 💜 par un développeur de merde dc: miserablepersonne .**

*"Parce que la vie est trop courte pour utiliser Discord de manière basique"*

</div>
