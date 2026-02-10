//go:build cli

/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"vencord/buildinfo"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var discords []any
var interactive = false

func showBanner() {
	color.HiRed(`                                                                                                    
                                           :::::::::::::                                            
                                             ::::::::::::    ::                                     
                                                ::::::::::-  ::::                                   
                                                  ::::::::::  :::::                                 
                                                    ----:--:: -::::::                               
                                              ::::----::::::: ---::::                               
                                           ::::----::::::-====------ +*                             
                                             :---+-::::-===+++=:--=--*+---                          
                                              :---+=::::=======:--::=#*-====                        
                                            ------+*=::.:------::::=*%#--=-    ==++                 
                                     :-    -==----+%##:.:::::::::--=#%*::-   -==++++*+              
                                    ------   ---::-===----:::::=--=+****+++:--==++++=   =+          
                                   ====----::  ::-----====+=::==---========-:---=  --===++          
                                  --=----:::..-:::::::-===-==---==--=++===+--::----=====++          
                               ------    :::...:.::-+*##===--=:=-:=***%%=-=-::::---======           
                                ------:::.........:++++%#----+---:=+**%#----*%%=                    
                                ------:::.....:....-=+**=:--:---=--.:::=--==**##%%#+---=            
                                        :::+##*-:...::.::::----=--====---+=**::::------             
                                  :::::-=+====+-::::.:::::--+=-+=--======+####*=-  :--              
                                    :-:::::. :-+*#+::::.:=:-=---+::=:+##*=---=+**=--                
                                     :::  :::+*+=====+-....-:::::=::::=+**+-:------                 
                                        :----::::-+++==:---:::::::=**#=::--==--=-                   
                                          -----:---::.:++**+::::::-=+#+------=+                     
                                            ...::::::::=+=-::::::::::-+++++                         
                                           :::....-::::-::::....::=---==+**                         
                                           ::::=*++==::::::+***=--++-::.+*                          
              --                      ::=   =-:::--===-::::--::-------..:-                          
             -=+                   ...::=++-:::::::--::::-=====-:-:::...:                           
            :--=                 --=++=---=*+=:::-:::::::::::----:.......                           
            :--=+*        :-:: =---::---**##*++=::--:::::.:-:.::--=+                                
           .::::-==+==   +=:.:-:::.....:==+****+=-:::==-::::::-:=*                                  
            ::::::----::::--::::...:::::-=++++++=--::===-----:==-                                   
             :::::::::..:....:..--::---====+====-:::--===+=====*                                    
          ..:    ..:::-::=:.:::-::..::-=======---:::--=++====+                                      
          :::-      ..:::--==--::....:::---:::::::::-=+=---+*                                       
          ::::=*      :---=++==::-:.:..:-:::--::::::-=----+*                                        
          .:::.-=+++++=-:---------::..:==--  ....:::::::-==                                         
          :::::..:---==-:+-+:::===.::::        .........:+                                          
          .::::.....::--::-===    ...              .:-=                                             
               .........::-=                                                                        
                     .::::*                                                                         `)
	fmt.Println()
}

func isValidBranch(branch string) bool {
	switch branch {
	case "", "stable", "ptb", "canary", "auto":
		return true
	default:
		return false
	}
}

func die(msg string) {
	Log.Error(msg)
	exitFailure()
}

