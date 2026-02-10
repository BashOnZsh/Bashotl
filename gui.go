//go:build !cli

/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	"vencord/buildinfo"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"

	// png decoder for icon
	_ "image/png"
	"os"
	"os/signal"
	path "path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

var (
	discords        []any
	radioIdx        int
	customChoiceIdx int

	customDir              string
	autoCompleteDir        string
	autoCompleteFile       string
	autoCompleteCandidates []string
	autoCompleteIdx        int
	lastAutoComplete       string
	didAutoComplete        bool

	modalId      = 0
	modalTitle   = "Oh Non :("
	modalMessage = "Vous ne devriez jamais voir ceci"

	acceptedOpenAsar   bool
	showedUpdatePrompt bool

	// Nouvelles variables pour les fonctionnalités avancées
	currentTheme      = "fishstick" // fishstick, dark, skullkid, sanglant, terminal, pepe, wumpus
	showAdvancedMode  = false
	autoUpdateEnabled = true
	showNotifications = true
	compactMode       = false
	animationEnabled  = true

	// Variables pour les statistiques
	installCount    = 0
	lastInstallTime = ""
	preferredBranch = "auto"

	win *g.MasterWindow

	// Variable pour le lecteur audio
	audioStarted          = false
	audioVolume   float64 = 0.05 // Volume entre 0.0 et 1.0 (réglé très bas au lancement ~5%)
	maxVolume     float64 = 0.3  // Limite maximale du volume (~30% pour éviter d'être trop fort)
	audioContext  *oto.Context
	audioPlayer   *oto.Player
	audioLoopDone = make(chan bool, 1)
)

//go:embed assets/icon_256.png
var iconBytes []byte

//go:embed bashcord.mp3
var bashcordMP3 []byte

// Couleurs du thème Fishstick
var (
	FishstickOrange = color.RGBA{R: 0xFF, G: 0x8C, B: 0x42, A: 0xFF} // Orange vif du poisson
	FishstickBlue   = color.RGBA{R: 0x4A, G: 0x90, B: 0xE2, A: 0xFF} // Bleu océan
	FishstickYellow = color.RGBA{R: 0xFF, G: 0xE1, B: 0x35, A: 0xFF} // Jaune des écailles
	FishstickGreen  = color.RGBA{R: 0x2E, G: 0xCC, B: 0x71, A: 0xFF} // Vert algue
	FishstickPurple = color.RGBA{R: 0x9B, G: 0x59, B: 0xB6, A: 0xFF} // Violet profond océan
	FishstickCyan   = color.RGBA{R: 0x1A, G: 0xBC, B: 0x9C, A: 0xFF} // Cyan tropical
	FishstickWhite  = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // Blanc des bulles
	FishstickDark   = color.RGBA{R: 0x1A, G: 0x1A, B: 0x2E, A: 0xFF} // Bleu nuit profond
	FishstickLight  = color.RGBA{R: 0x87, G: 0xCE, B: 0xEB, A: 0xFF} // Bleu ciel océan
)

func init() {
	LogLevel = LevelDebug
	loadUserPreferences()
}

