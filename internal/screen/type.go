package screen

type ScreenType int64

var (
	DefaultScreen                 ScreenType = ScreenType(0)
	RequestScreen                 ScreenType = ScreenType(1)
	AddDocumentScreen             ScreenType = ScreenType(2)
	E2ESelectorScreen             ScreenType = ScreenType(3)
	PostE2ESelectorScreen         ScreenType = ScreenType(4)
	AddDocumentIntoGroupScreen    ScreenType = ScreenType(5)
	ChangeDocumentsScreen         ScreenType = ScreenType(6)
	DeleteDocumentFromGroupScreen ScreenType = ScreenType(7)
	DeleteDocumentScreen          ScreenType = ScreenType(8)
	GroupAccessScreen             ScreenType = ScreenType(9)
	PickGroupScreen               ScreenType = ScreenType(10)
	GroupDocumentsScreen          ScreenType = ScreenType(11)
	GroupUsersScreen              ScreenType = ScreenType(12)
	CreateGroupScreen             ScreenType = ScreenType(13)
	ChangeGroupNameScreen         ScreenType = ScreenType(14)
	AddUserIntoGroupScreen        ScreenType = ScreenType(15)
	PickUserScreen                ScreenType = ScreenType(16)
)
