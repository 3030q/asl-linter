package main

import (
	"github.com/pterm/pterm"
	"strconv"
	"time"
)

func introScreen() {
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("COOL", pterm.NewStyle(pterm.FgLightGreen)),
		pterm.NewLettersFromStringWithStyle("MONKES", pterm.NewStyle(pterm.FgLightMagenta))).
		Render()

	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgDarkGray)).WithMargin(10).Println(
		"AWS-LINTER - static linter for amazon states language with usage grpc & proto3")

	pterm.Info.Println("Это консольная утилита для статичекой проверки плейбуков Amazon States Language" +
		"\nНа данный момент утилита работает лишь с inline автоматами" +
		"\nНе работает на плэйбуках с содержанием параллельностьи, map и choice" +
		"\n" +
		"\nДля запуска утилиты необходимо указать имя файла без их расширения для ASL структуры и proto3 файла" +
		"\nОни должны находиться в одной директории, что и утилита" +
		"\n" +
		"\nЗа большей информаций можно обращаться в tg @NikitaRybin888 :)" +
		"\n" +
		"\nActual date " + pterm.Green(time.Now().Format("02 Jan 2006 - 15:04:05 MST")))
	pterm.Println()
	introSpinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Doing something important...")
	time.Sleep(time.Second)
	for i := 3; i > 0; i-- {
		if i > 1 {
			introSpinner.UpdateText("Doing something important... " + strconv.Itoa(i) + " seconds...")
		} else {
			introSpinner.UpdateText("Doing something important... " + strconv.Itoa(i) + " second...")
		}
		time.Sleep(time.Second)
	}
	introSpinner.Stop()
}