// Fonctions pour la gestion des thèmes
func getThemeColors(theme string) map[string]color.RGBA {
	switch theme {
	case "dark":
		return map[string]color.RGBA{
			"primary":   color.RGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF}, // Noir profond amélioré
			"secondary": color.RGBA{R: 0x2D, G: 0x2D, B: 0x33, A: 0xFF}, // Gris bleuté plus profond
			"accent":    color.RGBA{R: 0x5B, G: 0x6E, B: 0xF7, A: 0xFF}, // Bleu accent plus vif
			"text":      color.RGBA{R: 0xF0, G: 0xF0, B: 0xF5, A: 0xFF}, // Blanc légèrement bleuté
			"success":   color.RGBA{R: 0x43, G: 0x9A, B: 0x68, A: 0xFF}, // Vert succès plus lumineux
			"warning":   color.RGBA{R: 0xFF, G: 0xD9, B: 0x4A, A: 0xFF}, // Jaune warning plus chaud
			"error":     color.RGBA{R: 0xF5, G: 0x4C, B: 0x4F, A: 0xFF}, // Rouge error plus vif
		}
	case "fishstick":
		return map[string]color.RGBA{
			"primary":   color.RGBA{R: 0x15, G: 0x1A, B: 0x28, A: 0xFF}, // Bleu nuit plus profond
			"secondary": color.RGBA{R: 0x2E, G: 0x4A, B: 0x7C, A: 0xFF}, // Bleu océan plus riche
			"accent":    color.RGBA{R: 0xFF, G: 0x95, B: 0x52, A: 0xFF}, // Orange plus chaud
			"text":      color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, // Blanc pur
			"success":   color.RGBA{R: 0x3E, G: 0xD4, B: 0x7A, A: 0xFF}, // Vert succès plus lumineux
			"warning":   color.RGBA{R: 0xFF, G: 0xE5, B: 0x4F, A: 0xFF}, // Jaune plus vibrant
			"error":     color.RGBA{R: 0xB8, G: 0x4D, B: 0xD6, A: 0xFF}, // Violet error plus doux
		}
	case "skullkid":
		return map[string]color.RGBA{
			"primary":   color.RGBA{R: 0x12, G: 0x08, B: 0x24, A: 0xFF}, // Violet très foncé amélioré
			"secondary": color.RGBA{R: 0x3A, G: 0x24, B: 0x5A, A: 0xFF}, // Violet moyen plus riche
			"accent":    color.RGBA{R: 0x9B, G: 0x55, B: 0xB8, A: 0xFF}, // Violet accent plus lumineux
			"text":      color.RGBA{R: 0xFF, G: 0xE8, B: 0x5C, A: 0xFF}, // Jaune doré plus brillant
			"success":   color.RGBA{R: 0x6B, G: 0xA0, B: 0x4A, A: 0xFF}, // Vert forêt plus lumineux
			"warning":   color.RGBA{R: 0xFF, G: 0xA5, B: 0x5C, A: 0xFF}, // Orange plus chaud
			"error":     color.RGBA{R: 0xB8, G: 0x28, B: 0x3A, A: 0xFF}, // Rouge sombre plus profond
		}
	case "sanglant":
		return map[string]color.RGBA{
			"primary":   color.RGBA{R: 0x15, G: 0x00, B: 0x00, A: 0xFF}, // Noir rougeâtre plus profond
			"secondary": color.RGBA{R: 0x3A, G: 0x0A, B: 0x0A, A: 0xFF}, // Rouge très foncé avec nuance
			"accent":    color.RGBA{R: 0xC8, G: 0x1A, B: 0x1A, A: 0xFF}, // Rouge accent plus vibrant
			"text":      color.RGBA{R: 0xFF, G: 0xF0, B: 0xF0, A: 0xFF}, // Blanc rosé plus doux
			"success":   color.RGBA{R: 0x5C, G: 0x1A, B: 0x1A, A: 0xFF}, // Rouge foncé amélioré
			"warning":   color.RGBA{R: 0xFF, G: 0x6B, B: 0x6B, A: 0xFF}, // Rouge vif plus lumineux
			"error":     color.RGBA{R: 0xD9, G: 0x1A, B: 0x1A, A: 0xFF}, // Rouge intense plus chaud
		}
	case "terminal":
		return map[string]color.RGBA{
			"primary":   color.RGBA{R: 0x0A, G: 0x0A, B: 0x0A, A: 0xFF}, // Noir légèrement adouci
			"secondary": color.RGBA{R: 0x1A, G: 0x1F, B: 0x1A, A: 0xFF}, // Gris très foncé avec nuance verte
			"accent":    color.RGBA{R: 0x32, G: 0xFF, B: 0x32, A: 0xFF}, // Vert terminal plus doux
			"text":      color.RGBA{R: 0x8F, G: 0xFF, B: 0x8F, A: 0xFF}, // Vert terminal plus doux pour le texte
			"success":   color.RGBA{R: 0x5C, G: 0xFF, B: 0x5C, A: 0xFF}, // Vert succès plus doux
			"warning":   color.RGBA{R: 0xFF, G: 0xFF, B: 0x5C, A: 0xFF}, // Jaune plus doux
			"error":     color.RGBA{R: 0xFF, G: 0x5C, B: 0x5C, A: 0xFF}, // Rouge plus doux
		}
	case "pepe":
		return map[string]color.RGBA{
			"primary":   color.RGBA{R: 0x1F, G: 0x3D, B: 0x15, A: 0xFF}, // Vert foncé plus profond
			"secondary": color.RGBA{R: 0x3A, G: 0x6B, B: 0x28, A: 0xFF}, // Vert moyen plus riche
			"accent":    color.RGBA{R: 0x7A, G: 0xC8, B: 0x4A, A: 0xFF}, // Vert clair plus vibrant
			"text":      color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, // Blanc pur
			"success":   color.RGBA{R: 0x5A, G: 0xA0, B: 0x3C, A: 0xFF}, // Vert succès plus lumineux
			"warning":   color.RGBA{R: 0xFF, G: 0xE0, B: 0x32, A: 0xFF}, // Jaune plus chaud
			"error":     color.RGBA{R: 0xD9, G: 0x4A, B: 0x1A, A: 0xFF}, // Rouge plus chaud
		}
	case "wumpus":
		return map[string]color.RGBA{
			"primary":   color.RGBA{R: 0x2C, G: 0x2F, B: 0x33, A: 0xFF}, // Gris Discord foncé amélioré
			"secondary": color.RGBA{R: 0x3E, G: 0x43, B: 0x4B, A: 0xFF}, // Gris Discord moyen plus profond
			"accent":    color.RGBA{R: 0x5B, G: 0x6E, B: 0xF7, A: 0xFF}, // Bleu Discord plus vif
			"text":      color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, // Blanc pur
			"success":   color.RGBA{R: 0x3C, G: 0xB8, B: 0x6A, A: 0xFF}, // Vert Discord plus lumineux
			"warning":   color.RGBA{R: 0xFF, G: 0xC9, B: 0x4A, A: 0xFF}, // Jaune Discord plus chaud
			"error":     color.RGBA{R: 0xF5, G: 0x4C, B: 0x4F, A: 0xFF}, // Rouge Discord plus vif
		}
	default:
		return getThemeColors("fishstick")
	}
}

func loadUserPreferences() {
	// TODO: Charger les préférences depuis un fichier de configuration
	// Pour l'instant, on utilise les valeurs par défaut
}

func saveUserPreferences() {
	// TODO: Sauvegarder les préférences dans un fichier de configuration
}

