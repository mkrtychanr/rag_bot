package app

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/logger"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	"github.com/mkrtychanr/rag_bot/internal/screen/base"
	adddocumentintogroupscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/add_document_into_group_screen"
	adddocumentscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/add_document_screen"
	changedocumentsscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/change_documents_screen"
	deletedocumentfromgroupscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/delete_document_from_group_screen"
	deletedocumentscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/delete_document_screen"
	mainmenu "github.com/mkrtychanr/rag_bot/internal/screen/blocks/main_menu"
	requestscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/request_screen"
	postselectormenu "github.com/mkrtychanr/rag_bot/internal/screen/post_selector_menu"
	changedocument "github.com/mkrtychanr/rag_bot/internal/usecase/documents/change_document"
	deletedocument "github.com/mkrtychanr/rag_bot/internal/usecase/documents/delete_document"
	getdocuments "github.com/mkrtychanr/rag_bot/internal/usecase/documents/get_documents"
	groupaccess "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_access"
)

type superObject struct{}

func (s *superObject) MakeRequest(ctx context.Context, request string, tgID int64) (string, error) {
	return "super object make request", nil
}

func (s *superObject) AddDocument(ctx context.Context, name string, tgID int64, documentTelegramID string) error {
	logger.GetLogger().Info().Msgf("document %s added", name)

	return nil
}

func (a *app) newTree() screen.Screen {
	so := &superObject{}
	mainScreen := &mainmenu.DefaultMenuScreen{
		Base: base.Base{
			Title: "Добро пожаловать",
			Text:  "Добро пожаловать",
		},
	}

	requestScreen := requestscreen.NewRequestScreen(so, a.container.TgGateway, base.Base{
		Title:          "Сделать запрос",
		Text:           "Сделать запрос",
		HeadScreen:     mainScreen,
		PreviousScreen: mainScreen,
	})

	myDocumentsScreen := &mainmenu.DefaultMenuScreen{
		Base: base.Base{
			Title:          "Мои документы",
			Text:           "Мои документы",
			HeadScreen:     mainScreen,
			PreviousScreen: mainScreen,
		},
	}

	groupsScreen := &mainmenu.DefaultMenuScreen{
		Base: base.Base{
			Title:          "Группы",
			Text:           "Группы",
			HeadScreen:     mainScreen,
			PreviousScreen: mainScreen,
		},
	}

	mainScreen.NextScreens = []screen.Screen{requestScreen, myDocumentsScreen, groupsScreen}

	addDocumentScreen := adddocumentscreen.NewAddDocumentScreen(so, base.Base{
		Title:          "Добавить документ",
		Text:           "Добавить документ",
		HeadScreen:     mainScreen,
		PreviousScreen: myDocumentsScreen,
	})

	changeDocumentsScreen := changedocumentsscreen.NewChangeDocumentsScreen(
		getdocuments.NewUseCase(a.container.Repository),
		base.Base{
			Title:          "Изменить документ",
			Text:           "Изменить документ",
			HeadScreen:     mainScreen,
			PreviousScreen: myDocumentsScreen,
		},
	)

	postChangeDocumentsScreen := &postselectormenu.PostSelectorMenu{
		Base: base.Base{
			HeadScreen:     mainScreen,
			PreviousScreen: changeDocumentsScreen,
		},
		Text: "Как изменить",
	}

	changeDocumentUseCase := changedocument.NewUseCase(a.container.Repository)
	groupAccessUseCase := groupaccess.NewUseCase(a.container.Repository)

	addDocumentIntoGroupScreen := adddocumentintogroupscreen.NewAddDocumentIntoGroupScreen(
		changeDocumentUseCase,
		groupAccessUseCase,
		base.Base{
			Text:           "Добавить в группу",
			Title:          "Добавить в группу",
			HeadScreen:     mainScreen,
			PreviousScreen: postChangeDocumentsScreen,
		},
	)

	addDocumentIntoGroupScreen.NextScreens = []screen.Screen{addDocumentIntoGroupScreen}

	deleteDocumentFromGroupScreen := deletedocumentfromgroupscreen.NewDeleteDocumentFromGroupScreen(
		changeDocumentUseCase,
		groupAccessUseCase,
		base.Base{
			Text:           "Удалить из группы",
			Title:          "Удалить из группы",
			HeadScreen:     mainScreen,
			PreviousScreen: postChangeDocumentsScreen,
		},
	)

	deleteDocumentFromGroupScreen.NextScreens = []screen.Screen{deleteDocumentFromGroupScreen}

	postChangeDocumentsScreen.NextScreens = []screen.Screen{addDocumentIntoGroupScreen, deleteDocumentFromGroupScreen}

	changeDocumentsScreen.NextScreens = []screen.Screen{postChangeDocumentsScreen}

	deleteDocumentUseCase := deletedocument.NewUseCase(a.container.Repository)

	getDocumentsUseCase := getdocuments.NewUseCase(a.container.Repository)

	deleteDocumentScreen := deletedocumentscreen.NewDeleteDocumentScreen(
		deleteDocumentUseCase,
		getDocumentsUseCase,
		base.Base{
			Text:           "Удалить документ",
			Title:          "Удалить документ",
			HeadScreen:     mainScreen,
			PreviousScreen: myDocumentsScreen,
		},
	)

	deleteDocumentScreen.NextScreens = []screen.Screen{deleteDocumentScreen}

	myDocumentsScreen.NextScreens = []screen.Screen{addDocumentScreen, changeDocumentsScreen, deleteDocumentScreen}

	return mainScreen
}
