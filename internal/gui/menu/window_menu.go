package menu

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/gui/settings"
)

func MakeMenu(guiApp fyne.App, ollamaClient *ollamaApi.Client, w fyne.Window) *fyne.MainMenu {
	openSettings := func() {
		settings.ShowSettings(guiApp, ollamaClient)
	}
	showAbout := func() {
		settings.ShowAbout(guiApp)
	}
	aboutItem := fyne.NewMenuItem("About", showAbout)
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	performFind := func() { fmt.Println("Menu Find") }
	findItem := fyne.NewMenuItem("Find", performFind)
	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
		performFind()
	})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://ctrlplusrevise.com/docs/tutorials/")
			_ = guiApp.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://discord.gg/TYBtGUdVBU")
			_ = guiApp.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://www.patreon.com/SalmonsStudios")
			_ = guiApp.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File")
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	file.Items = append(file.Items, aboutItem)
	main := fyne.NewMainMenu(
		file,
		helpMenu,
	)
	return main
}