// Fonction pour démarrer la musique en arrière-plan
func startBackgroundMusic() {
	if audioStarted {
		return
	}

	go func() {
		// Décoder le MP3 depuis les bytes embarqués
		decoder, err := mp3.NewDecoder(bytes.NewReader(bashcordMP3))
		if err != nil {
			Log.Warn("Failed to decode MP3", err)
			return
		}

		// Créer le contexte audio
		op := &oto.NewContextOptions{
			SampleRate:   44100,
			ChannelCount: 2,
			Format:       oto.FormatSignedInt16LE,
		}
		ctx, ready, err := oto.NewContext(op)
		if err != nil {
			Log.Warn("Failed to create audio context", err)
			return
		}
		<-ready
		audioContext = ctx

		audioStarted = true
		Log.Debug("Background music started")

		// Boucle de lecture infinie
		for audioStarted {
			// Réinitialiser le décoder pour chaque boucle
			decoder, err = mp3.NewDecoder(bytes.NewReader(bashcordMP3))
			if err != nil {
				Log.Warn("Failed to decode MP3 in loop", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Créer un nouveau player pour cette lecture
			player := ctx.NewPlayer(decoder)
			audioPlayer = player

			// Appliquer le volume
			player.SetVolume(audioVolume)

			// Jouer la musique
			player.Play()

			// Attendre la fin de la lecture
			for player.IsPlaying() && audioStarted {
				// Mettre à jour le volume en temps réel
				if audioPlayer != nil {
					audioPlayer.SetVolume(audioVolume)
				}
				time.Sleep(100 * time.Millisecond)
			}

			// Fermer le player
			if err := player.Close(); err != nil {
				Log.Warn("Failed to close player", err)
			}

			if !audioStarted {
				break
			}
		}

		// Le contexte audio n'a pas besoin d'être fermé explicitement
		audioContext = nil
		audioPlayer = nil
	}()
}

// Fonction pour arrêter la musique et nettoyer
func stopBackgroundMusic() {
	if !audioStarted {
		return
	}

	audioStarted = false

	// Arrêter le player
	if audioPlayer != nil {
		audioPlayer.Close()
		audioPlayer = nil
	}

	// Le contexte audio se fermera automatiquement
	audioContext = nil

	Log.Debug("Background music stopped")
}

func main() {
	InitGithubDownloader()
	discords = FindDiscords()

	customChoiceIdx = len(discords)

	go func() {
		<-GithubDoneChan
		g.Update()
	}()

	go func() {
		<-SelfUpdateCheckDoneChan
		g.Update()
	}()

	var linuxFlags g.MasterWindowFlags = 0
	if runtime.GOOS == "linux" {
		os.Setenv("GDK_SCALE", "1")
		os.Setenv("GDK_DPI_SCALE", "1")
	}

	// Créer la fenêtre avec une taille raisonnable (non plein écran)
	win = g.NewMasterWindow("Bashcord Installer", 1000, 700, linuxFlags)

	icon, _, err := image.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		Log.Warn("Failed to load application icon", err)
		Log.Debug(iconBytes, len(iconBytes))
	} else {
		win.SetIcon([]image.Image{icon})
	}

	// Démarrer la musique en arrière-plan
	startBackgroundMusic()

	// Gérer les signaux de fermeture pour nettoyer la musique
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		stopBackgroundMusic()
		os.Exit(0)
	}()

	win.Run(loop)

	// Nettoyer quand la fenêtre se ferme
	stopBackgroundMusic()
}

type CondWidget struct {
	predicate  bool
	ifWidget   func() g.Widget
	elseWidget func() g.Widget
}

func (w *CondWidget) Build() {
	if w.predicate {
		w.ifWidget().Build()
	} else if w.elseWidget != nil {
		w.elseWidget().Build()
	}
}

func getChosenInstall() *DiscordInstall {
	var choice *DiscordInstall
	if radioIdx == customChoiceIdx {
		choice = ParseDiscord(customDir, "")
		if choice == nil {
			g.OpenPopup("#invalid-custom-location")
		}
	} else {
		choice = discords[radioIdx].(*DiscordInstall)
	}
	return choice
}

func InstallLatestBuilds() (err error) {
	if IsDevInstall {
		return
	}

	err = installLatestBuilds()
	if err != nil {
		ShowModal("Oups !", "Échec de l'installation des dernières versions de Bashcord depuis GitHub :\n"+err.Error())
	}
	return
}

func handlePatch() {
	choice := getChosenInstall()
	if choice != nil {
		choice.Patch()
	}
}

func handleUnpatch() {
	choice := getChosenInstall()
	if choice != nil {
		choice.Unpatch()
	}
}

func handleOpenAsar() {
	if acceptedOpenAsar || getChosenInstall().IsOpenAsar() {
		handleOpenAsarConfirmed()
		return
	}

	g.OpenPopup("#openasar-confirm")
}

func handleOpenAsarConfirmed() {
	choice := getChosenInstall()
	if choice != nil {
		if choice.IsOpenAsar() {
			if err := choice.UninstallOpenAsar(); err != nil {
				handleErr(choice, err, "désinstaller OpenAsar de")
			} else {
				g.OpenPopup("#openasar-unpatched")
				g.Update()
			}
		} else {
			if err := choice.InstallOpenAsar(); err != nil {
				handleErr(choice, err, "installer OpenAsar sur")
			} else {
				g.OpenPopup("#openasar-patched")
				g.Update()
			}
		}
	}
}

