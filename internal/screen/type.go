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
)
