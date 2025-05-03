// Package main is a one-file package. Here's only main function.
package main

import (
	// "github.com/mkrtychanr/rag_bot/cmd/commands"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mkrtychanr/rag_bot/cmd/commands"
	"github.com/mkrtychanr/rag_bot/internal/logger"
	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	requestscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/request_screen"
	_ "go.uber.org/automaxprocs"
)

type superObject struct{}

func (s *superObject) MakeRequest(ctx context.Context, request string, tgID int64) (string, error) {
	return "super object make request", nil
}

func (s *superObject) Send(ctx context.Context, userID int64, message string) error {
	logger.GetLogger().Info().Msgf("send %s", message)

	return nil
}

func printScreen(sc model.Screen) {
	fmt.Println(sc.Text)
	for _, buttons := range sc.Buttons {
		for _, button := range buttons {
			fmt.Printf("|%-20s", button.Text)
		}

		if len(buttons) != 0 {
			fmt.Printf("|")
		}

		fmt.Println()
	}
	fmt.Println("*******************************************")
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func test(curr screen.Screen) {
	sc, err := curr.Render()
	fatal(err)

	printScreen(sc)

	curr, err = curr.Next(context.Background(), model.MenuOption{
		Option: 0,
	})
	fatal(err)

	sc, err = curr.Render()
	fatal(err)

	printScreen(sc)

	data, err := json.Marshal(requestscreen.PerformModel{
		UserID:  123,
		Request: "aboba",
	})
	fatal(err)

	curr, err = curr.Perform(context.Background(), data)
	fatal(err)

	sc, err = curr.Render()
	fatal(err)

	printScreen(sc)

	curr, err = curr.Next(context.Background(), model.MenuOption{
		Option: -1,
	})
	fatal(err)

	sc, err = curr.Render()
	fatal(err)

	printScreen(sc)
}

func main() {
	commands.Execute()

	// so := &superObject{}

	// mainScreen := &mainmenu.DefaultMenuScreen{}

	// requestScreen := requestscreen.NewRequestScreen(so, so, base.Base{
	// 	Title:          "Сделать запрос",
	// 	Text:           "Сделать запрос",
	// 	HeadScreen:     mainScreen,
	// 	PreviousScreen: mainScreen,
	// })

	// myDocumentsScreen := &mainmenu.DefaultMenuScreen{
	// 	Base: base.Base{
	// 		Title:          "Мои документы",
	// 		Text:           "Мои документы",
	// 		HeadScreen:     mainScreen,
	// 		PreviousScreen: mainScreen,
	// 	},
	// }

	// groupsScreen := &mainmenu.DefaultMenuScreen{
	// 	Base: base.Base{
	// 		Title:          "Группы",
	// 		Text:           "Группы",
	// 		HeadScreen:     mainScreen,
	// 		PreviousScreen: mainScreen,
	// 	},
	// }

	// mainScreen.NextScreens = []screen.Screen{requestScreen, myDocumentsScreen, groupsScreen}

	// test(mainScreen)
}