func main() {
	// Agrandir la console sur Windows pour mieux afficher l'ASCII art
	if runtime.GOOS == "windows" {
		ResizeConsoleWindow()
	}
	
	InitGithubDownloader()
	discords = FindDiscords()

	// Used by log.go init func
	flag.Bool("debug", false, "Activer les infos de debug (pour les masochistes)")

	var helpFlag = flag.Bool("help", false, "Afficher les instructions d'usage (si tu sais pas lire)")
	var versionFlag = flag.Bool("version", false, "Voir la version du programme (passionnant)")
	var updateSelfFlag = flag.Bool("update-self", false, "Me mettre à jour (j'en ai besoin)")
	var installFlag = flag.Bool("install", false, "Installer BASHCORD (enfin !)")
	var updateFlag = flag.Bool("repair", false, "Réparer BASHCORD (encore cassé ?)")
	var uninstallFlag = flag.Bool("uninstall", false, "Désinstaller BASHCORD (tu abandonnes déjà ?)")
	var installOpenAsarFlag = flag.Bool("install-openasar", false, "Installer OpenAsar (pour les vrais)")
	var uninstallOpenAsarFlag = flag.Bool("uninstall-openasar", false, "Désinstaller OpenAsar (retour aux basiques)")
	var locationFlag = flag.String("location", "", "L'emplacement de Discord à modifier")
	var branchFlag = flag.String("branch", "", "La branche Discord à modifier [auto|stable|ptb|canary]")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *versionFlag {
		fmt.Println("Equilotl Cli", buildinfo.InstallerTag, "("+buildinfo.InstallerGitHash+")")
		fmt.Println("Copyright (C) 2025 Vendicated et les contributeurs Vencord")
		fmt.Println("Licence GPLv3+ : GNU GPL version 3 ou plus récente <https://gnu.org/licenses/gpl.html>.")
		return
	}

	if *updateSelfFlag {
		if !<-SelfUpdateCheckDoneChan {
			die("Impossible de me mettre à jour car la vérification des mises à jour a échoué (bravo)")
		}
		if err := UpdateSelf(); err != nil {
			Log.Error("Échec de la mise à jour automatique :", err)
			exitFailure()
		}
		exitSuccess()
	}

	if *locationFlag != "" && *branchFlag != "" {
		die("Les flags 'location' et 'branch' sont mutuellement exclusifs (choisis-en un, génie).")
	}

	if !isValidBranch(*branchFlag) {
		die("Le flag 'branch' doit être l'un des suivants : [auto|stable|ptb|canary] (pas si compliqué)")
	}

	if *installFlag || *updateFlag {
		if !<-GithubDoneChan {
			die("Pas d'" + Ternary(*installFlag, "installation", "mise à jour") + " car la récupération des données de release a échoué (GitHub nous boude)")
		}
	}

	install, uninstall, update, installOpenAsar, uninstallOpenAsar := *installFlag, *uninstallFlag, *updateFlag, *installOpenAsarFlag, *uninstallOpenAsarFlag
	switches := []*bool{&install, &update, &uninstall, &installOpenAsar, &uninstallOpenAsar}
	if !SliceContainsFunc(switches, func(b *bool) bool { return *b }) {
		interactive = true

		// Afficher le banner ASCII seulement en mode interactif
		showBanner()

		go func() {
			<-SelfUpdateCheckDoneChan
			if IsSelfOutdated {
				Log.Warn("Ton installateur est obsolète (comme ton PC probablement).")
				Log.Warn("Pour mettre à jour, sélectionne l'option 'Mettre à jour Bashcord_CLI' ou lance avec --update-self")
			}
		}()

		choices := []string{
			"Installer FILS DE PUTE",
			"Réparer SALE DOG ",
			"Désinstaller FAIS PAS STP ",
			"Installer OpenAsar (pour les connaisseurs)",
			"Désinstaller OpenAsar (retour en arrière)",
			"Voir le menu d'aide (RTFM)",
			"Mettre à jour Bashcord_CLI (fais-le !)",
			"Quitter (fuyaaaaard !)",
		}
		_, choice, err := (&promptui.Select{
			Label: "Que veux-tu faire ? (Appuie sur Entrée sois pas con)",
			Items: choices,
			HideHelp: true,
		}).Run()
		handlePromptError(err)

		switch choice {
		case "Voir le menu d'aide (RTFM)":
			flag.Usage()
			return
		case "Quitter (fuyaaaaard !)":
			return
		case "Mettre à jour Equilotl (fais-le !)":
			if err := UpdateSelf(); err != nil {
				Log.Error("Échec de la mise à jour automatique :", err)
				exitFailure()
			}
			exitSuccess()
		}

		*switches[SliceIndex(choices, choice)] = true
	}

	var err error
	var errSilent error
	if install {
		errSilent = PromptDiscord("patcher", *locationFlag, *branchFlag).patch()
	} else if uninstall {
		errSilent = PromptDiscord("dépatcher", *locationFlag, *branchFlag).unpatch()
	} else if update {
		Log.Info("Téléchargement des derniers fichiers Bashcord... (patience, petit scarabée)")
		err := installLatestBuilds()
		Log.Info("Terminé ! (miracle)")
		if err == nil {
			errSilent = PromptDiscord("réparer", *locationFlag, *branchFlag).patch()
		}
	} else if installOpenAsar {
		discord := PromptDiscord("patcher", *locationFlag, *branchFlag)
		if !discord.IsOpenAsar() {
			err = discord.InstallOpenAsar()
		} else {
			die("OpenAsar déjà installé (tu dors ou quoi ?)")
		}
	} else if uninstallOpenAsar {
		discord := PromptDiscord("patcher", *locationFlag, *branchFlag)
		if discord.IsOpenAsar() {
			err = discord.UninstallOpenAsar()
		} else {
			die("OpenAsar pas installé (logique, non ?)")
		}
	}

	if err != nil {
		Log.Error(err)
		exitFailure()
	}
	if errSilent != nil {
		exitFailure()
	}

	exitSuccess()
}