func handleErr(di *DiscordInstall, err error, action string) {
	if errors.Is(err, os.ErrPermission) {
		switch runtime.GOOS {
		case "windows":
			err = errors.New("Permission refusée. Assurez-vous que Discord est complètement fermé (depuis la barre système) !")
		case "darwin":
			// FIXME: This text is not selectable which is a bit mehhh
			command := "sudo chown -R \"${USER}:wheel\" " + di.path
			err = errors.New("Permission refusée. Veuillez accorder à l'installateur l'accès complet au disque dans les paramètres système (page confidentialité et sécurité).\n\nSi cela ne fonctionne toujours pas, essayez d'exécuter la commande suivante dans votre terminal :\n" + command)
		case "linux":
			command := "sudo chown -R \"$USER:$USER\" " + di.path
			err = errors.New("Permission refusée. Essayez d'exécuter l'installateur avec les privilèges sudo.\n\nSi cela ne fonctionne toujours pas, essayez d'exécuter la commande suivante dans votre terminal :\n" + command)
		default:
			err = errors.New("Permission refusée. Essayez peut-être de m'exécuter en tant qu'Administrateur/Root ?")
		}
	}

	ShowModal("Échec de "+action+" cette installation", err.Error())
}

func HandleScuffedInstall() {
	g.OpenPopup("#scuffed-install")
}

func (di *DiscordInstall) Patch() {
	if CheckScuffedInstall() {
		return
	}
	if err := di.patch(); err != nil {
		handleErr(di, err, "patcher")
	} else {
		g.OpenPopup("#patched")
	}
}

func (di *DiscordInstall) Unpatch() {
	if err := di.unpatch(); err != nil {
		handleErr(di, err, "dépatcher")
	} else {
		g.OpenPopup("#unpatched")
	}
}

func onCustomInputChanged() {
	p := customDir
	if len(p) != 0 {
		// Select the custom option for people
		radioIdx = customChoiceIdx
	}

	dir := path.Dir(p)

	isNewDir := strings.HasSuffix(p, "/")
	wentUpADir := !isNewDir && dir != autoCompleteDir

	if isNewDir || wentUpADir {
		autoCompleteDir = dir
		// reset all the funnies
		autoCompleteIdx = 0
		lastAutoComplete = ""
		autoCompleteFile = ""
		autoCompleteCandidates = nil

		// Generate autocomplete items
		files, err := os.ReadDir(dir)
		if err == nil {
			for _, file := range files {
				autoCompleteCandidates = append(autoCompleteCandidates, file.Name())
			}
		}
	} else if !didAutoComplete {
		// reset auto complete and update our file
		autoCompleteFile = path.Base(p)
		lastAutoComplete = ""
	}

	if wentUpADir {
		autoCompleteFile = path.Base(p)
	}

	didAutoComplete = false
}

// go can you give me []any?
// to pass to giu RangeBuilder?
// yeeeeees
// actually returns []string like a boss
func makeAutoComplete() []any {
	input := strings.ToLower(autoCompleteFile)

	var candidates []any
	for _, e := range autoCompleteCandidates {
		file := strings.ToLower(e)
		if autoCompleteFile == "" || strings.HasPrefix(file, input) {
			candidates = append(candidates, e)
		}
	}
	return candidates
}

func makeRadioOnChange(i int) func() {
	return func() {
		radioIdx = i
	}
}

func Tooltip(label string) g.Widget {
	return g.Style().
		SetStyle(g.StyleVarWindowPadding, 10, 8).
		SetStyleFloat(g.StyleVarWindowRounding, 8).
		To(
			g.Tooltip(label),
		)
}

func InfoModal(id, title, description string) g.Widget {
	return RawInfoModal(id, title, description, false)
}

func RawInfoModal(id, title, description string, isOpenAsar bool) g.Widget {
	isDynamic := strings.HasPrefix(id, "#modal") && !strings.Contains(description, "\n")
	return g.Style().
		SetStyle(g.StyleVarWindowPadding, 30, 30).
		SetStyleFloat(g.StyleVarWindowRounding, 12).
		To(
			g.PopupModal(id).
				Flags(g.WindowFlagsNoTitleBar | Ternary(isDynamic, g.WindowFlagsAlwaysAutoResize, 0)).
				Layout(
					g.Align(g.AlignCenter).To(
						g.Style().SetFontSize(30).To(
							g.Label(title),
						),
						g.Style().SetFontSize(20).To(
							g.Label(description).Wrapped(isDynamic),
						),
						&CondWidget{id == "#scuffed-install", func() g.Widget {
							return g.Column(
								g.Dummy(0, 10),
								g.Button("Emmène-moi là !").OnClick(func() {
									// this issue only exists on windows so using Windows specific path is oki
									username := os.Getenv("USERNAME")
									programData := os.Getenv("PROGRAMDATA")
									g.OpenURL("file://" + path.Join(programData, username))
								}).Size(200, 30),
							)
						}, nil},
						g.Dummy(0, 20),
						&CondWidget{isOpenAsar,
							func() g.Widget {
								return g.Row(
									g.Button("Accepter").
										OnClick(func() {
											acceptedOpenAsar = true
											g.CloseCurrentPopup()
										}).
										Size(100, 30),
									g.Button("Annuler").
										OnClick(func() {
											g.CloseCurrentPopup()
										}).
										Size(100, 30),
								)
							},
							func() g.Widget {
								return g.Button("Ok").
									OnClick(func() {
										g.CloseCurrentPopup()
									}).
									Size(100, 30)
							},
						},
					),
				),
		)
}

