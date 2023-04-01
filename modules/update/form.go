package update

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
	"os"
)

const (
	buttonConfirm = "button_confirm"
	buttonCancel  = "button_cancel"
)

const (
	labelQuestion = "label_question"
	labelVersion  = "label_version"
)

const (
	imageLogo = "image_logo"
)

type Updater struct {
	app     fyne.App
	windows *Windows
}

func NewUpdater(app fyne.App, windows *Windows) *Updater {
	return &Updater{app: app, windows: windows}
}

func (u *Updater) StartApp() {
	u.windows.mw.Window.Show()
	u.app.Run()
}

type Windows struct {
	mw *MainWindow
}

func NewWindows(mw *MainWindow) *Windows {
	return &Windows{mw: mw}
}

type MainWindow struct {
	Window  fyne.Window
	Images  map[string]*fyne.Resource
	Buttons map[string]*widget.Button
	Labels  map[string]*canvas.Text
	Box     *fyne.Container
	Config  struct {
		version    string
		execPath   string
		updatePath string
		owner      string
		repo       string
		PID        string
	}
}

func NewMainWindow(window fyne.Window) *MainWindow {
	return &MainWindow{
		Window:  window,
		Buttons: make(map[string]*widget.Button),
		Labels:  make(map[string]*canvas.Text),
		Images:  make(map[string]*fyne.Resource),
	}
}

func (w *MainWindow) ConfigurateWindow() {
	w.Window.Resize(fyne.NewSize(400, 450))
	w.Window.SetFixedSize(true)
	w.Window.SetMaster()
}

func (w *MainWindow) InitButtons(names []string, buttons ...*widget.Button) {
	for key, value := range buttons {
		w.Buttons[names[key]] = value
	}
}

func (w *MainWindow) InitLabels(names []string, labels ...*canvas.Text) {
	for key, value := range labels {
		w.Labels[names[key]] = value
	}
}

func (w *MainWindow) InitImages(names []string, images ...*fyne.Resource) {
	for key, value := range images {
		w.Images[names[key]] = value
	}
}

func (w *MainWindow) InitBox(box *fyne.Container) {
	w.Box = box
	w.Window.SetContent(box)
}

func (w *MainWindow) EnvError(err string) {
	errDialog := dialog.NewInformation("Ошибка", err, w.Window)
	errDialog.Show()
	errDialog.SetOnClosed(func() {
		io.WriteString(os.Stderr, err)
		os.Exit(1)
	})
	w.Window.SetOnClosed(func() {
		io.WriteString(os.Stderr, err)
		os.Exit(1)
	})
	w.Window.Resize(fyne.NewSize(600, 200))
	w.Window.RequestFocus()
	w.Window.CenterOnScreen()
	w.Window.ShowAndRun()
}
