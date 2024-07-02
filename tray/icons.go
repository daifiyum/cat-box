package tray

import "os"

var (
	AppIcon   []byte
	CloseIcon []byte
	HomeIcon  []byte
	ProxyIcon []byte
	SubIcon   []byte
)

func InitIcons() {
	AppIcon, _ = os.ReadFile("resources/icons/box.ico")
	CloseIcon, _ = os.ReadFile("resources/icons/close.ico")
	HomeIcon, _ = os.ReadFile("resources/icons/home.ico")
	ProxyIcon, _ = os.ReadFile("resources/icons/proxy.ico")
	SubIcon, _ = os.ReadFile("resources/icons/sub.ico")
}