func UpdateModal() g.Widget {
	return g.Style().
		SetStyle(g.StyleVarWindowPadding, 30, 30).
		SetStyleFloat(g.StyleVarWindowRounding, 12).
		To(
			g.PopupModal("#update-prompt").
				Flags(g.WindowFlagsNoTitleBar | g.WindowFlagsAlwaysAutoResize).
				Layout(
					g.Align(g.AlignCenter).To(
						g.Style().SetFontSize(30).To(
							g.Label("Votre installateur est obsolète !"),
						),
						g.Style().SetFontSize(20).To(
							g.Label(
								"Souhaitez-vous mettre à jour maintenant ?\n\n"+
									"Une fois que vous appuyez sur Mettre à jour maintenant, le nouvel installateur sera automatiquement téléchargé.\n"+
									"L'installateur semblera temporairement ne plus répondre. Attendez simplement !\n"+
									"Une fois la mise à jour terminée, l'installateur se rouvrira automatiquement.\n\n"+
									"Sur MacOS, les mises à jour automatiques ne sont pas prises en charge, il s'ouvrira donc dans le navigateur.",
							),
						),
						g.Row(
							g.Button("Mettre à jour maintenant").
								OnClick(func() {
									if runtime.GOOS == "darwin" {
										g.CloseCurrentPopup()
										g.OpenURL(GetInstallerDownloadLink())
										return
									}

									err := UpdateSelf()
									g.CloseCurrentPopup()

									if err != nil {
										ShowModal("Échec de la mise à jour automatique !", err.Error())
									} else {
										if err = RelaunchSelf(); err != nil {
											ShowModal("Échec du redémarrage automatique ! Veuillez le faire manuellement.", err.Error())
										}
									}
								}).
								Size(150, 30),
							g.Button("Plus tard").
								OnClick(func() {
									g.CloseCurrentPopup()
								}).
								Size(100, 30),
						),
					),
				),
		)
}

func ShowModal(title, desc string) {
	modalTitle = title
	modalMessage = desc
	modalId++
	g.OpenPopup("#modal" + strconv.Itoa(modalId))
}

// Fonction pour créer un bouton stylisé
func createStyledButton(text string, onClick func(), colors map[string]color.RGBA, width, height float32) g.Widget {
	return g.Style().
		SetColor(g.StyleColorButton, colors["accent"]).
		SetColor(g.StyleColorButtonHovered, color.RGBA{
			R: colors["accent"].R,
			G: colors["accent"].G,
			B: colors["accent"].B,
			A: 200,
		}).
		SetColor(g.StyleColorButtonActive, color.RGBA{
			R: colors["accent"].R,
			G: colors["accent"].G,
			B: colors["accent"].B,
			A: 255,
		}).
		SetColor(g.StyleColorText, colors["text"]).
		SetStyleFloat(g.StyleVarFrameRounding, 8).
		SetStyle(g.StyleVarFramePadding, 12, 8).
		To(
			g.Button(text).
				OnClick(onClick).
				Size(width, height),
		)
}

// Fonction pour créer une carte d'information stylisée
func createInfoCard(title, content string, colors map[string]color.RGBA, height float32) g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, colors["secondary"]).
		SetStyleFloat(g.StyleVarAlpha, 0.9).
		SetStyle(g.StyleVarWindowPadding, 15, 15).
		SetStyleFloat(g.StyleVarChildRounding, 12).
		To(
			g.Child().
				Size(g.Auto, height).
				Layout(
					g.Style().
						SetColor(g.StyleColorText, colors["text"]).
						SetFontSize(18).
						To(
							g.Label(title),
						),
					g.Dummy(0, 8),
					g.Style().
						SetColor(g.StyleColorText, color.RGBA{
							R: colors["text"].R,
							G: colors["text"].G,
							B: colors["text"].B,
							A: 200,
						}).
						SetFontSize(14).
						To(
							g.Label(content).Wrapped(true),
						),
				),
		)
}

// Fonction pour créer le header avec statistiques
// Fonction pour créer le switcher de thème
func renderThemeSwitcher(colors map[string]color.RGBA) g.Widget {
	themes := []string{"fishstick", "dark", "skullkid", "sanglant", "terminal", "pepe", "wumpus"}
	var currentIdx int32
	for i, theme := range themes {
		if theme == currentTheme {
			currentIdx = int32(i)
			break
		}
	}

	return g.Style().
		SetColor(g.StyleColorFrameBg, colors["secondary"]).
		SetColor(g.StyleColorFrameBgHovered, colors["accent"]).
		SetColor(g.StyleColorText, colors["text"]).
		SetStyleFloat(g.StyleVarFrameRounding, 8).
		SetStyle(g.StyleVarFramePadding, 8, 8).
		To(
			g.Row(
				g.Label("Thème:"),
				g.Dummy(5, 0),
				g.Combo("##theme", themes[currentIdx], themes, &currentIdx).
					OnChange(func() {
						currentTheme = themes[currentIdx]
						g.Update()
					}).
					Size(150),
			),
		)
}

// Fonction pour créer le contrôle de volume
func renderVolumeControl(colors map[string]color.RGBA) g.Widget {
	volume32 := float32(audioVolume)
	maxVolume32 := float32(maxVolume)
	volumePercent := int(audioVolume * 100.0)
	return g.Style().
		SetColor(g.StyleColorFrameBg, colors["secondary"]).
		SetColor(g.StyleColorFrameBgHovered, colors["accent"]).
		SetColor(g.StyleColorText, colors["text"]).
		SetStyleFloat(g.StyleVarFrameRounding, 8).
		SetStyle(g.StyleVarFramePadding, 8, 8).
		To(
			g.Row(
				g.Label("Volume:"),
				g.Dummy(5, 0),
				g.SliderFloat(&volume32, 0.0, maxVolume32).
					Size(100).
					Label("##volume").
					OnChange(func() {
						// Limiter le volume au maximum autorisé
						if volume32 > maxVolume32 {
							volume32 = maxVolume32
						}
						// Mettre à jour le volume en temps réel
						audioVolume = float64(volume32)
						if audioPlayer != nil {
							audioPlayer.SetVolume(audioVolume)
						}
						g.Update()
					}),
				g.Dummy(5, 0),
				g.Style().
					SetColor(g.StyleColorText, colors["accent"]).
					To(
						g.Label(fmt.Sprintf("%d%%", volumePercent)),
					),
			),
		)
}

