# 🌤 Application Météo en Go avec Fyne

Cette application de bureau, développée en **Go** avec le framework **Fyne**, permet d'afficher la météo actuelle et les prévisions pour plusieurs jours, en utilisant l'API **OpenWeatherMap**.

---

## ✨ Fonctionnalités

- Affichage de la **température actuelle**
- Description des conditions météo
- Icône météo dynamique
- **Prévisions sur 5 jours** avec températures et icônes
- Interface graphique responsive et multiplateforme (Windows, Linux, macOS)

---

## 📦 Installation

### 1. Prérequis
- Go 1.18 ou plus récent
- Une clé API gratuite [OpenWeatherMap](https://openweathermap.org/api)

### 2. Cloner le dépôt
```bash
git clone https://github.com/Ndxdeveloper/meteo-fyne.git
cd meteo-fyne
````

### 3. Installer les dépendances

```bash
go mod tidy
```

### 4. Configurer la clé API

Modifier dans `main.go` :

```go
const apiKey = "VOTRE_CLE_API_OPENWEATHERMAP"
const city   = "VotreVille"
```

### 5. Lancer l'application

```bash
go run main.go
```

---

## 🖼 Aperçu

![Aperçu de l'application](screenshot.png)

---

## ⚙️ Technologies utilisées

* **Go**
* **Fyne** (interface graphique)
* **OpenWeatherMap API**

---

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

© 2025 - [NDXdev](https://github.com/ndxDeveloper) - [ndxdev@gmail.com](mailto:ndxdev@gmail.com)