func exit(status int) {
	if runtime.GOOS == "windows" && IsDoubleClickRun() && interactive {
		fmt.Print("Appuie sur Entrée pour quitter (si tu y arrives)")
		var b byte
		_, _ = fmt.Scanf("%v", &b)
	}
	os.Exit(status)
}

func exitSuccess() {
	color.HiGreen("✔ Succès ! (incroyable)")
	exit(0)
}

func exitFailure() {
	color.HiRed("❌ Échec ! (comme d'habitude)")
	exit(1)
}

func handlePromptError(err error) {
	if errors.Is(err, promptui.ErrInterrupt) {
		exit(0)
	}

	Log.FatalIfErr(err)
}

func PromptDiscord(action, dir, branch string) *DiscordInstall {
	if branch == "auto" {
		for _, b := range []string{"stable", "canary", "ptb"} {
			for _, discord := range discords {
				install := discord.(*DiscordInstall)
				if install.branch == b {
					return install
				}
			}
		}
		die("Aucune installation Discord trouvée. Essaie de la spécifier manuellement avec le flag --dir. Indice : snap n'est pas supporté (évidemment)")
	}

	if branch != "" {
		for _, discord := range discords {
			install := discord.(*DiscordInstall)
			if install.branch == branch {
				return install
			}
		}
		die("Discord " + branch + " introuvable (tu es sûr qu'il existe ?)")
	}

	if dir != "" {
		if discord := ParseDiscord(dir, branch); discord != nil {
			return discord
		} else {
			die(dir + " n'est pas une installation Discord valide. Indice : snap n'est pas supporté (on t'avait prévenu)")
		}
	}

	items := SliceMap(discords, func(d any) string {
		install := d.(*DiscordInstall)
		//goland:noinspection GoDeprecation
		return fmt.Sprintf("%s - %s%s", strings.Title(install.branch), install.path, Ternary(install.isPatched, " [PATCHÉ]", ""))
	})
	items = append(items, "Emplacement personnalisé (pour les rebelles)")

	_, choice, err := (&promptui.Select{
		Label: "Sélectionne l'installation Discord à " + action + " (Appuie sur Entrée pour confirmer, courage !)",
		Items: items,
		HideHelp: true,
	}).Run()
	handlePromptError(err)

	if choice != "Emplacement personnalisé (pour les rebelles)" {
		return discords[SliceIndex(items, choice)].(*DiscordInstall)
	}

	for {
		custom, err := (&promptui.Prompt{
			Label: "Emplacement Discord personnalisé (j'espère que tu sais ce que tu fais)",
		}).Run()
		handlePromptError(err)

		if di := ParseDiscord(custom, ""); di != nil {
			return di
		}

		Log.Error("Installation Discord invalide ! (surprise)")
	}
}

func InstallLatestBuilds() error {
	return installLatestBuilds()
}

func HandleScuffedInstall() {
	fmt.Println("Attends un peu !")
	fmt.Println("Tu as une installation Discord cassée (bravo l'artiste).")
	fmt.Println("Veuillez réinstaller Discord avant de continuer !")
	fmt.Println("Sinon, Equicord ne fonctionnera probablement pas (logique).")
}