func renderHeader(colors map[string]color.RGBA) g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, colors["primary"]).
		SetStyleFloat(g.StyleVarChildRounding, 15).
		SetStyle(g.StyleVarWindowPadding, 20, 20).
		To(
			g.Child().
				Size(g.Auto, 120).
				Layout(
					g.Row(
						g.Style().
							SetColor(g.StyleColorText, colors["accent"]).
							SetFontSize(36).
							To(
								g.Label("BASHCORD"),
							),
						g.Dummy(20, 0),
						renderThemeSwitcher(colors),
						g.Dummy(20, 0),
						renderVolumeControl(colors),
					),
					g.Dummy(0, 10),
					g.Row(
						g.Style().
							SetColor(g.StyleColorText, colors["success"]).
							SetFontSize(15).
							To(
								g.Label("Installations: "+strconv.Itoa(installCount)),
							),
						g.Dummy(20, 0),
						g.Style().
							SetColor(g.StyleColorText, colors["warning"]).
							SetFontSize(15).
							To(
								g.Label("Derniere installation: "+Ternary(lastInstallTime != "", lastInstallTime, "Jamais")),
							),
						g.Dummy(20, 0),
						g.Style().
							SetColor(g.StyleColorText, colors["accent"]).
							SetFontSize(15).
							To(
								g.Label("Branche preferee: "+preferredBranch),
							),
					),
				),
		)
}

// Fonction pour créer le panneau de contrôle avancé
func renderAdvancedPanel(colors map[string]color.RGBA) g.Widget {
	if !showAdvancedMode {
		return g.Dummy(0, 0)
	}

	return g.Style().
		SetColor(g.StyleColorChildBg, colors["secondary"]).
		SetStyleFloat(g.StyleVarChildRounding, 10).
		SetStyle(g.StyleVarWindowPadding, 15, 15).
		To(
			g.Child().
				Size(g.Auto, 100).
				Layout(
					g.Style().
						SetColor(g.StyleColorText, colors["text"]).
						SetFontSize(16).
						To(
							g.Label("Paramètres Avancés"),
						),
					g.Dummy(0, 8),
					g.Row(
						g.Checkbox("Mise à jour automatique", &autoUpdateEnabled),
						g.Dummy(20, 0),
						g.Checkbox("Notifications", &showNotifications),
						g.Dummy(20, 0),
						g.Checkbox("Mode compact", &compactMode),
						g.Dummy(20, 0),
						g.Checkbox("Animations", &animationEnabled),
					),
				),
		)
}

