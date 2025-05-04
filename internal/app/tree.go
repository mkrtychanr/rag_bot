package app

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/logger"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	"github.com/mkrtychanr/rag_bot/internal/screen/base"
	adddocumentintogroupscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/add_document_into_group_screen"
	adddocumentscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/add_document_screen"
	changedocumentsscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/change_documents_screen"
	deletedocumentfromgroupscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/delete_document_from_group_screen"
	deletedocumentscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/delete_document_screen"
	groupaccessscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_access_screen"
	groupadduserscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_add_user_screen"
	groupchangenamescreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_change_name_screen"
	groupcreatescreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_create_screen"
	groupdocumentsscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_documents_screen"
	groupusersscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/group_users_screen"
	pickgroupscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/pick_group_screen"
	pickuserscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/pick_user_screen"
	mainmenu "github.com/mkrtychanr/rag_bot/internal/screen/blocks/main_menu"
	requestscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/request_screen"
	postselectormenu "github.com/mkrtychanr/rag_bot/internal/screen/post_selector_menu"
	changedocument "github.com/mkrtychanr/rag_bot/internal/usecase/documents/change_document"
	deletedocument "github.com/mkrtychanr/rag_bot/internal/usecase/documents/delete_document"
	getdocuments "github.com/mkrtychanr/rag_bot/internal/usecase/documents/get_documents"
	groupaccess "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_access"
	groupchange "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_change"
	groupcreate "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_create"
	groupdelete "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_delete"
	groupinfo "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_info"
	pickgroupchange "github.com/mkrtychanr/rag_bot/internal/usecase/groups/group_pick/pick_group_change"
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

	groupCreateUseCase := groupcreate.NewUseCase(a.container.Repository)

	createGroupScreen := groupcreatescreen.NewCreateGroupScreen(
		groupCreateUseCase,
		base.Base{
			HeadScreen:     mainScreen,
			PreviousScreen: myGroupsScreen,
		},
	)

	pickGroupForChangeScreen := pickgroupscreen.NewPickGroupScreen(
		pickgroupchange.NewUseCase(a.container.Repository),
		nil,
		base.Base{
			Title:          "Изменить группу",
			Text:           "Изменить группу",
			HeadScreen:     mainScreen,
			PreviousScreen: myGroupsScreen,
		},
	)

	postPickGroupForChangeScreen := &postselectormenu.PostSelectorMenu{
		Base: base.Base{
			HeadScreen:     mainScreen,
			PreviousScreen: pickGroupForChangeScreen,
		},
		Text: "Изменить группу",
	}

	changeGroupUseCase := groupchange.NewUseCase(a.container.Repository)

	changeGroupNameScreen := groupchangenamescreen.NewChangeGroupNameScreen(
		changeGroupUseCase,
		base.Base{
			HeadScreen:     mainScreen,
			PreviousScreen: postPickGroupForChangeScreen,
		},
	)

	addUserIntoGroupScreen := groupadduserscreen.NewAddUserIntoGroupScreen(
		changeGroupUseCase,
		base.Base{
			HeadScreen:     mainScreen,
			PreviousScreen: postPickGroupForChangeScreen,
		},
	)

	upRightsPolicyScreen := pickuserscreen.NewPickUserScreen(
		changeGroupUseCase.GetUsersWithReadOnlyRightPolicy,
		changeGroupUseCase.SetReadWriteRightsForUserInGroup,
		base.Base{
			Title:          "Повысить права",
			Text:           "Повысить права",
			HeadScreen:     mainScreen,
			PreviousScreen: postPickGroupForChangeScreen,
		},
	)

	upRightsPolicyScreen.NextScreens = []screen.Screen{upRightsPolicyScreen}

	downRigthsPolicyScreen := pickuserscreen.NewPickUserScreen(
		changeGroupUseCase.GetUsersWithReadWriteRightPolicy,
		changeGroupUseCase.SetReadOnlyRightsForUserInGroup,
		base.Base{
			Title:          "Понизить права",
			Text:           "Понизить права",
			HeadScreen:     mainScreen,
			PreviousScreen: postPickGroupForChangeScreen,
		},
	)

	downRigthsPolicyScreen.NextScreens = []screen.Screen{downRigthsPolicyScreen}

	deleteUserFromGroupScreen := pickuserscreen.NewPickUserScreen(
		changeGroupUseCase.GetGroupUsersToDelete,
		changeGroupUseCase.DeleteUserFromGroup,
		base.Base{
			Title:          "Удалить из группы",
			Text:           "Удалить из группы",
			HeadScreen:     mainScreen,
			PreviousScreen: postPickGroupForChangeScreen,
		},
	)

	deleteUserFromGroupScreen.NextScreens = []screen.Screen{deleteUserFromGroupScreen}

	postPickGroupForChangeScreen.NextScreens = []screen.Screen{changeGroupNameScreen, addUserIntoGroupScreen, upRightsPolicyScreen, downRigthsPolicyScreen, deleteUserFromGroupScreen}

	pickGroupForChangeScreen.NextScreens = []screen.Screen{postPickGroupForChangeScreen}

	deleteGroupUseCase := groupdelete.NewUseCase(a.container.Repository)

	f := func(ctx context.Context, payload map[string]any) error {
		groupID, ok := payload["group_id"].(int64)
		if !ok {
			return screen.ErrWrongType
		}

		if err := deleteGroupUseCase.DeleteGroup(ctx, groupID); err != nil {
			return fmt.Errorf("failed to delete group: %w", err)
		}

		return nil
	}

	pickGroupForDelete := pickgroupscreen.NewPickGroupScreen(
		pickgroupchange.NewUseCase(a.container.Repository),
		f,
		base.Base{
			Title:          "Удалить группу",
			Text:           "Удалить группу",
			HeadScreen:     mainScreen,
			PreviousScreen: myGroupsScreen,
		},
	)

	pickGroupForDelete.NextScreens = []screen.Screen{pickGroupForDelete}

	myGroupsScreen.NextScreens = []screen.Screen{pickGroupForInfoScreen, createGroupScreen, pickGroupForChangeScreen, pickGroupForDelete}

	groupsScreen.NextScreens = []screen.Screen{accessGroupsScreen, myGroupsScreen}

	return mainScreen
}
