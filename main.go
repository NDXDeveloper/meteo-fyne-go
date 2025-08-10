// main.go
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	// Import des packages Fyne nécessaires pour créer l'interface graphique
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	// Package interne pour la gestion de la météo
	"meteo/weather"
)

const (
	// Clé API OpenWeatherMap + ville ciblée
	apiKey = ""
	city   = ""
)

func main() {
	// Vérifie si la clé API a été remplacée
	if apiKey == "" {
		log.Fatal("Erreur : Veuillez remplacer la clé API.")
	}

	// --- 1. Création de l'application Fyne ---
	// app.New() crée une instance de l'application graphique
	myApp := app.New()

	// myApp.NewWindow() crée une nouvelle fenêtre avec un titre dynamique
	myWindow := myApp.NewWindow(fmt.Sprintf("Météo pour %s", city))
	// On fixe la taille initiale de la fenêtre
	myWindow.Resize(fyne.NewSize(600, 400))

	// --- 2. Initialisation du client météo ---
	weatherClient := weather.NewClient(apiKey)

	// --- 3. Widgets pour la météo actuelle ---
	// Image pour l’icône météo (vide au départ)
	currentIcon := canvas.NewImageFromResource(nil)
	currentIcon.SetMinSize(fyne.NewSize(100, 100))
	currentIcon.FillMode = canvas.ImageFillContain

	// Température (en gras)
	currentTemp := widget.NewLabel("Chargement...")
	currentTemp.TextStyle.Bold = true

	// Description météo (ex : "Nuageux")
	currentDesc := widget.NewLabel("")

	// Conteneur vertical pour la météo actuelle
	currentWeatherBox := container.NewVBox(
		// Titre : nom de la ville en gras + italique, centré
		widget.NewLabelWithStyle(city, fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true}),
		// Ligne horizontale : icône à gauche, température + description à droite
		container.NewHBox(
			currentIcon,
			container.NewVBox(
				currentTemp,
				currentDesc,
			),
		),
	)

	// --- 4. Widgets pour les prévisions météo ---
	// Grille avec 5 colonnes (un jour par colonne)
	forecastContainer := container.NewGridWithColumns(5)

	// --- 5. Conteneur principal ---
	// On empile verticalement : météo actuelle, séparateur, prévisions
	content := container.NewVBox(
		currentWeatherBox,
		widget.NewSeparator(),
		forecastContainer,
	)

	// --- 6. Lancement de la récupération des données météo ---
	// On utilise "go" pour exécuter la mise à jour dans une goroutine
	go updateWeatherData(weatherClient, currentIcon, currentTemp, currentDesc, forecastContainer)

	// --- 7. Affichage ---
	// On définit le contenu de la fenêtre
	myWindow.SetContent(content)
	// On lance l'application (bloquant tant que la fenêtre est ouverte)
	myWindow.ShowAndRun()
}

func updateWeatherData(client *weather.Client, icon *canvas.Image, temp, desc *widget.Label, forecastBox *fyne.Container) {
	// --- 1. Récupération météo actuelle ---
	current, err := client.GetCurrentWeather(city)
	if err != nil {
		log.Printf("Erreur météo actuelle : %v", err)
		// fyne.Do permet d'exécuter du code sur le thread graphique principal
		fyne.Do(func() {
			desc.SetText("Erreur de chargement")
		})
		return
	}

	// --- 2. Mise à jour des widgets (météo actuelle) ---
	fyne.Do(func() {
		// Affiche la température formatée avec 1 décimale
		temp.SetText(fmt.Sprintf("%.1f °C", current.Main.Temp))

		// Si on a des infos météo...
		if len(current.Weather) > 0 {
			description := strings.Title(current.Weather[0].Description)
			desc.SetText(description)

			// On récupère l'icône météo
			iconData, err := client.GetIcon(current.Weather[0].Icon)
			if err == nil {
				// On crée une ressource Fyne à partir des données binaires
				icon.Resource = fyne.NewStaticResource(current.Weather[0].Icon+".png", iconData)
				// icon.Refresh() redessine l'image
				icon.Refresh()
			}
		}
	})

	// --- 3. Récupération des prévisions ---
	forecast, err := client.GetForecast(city)
	if err != nil {
		log.Printf("Erreur prévisions : %v", err)
		return
	}

	// --- 4. Affichage des prévisions (5 jours max) ---
	for i, f := range forecast.List {
		if i >= 5 {
			break
		}

		// Jour (format abrégé ex : "lun. 12")
		dayLabel := widget.NewLabel(time.Unix(f.Dt, 0).Format("lun. 02"))
		dayLabel.TextStyle.Bold = true

		// Température du jour
		tempLabel := widget.NewLabel(fmt.Sprintf("%.0f °C", f.Main.Temp))

		// Icône météo du jour
		forecastIcon := canvas.NewImageFromResource(nil)
		forecastIcon.SetMinSize(fyne.NewSize(50, 50))
		forecastIcon.FillMode = canvas.ImageFillContain

		// Si on a une icône, on la charge dans une goroutine
		if len(f.Weather) > 0 {
			iconCode := f.Weather[0].Icon
			go func(iconCode string, forecastIcon *canvas.Image) {
				iconData, err := client.GetIcon(iconCode)
				if err == nil {
					fyne.Do(func() {
						forecastIcon.Resource = fyne.NewStaticResource(iconCode+".png", iconData)
						forecastIcon.Refresh()
					})
				}
			}(iconCode, forecastIcon)
		}

		// On ajoute la prévision au conteneur (dans le thread principal)
		fyne.Do(func() {
			dayBox := container.NewVBox(
				dayLabel,
				forecastIcon,
				tempLabel,
			)
			forecastBox.Add(dayBox)
			forecastBox.Refresh()
		})
	}
}