func renderInstaller() g.Widget {
	candidates := makeAutoComplete()
	wi, _ := win.GetSize()
	w := float32(wi) - 96
	colors := getThemeColors(currentTheme)

	var currentDiscord *DiscordInstall
	if radioIdx != customChoiceIdx {
		currentDiscord = discords[radioIdx].(*DiscordInstall)
	}
	var isOpenAsar = currentDiscord != nil && currentDiscord.IsOpenAsar()

	if CanUpdateSelf() && !showedUpdatePrompt {
		showedUpdatePrompt = true
		g.OpenPopup("#update-prompt")
	}

	layout := g.Layout{
		// Header avec statistiques
		renderHeader(colors),
		g.Dummy(0, 15),

		// Panneau de contrôle avancé
		renderAdvancedPanel(colors),
		g.Dummy(0, 10),

		// Séparateur stylisé
		g.Style().
			SetColor(g.StyleColorSeparator, colors["accent"]).
			To(
				g.Separator(),
			),
		g.Dummy(0, 10),

		// Carte d'information de sécurité
		createInfoCard(
			"Sécurité",
			"**Github** est le seul endroit officiel pour obtenir Bashcord. Tout autre site prétendant être nous est malveillant.\n"+
				"Si vous avez téléchargé depuis une autre source, vous devriez tout supprimer/désinstaller immédiatement, effectuer une analyse anti-malware et changer votre mot de passe Discord.",
			colors,
			100,
		),

		g.Dummy(0, 15),

		// Titre de sélection
		g.Style().
			SetColor(g.StyleColorText, colors["accent"]).
			SetFontSize(24).
			To(
				g.Label("Sélectionnez une installation Discord à patcher"),
			),

		// Message d'erreur si aucune installation trouvée
		&CondWidget{len(discords) == 0, func() g.Widget {
			s := "Aucune installation Discord trouvee. Vous devez d'abord installer Discord."
			if runtime.GOOS == "linux" {
				s += " snap n'est pas pris en charge."
			}
			return createInfoCard("Aucune Installation", s, colors, 80)
		}, nil},

		// Liste des installations Discord
		g.Style().
			SetColor(g.StyleColorText, colors["text"]).
			SetFontSize(16).
			To(
				g.RangeBuilder("Discords", discords, func(i int, v any) g.Widget {
					d := v.(*DiscordInstall)
					//goland:noinspection GoDeprecation
					text := strings.Title(d.branch) + " - " + d.path
					if d.isPatched {
						text += " [PATCHE]"
					}
					return g.Style().
						SetColor(g.StyleColorCheckMark, colors["accent"]).
						SetStyleFloat(g.StyleVarFrameRounding, 6).
						To(
							g.RadioButton(text, radioIdx == i).
								OnChange(makeRadioOnChange(i)),
						)
				}),

				g.Style().
					SetColor(g.StyleColorCheckMark, colors["accent"]).
					SetStyleFloat(g.StyleVarFrameRounding, 6).
					To(
						g.RadioButton("Emplacement d'installation personnalise", radioIdx == customChoiceIdx).
							OnChange(makeRadioOnChange(customChoiceIdx)),
					),
			),

		g.Dummy(0, 10),

		// Champ de saisie personnalisé stylisé
		g.Style().
			SetStyle(g.StyleVarFramePadding, 16, 16).
			SetColor(g.StyleColorFrameBg, colors["secondary"]).
			SetColor(g.StyleColorFrameBgHovered, colors["accent"]).
			SetColor(g.StyleColorFrameBgActive, colors["accent"]).
			SetColor(g.StyleColorText, colors["text"]).
			SetFontSize(16).
			SetStyleFloat(g.StyleVarFrameRounding, 8).
			To(
				g.InputText(&customDir).Hint("Chemin personnalise vers Discord").
					Size(w - 16).
					Flags(g.InputTextFlagsCallbackCompletion).
					OnChange(onCustomInputChanged).
					Callback(
						func(data imgui.InputTextCallbackData) int32 {
							if len(candidates) == 0 {
								return 0
							}
							if autoCompleteIdx >= len(candidates) {
								autoCompleteIdx = 0
							}
							didAutoComplete = true
							start := len(customDir)
							if lastAutoComplete != "" {
								start -= len(lastAutoComplete)
								data.DeleteBytes(start, len(lastAutoComplete))
							} else if autoCompleteFile != "" {
								start -= len(autoCompleteFile)
								data.DeleteBytes(start, len(autoCompleteFile))
							}
							lastAutoComplete = candidates[autoCompleteIdx].(string)
							data.InsertBytes(start, []byte(lastAutoComplete))
							autoCompleteIdx++
							return 0
						},
					),
			),

		g.Dummy(0, 20),

		// Boutons d'action stylisés
		g.Style().SetFontSize(16).To(
			g.Row(
				createStyledButton("Installer", handlePatch, colors, (w-60)/4, 50),
				createStyledButton("Réparer", func() {
					if IsDevInstall {
						handlePatch()
					} else {
						err := InstallLatestBuilds()
						if err == nil {
							handlePatch()
						}
					}
				}, colors, (w-60)/4, 50),
				createStyledButton("Désinstaller", handleUnpatch, colors, (w-60)/4, 50),
				createStyledButton(Ternary(isOpenAsar, "Désinstaller OpenAsar", "Installer OpenAsar"), handleOpenAsar, colors, (w-60)/4, 50),
			),
		),

		InfoModal("#patched", "Patché avec succès", "Si Discord est encore ouvert, fermez-le complètement d'abord.\n"+
			"Ensuite, démarrez-le et vérifiez que Bashcord s'est installé avec succès en cherchant sa catégorie dans les Paramètres Discord"),
		InfoModal("#unpatched", "Dépatché avec succès", "Si Discord est encore ouvert, fermez-le complètement d'abord. Ensuite redémarrez-le, il devrait être revenu à l'état d'origine !"),
		InfoModal("#scuffed-install", "Attendez !", "Vous avez une installation Discord cassée.\n"+
			"Parfois Discord décide de s'installer au mauvais endroit pour une raison quelconque !\n"+
			"Vous devez corriger cela avant de patcher, sinon Bashcord ne fonctionnera probablement pas.\n\n"+
			"Utilisez le bouton ci-dessous pour y aller et supprimer tout dossier appelé Discord ou Squirrel.\n"+
			"Si le dossier est maintenant vide, n'hésitez pas à revenir en arrière et supprimer ce dossier aussi.\n"+
			"Ensuite voyez si Discord démarre toujours. Sinon, réinstallez-le"),
		RawInfoModal("#openasar-confirm", "OpenAsar", "OpenAsar est une alternative open-source de l'app.asar du bureau Discord.\n"+
			"Bashcord n'est en aucun cas affilié à OpenAsar.\n"+
			"Vous installez OpenAsar à vos propres risques. Si vous rencontrez des problèmes avec OpenAsar,\n"+
			"aucun support ne sera fourni, rejoignez plutôt le serveur OpenAsar !\n\n"+
			"Pour installer OpenAsar, appuyez sur Accepter et cliquez à nouveau sur 'Installer OpenAsar'.", true),
		InfoModal("#openasar-patched", "OpenAsar installé avec succès", "Si Discord est encore ouvert, fermez-le complètement d'abord. Ensuite redémarrez-le et vérifiez qu'OpenAsar s'est installé avec succès !"),
		InfoModal("#openasar-unpatched", "OpenAsar désinstallé avec succès", "Si Discord est encore ouvert, fermez-le complètement d'abord. Ensuite redémarrez-le et il devrait être revenu à l'état d'origine !"),
		InfoModal("#invalid-custom-location", "Emplacement invalide", "L'emplacement spécifié n'est pas une installation Discord valide.\nAssurez-vous de sélectionner le dossier de base.\n\nAstuce : Discord snap n'est pas pris en charge. utilisez flatpak ou .deb"),
		InfoModal("#modal"+strconv.Itoa(modalId), modalTitle, modalMessage),

		UpdateModal(),
	}

	return layout
}

