package windows

import (
	"syscall"
	"unsafe"
)

// 创建菜单
func createPopupMenu() uintptr {
	ret, _, _ := CreatePopupMenu.Call()
	return ret
}

type Menu struct {
	Handle    uintptr
	Callbacks map[uint32]func()
}

func NewMenu() *Menu {
	return &Menu{
		Handle:    createPopupMenu(),
		Callbacks: make(map[uint32]func()),
	}
}

// 添加菜单项
func (m *Menu) AddItem(id uint32, label string, callback func()) {
	text, _ := syscall.UTF16PtrFromString(label)
	AppendMenu.Call(m.Handle, uintptr(MF_STRING), uintptr(id), uintptr(unsafe.Pointer(text)))
	m.Callbacks[id] = callback
}

// 添加子菜单
func (m *Menu) AddSubMenu(label string, subMenu *Menu) {
	text, _ := syscall.UTF16PtrFromString(label)
	AppendMenu.Call(m.Handle, uintptr(MF_POPUP), subMenu.Handle, uintptr(unsafe.Pointer(text)))
	for id, callback := range subMenu.Callbacks {
		m.Callbacks[id] = callback
	}
}

// 添加复选菜单项
func (m *Menu) AddCheckMenu(id uint32, label string, check bool, callback func()) {
	text, _ := syscall.UTF16PtrFromString(label)
	if check {
		AppendMenu.Call(m.Handle, uintptr(MF_CHECKED), uintptr(id), uintptr(unsafe.Pointer(text)))
	} else {
		AppendMenu.Call(m.Handle, uintptr(MF_UNCHECKED), uintptr(id), uintptr(unsafe.Pointer(text)))
	}
	m.Callbacks[id] = callback
}

// 分隔线
func (m *Menu) AddSeparator() {
	AppendMenu.Call(m.Handle, uintptr(MF_SEPARATOR))
}

// 检查并切换复选菜单状态，返回新的状态
func (m *Menu) ToggleCheck(id int) bool {
	ret, _, _ := CheckMenuItem.Call(
		m.Handle,
		uintptr(id),
		uintptr(MF_BYCOMMAND),
	)

	newState := ret != MF_CHECKED
	v := MF_UNCHECKED
	if newState {
		v = MF_CHECKED
	}

	CheckMenuItem.Call(
		m.Handle,
		uintptr(id),
		uintptr(v),
	)

	return newState
}
