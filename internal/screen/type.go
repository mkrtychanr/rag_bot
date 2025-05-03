package screen

type ScreenType int64

var (
	DefaultScreen     ScreenType = ScreenType(0)
	RequestScreen     ScreenType = ScreenType(1)
	AddDocumentScreen ScreenType = ScreenType(2)
	E2ESelectorScreen ScreenType = ScreenType(3)
)
