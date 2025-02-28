package command

type App struct {
	configPath string
}

func NewApp() *App {
	return &App{}
}