func renderErrorCard(col color.Color, message string, height float32) g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, col).
		SetStyleFloat(g.StyleVarAlpha, 0.9).
		SetStyle(g.StyleVarWindowPadding, 10, 10).
		SetStyleFloat(g.StyleVarChildRounding, 5).
		To(
			g.Child().
				Size(g.Auto, height).
				Layout(
					g.Row(
						g.Style().SetColor(g.StyleColorText, color.Black).To(
							g.Markdown(&message),
						),
					),
				),
		)
}

func loop() {
	g.PushWindowPadding(48, 48)
	colors := getThemeColors(currentTheme)

	g.SingleWindow().
		RegisterKeyboardShortcuts(
			g.WindowShortcut{Key: g.KeyUp, Callback: func() {
				if radioIdx > 0 {
					radioIdx--
				}
			}},
			g.WindowShortcut{Key: g.KeyDown, Callback: func() {
				if radioIdx < customChoiceIdx {
					radioIdx++
				}
			}},
		).
		Layout(
			// Appliquer le thème sélectionné dynamiquement
			g.Style().
				SetColor(g.StyleColorWindowBg, color.RGBA{R: colors["primary"].R, G: colors["primary"].G, B: colors["primary"].B, A: 0xFF}).
				SetColor(g.StyleColorChildBg, color.RGBA{R: colors["secondary"].R, G: colors["secondary"].G, B: colors["secondary"].B, A: 0x90}).
				SetColor(g.StyleColorFrameBg, color.RGBA{R: colors["secondary"].R, G: colors["secondary"].G, B: colors["secondary"].B, A: 0x80}).
				SetColor(g.StyleColorFrameBgHovered, color.RGBA{R: colors["accent"].R, G: colors["accent"].G, B: colors["accent"].B, A: 0x60}).
				SetColor(g.StyleColorFrameBgActive, color.RGBA{R: colors["accent"].R, G: colors["accent"].G, B: colors["accent"].B, A: 0x80}).
				SetColor(g.StyleColorCheckMark, colors["accent"]).
				SetStyleFloat(g.StyleVarChildRounding, 12).
				SetStyleFloat(g.StyleVarFrameRounding, 8).
				To(
					g.Dummy(0, 20),
					g.Style().
						SetColor(g.StyleColorText, colors["text"]).
						SetFontSize(20).
						To(
							g.Row(
								g.Label(Ternary(IsDevInstall, "Installation de développement : ", "Bashcord sera téléchargé vers : ")+EquicordDirectory),
								g.Style().
									SetColor(g.StyleColorButton, colors["accent"]).
									SetColor(g.StyleColorButtonHovered, colors["secondary"]).
									SetColor(g.StyleColorButtonActive, colors["secondary"]).
									SetColor(g.StyleColorText, colors["primary"]).
									SetStyle(g.StyleVarFramePadding, 4, 4).
									To(
										g.Button("Ouvrir le repertoire").OnClick(func() {
											g.OpenURL("file://" + path.Dir(EquicordDirectory))
										}),
									),
							),
							&CondWidget{!IsDevInstall, func() g.Widget {
								return g.Style().
									SetColor(g.StyleColorText, color.RGBA{
										R: colors["text"].R,
										G: colors["text"].G,
										B: colors["text"].B,
										A: 180,
									}).
									To(
										g.Label("Pour personnaliser cet emplacement, définissez la variable d'environnement 'BASHCORD_USER_DATA_DIR' et redémarrez-moi").Wrapped(true),
									)
							}, nil},
							g.Dummy(0, 10),
							g.Style().
								SetColor(g.StyleColorText, color.RGBA{
									R: colors["text"].R,
									G: colors["text"].G,
									B: colors["text"].B,
									A: 160,
								}).
								To(
									g.Label("Version de Bashcord : "+buildinfo.InstallerTag+" ("+buildinfo.InstallerGitHash+")"+Ternary(IsSelfOutdated, " - OBSOLÈTE", "")),
									g.Label("Version locale de Bashcord : "+InstalledHash),
								),
							&CondWidget{
								GithubError == nil,
								func() g.Widget {
									if IsDevInstall {
										return g.Style().
											SetColor(g.StyleColorText, colors["warning"]).
											To(
												g.Label("Pas de mise à jour de Bashcord car en mode développement"),
											)
									}
									return g.Style().
										SetColor(g.StyleColorText, colors["success"]).
										To(
											g.Label("Dernière version de Bashcord : " + LatestHash),
										)
								}, func() g.Widget {
									return createInfoCard("Erreur GitHub", "Echec de recuperation des informations depuis GitHub : "+GithubError.Error(), colors, 60)
								},
							},
						),

					renderInstaller(),
					
					// Crédit en bas de page
					g.Dummy(0, 20),
					g.Style().
						SetColor(g.StyleColorText, color.RGBA{
							R: colors["text"].R,
							G: colors["text"].G,
							B: colors["text"].B,
							A: 140,
						}).
						SetFontSize(12).
						To(
							g.Row(
								g.Dummy(0, 0),
								g.Align(g.AlignCenter).To(
									g.Label("Prod by enfant divin. Contact discord : 9mf"),
								),
							),
						),
				),
		)

	g.PopStyle()
}
