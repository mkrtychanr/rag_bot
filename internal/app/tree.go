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
	groupaccessscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_access_screen"
	groupdocumentsscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_documents_screen"
	groupusersscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_users_screen"
	pickgroupscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/pick_group_screen"
	mainmenu "github.com/mkrtychanr/rag_bot/internal/screen/blocks/main_menu"
	requestscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/request_screen"
	postselectormenu "github.com/mkrtychanr/rag_bot/internal/screen/post_selector_menu"
	changedocument "github.com/mkrtychanr/rag_bot/internal/usecase/documents/change_document"
	deletedocument "github.com/mkrtychanr/rag_bot/internal/usecase/documents/delete_document"
	getdocuments "github.com/mkrtychanr/rag_bot/internal/usecase/documents/get_documents"
	groupaccess "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_access"
	groupinfo "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_info"
	pickinfogroup "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_pick/pick_info_group"
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

	accessGroupsScreen := groupaccessscreen.NewGroupAccessScreen(
		groupAccessUseCase,
		base.Base{
			Title:          "В каких я группах",
			Text:           "В каких я группах",
			HeadScreen:     mainScreen,
			PreviousScreen: groupsScreen,
		},
	)

	myGroupsScreen := &mainmenu.DefaultMenuScreen{
		Base: base.Base{
			Title:          "Мои группы",
			Text:           "Мои группы",
			HeadScreen:     mainScreen,
			PreviousScreen: groupsScreen,
		},
	}

	pickGroupForInfoScreen := pickgroupscreen.NewPickGroupScreen(
		pickinfogroup.NewUseCase(a.container.Repository),
		nil,
		base.Base{
			Title:          "Информация по группе",
			Text:           "Информация по группе",
			HeadScreen:     mainScreen,
			PreviousScreen: myGroupsScreen,
		},
	)

	postPickGroupForInfoScreen := &postselectormenu.PostSelectorMenu{
		Base: base.Base{
			HeadScreen:     mainScreen,
			PreviousScreen: pickGroupForInfoScreen,
		},
		Text: "Информация о группе",
	}

	groupInfoUseCase := groupinfo.NewUseCase(a.container.Repository)

	groupDocumentsScreen := groupdocumentsscreen.NewGroupDocumentsScreen(
		groupInfoUseCase,
		base.Base{
			Title:          "Документы группы",
			Text:           "Документы группы",
			HeadScreen:     mainScreen,
			PreviousScreen: postPickGroupForInfoScreen,
		},
	)

	groupUsersScreen := groupusersscreen.NewGroupUsersScreen(
		groupInfoUseCase,
		base.Base{
			Title:          "Пользователи группы",
			Text:           "Пользователи группы",
			HeadScreen:     mainScreen,
			PreviousScreen: postPickGroupForInfoScreen,
		},
	)

	postPickGroupForInfoScreen.NextScreens = []screen.Screen{groupDocumentsScreen, groupUsersScreen}

	pickGroupForInfoScreen.NextScreens = []screen.Screen{postPickGroupForInfoScreen}

	myGroupsScreen.NextScreens = []screen.Screen{pickGroupForInfoScreen}

	groupsScreen.NextScreens = []screen.Screen{accessGroupsScreen, myGroupsScreen}

	return mainScreen
}
