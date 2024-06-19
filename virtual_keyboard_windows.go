package virtual_keyboard

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unicode"
	"unsafe"
)

var user32DLL = windows.NewLazyDLL("user32.dll")
var procPostMsg = user32DLL.NewProc("PostMessageW")
var findWindowW = user32DLL.NewProc("FindWindowW")

// Press key(s)
func (k *KeyBonding) Press() error {
	if k.hasALT {
		k.downKey(VK_ALT)
	}
	if k.hasALTGR {
		k.downKey(VK_ALT)
		k.downKey(VK_CTRL)
	}
	if k.hasSHIFT {
		k.downKey(VK_SHIFT)
	}
	if k.hasCTRL {
		k.downKey(VK_CTRL)
	}
	if k.hasRSHIFT {
		k.downKey(VK_RSHIFT)
	}
	if k.hasRCTRL {
		k.downKey(VK_RCTRL)
	}

	for _, key := range k.keys {
		k.downKey(key)
	}
	return nil
}

// Release key(s)
func (k *KeyBonding) Release() error {
	if k.hasALT {
		k.upKey(VK_ALT)
	}
	if k.hasALTGR {
		k.upKey(VK_ALT)
		k.upKey(VK_CTRL)
	}
	if k.hasSHIFT {
		k.upKey(VK_SHIFT)
	}
	if k.hasCTRL {
		k.upKey(VK_CTRL)
	}
	if k.hasRSHIFT {
		k.upKey(VK_RSHIFT)
	}
	if k.hasRCTRL {
		k.upKey(VK_RCTRL)
	}
	for _, key := range k.keys {
		k.upKey(key)
	}
	return nil
}

// Launch key bounding
func (k *KeyBonding) Launch() error {
	err := k.Press()
	if err != nil {
		return err
	}
	err = k.Release()
	return err
}

func (k *KeyBonding) upKey(key interface{}) error {
	keyCode, err := getKeyCode(key)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, _, _ = procPostMsg.Call(k.hwnd, 0x0101, keyCode, 0)

	return nil
}

func (k *KeyBonding) downKey(key interface{}) error {
	keyCode, err := getKeyCode(key)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, _, _ = procPostMsg.Call(k.hwnd, 0x0100, keyCode, 0)

	return nil
}

func (k *KeyBonding) AddHWND(hwnd uintptr) {
	k.hwnd = hwnd
}

func (k *KeyBonding) AddHWNDByTitle(title string) {
	k.hwnd = findWindowByTitle(title)
}

func findWindowByTitle(title string) uintptr {
	hwnd, _, _ := findWindowW.Call(0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))))

	return hwnd
}

func getKeyCode(input interface{}) (uintptr, error) {
	switch v := input.(type) {
	case rune:
		if (v >= '0' && v <= '9') || (v >= 'A' && v <= 'Z') || (v >= 'a' && v <= 'z') {
			return uintptr(unicode.ToUpper(v)), nil
		}
		return 0, fmt.Errorf("Invalid character: %c", v)
	case int:
		return uintptr(v), nil
	default:
		return 0, fmt.Errorf("Invalid type: only rune and int are allowed")
	}
}

const (
	VK_LBUTTON    = 0x01
	VK_RBUTTON    = 0x02
	VK_CANCEL     = 0x03
	VK_MBUTTON    = 0x04
	VK_XBUTTON1   = 0x05
	VK_XBUTTON2   = 0x06
	VK_BACK       = 0x08
	VK_TAB        = 0x09
	VK_CLEAR      = 0x0C
	VK_RETURN     = 0x0D
	VK_SHIFT      = 0x10
	VK_CONTROL    = 0x11
	VK_CTRL       = 0x11
	VK_MENU       = 0x12
	VK_ALT        = 0x12
	VK_PAUSE      = 0x13
	VK_CAPITAL    = 0x14
	VK_KANA       = 0x15
	VK_HANGUEL    = 0x15
	VK_HANGUL     = 0x15
	VK_JUNJA      = 0x17
	VK_FINAL      = 0x18
	VK_HANJA      = 0x19
	VK_KANJI      = 0x19
	VK_ESCAPE     = 0x1B
	VK_CONVERT    = 0x1C
	VK_NONCONVERT = 0x1D
	VK_ACCEPT     = 0x1E
	VK_MODECHANGE = 0x1F
	VK_SPACE      = 0x20
	VK_PRIOR      = 0x21
	VK_NEXT       = 0x22
	VK_END        = 0x23
	VK_HOME       = 0x24
	VK_LEFT       = 0x25
	VK_UP         = 0x26
	VK_RIGHT      = 0x27
	VK_DOWN       = 0x28
	VK_SELECT     = 0x29
	VK_PRINT      = 0x2A
	VK_EXECUTE    = 0x2B
	VK_SNAPSHOT   = 0x2C
	VK_INSERT     = 0x2D
	VK_DELETE     = 0x2E
	VK_HELP       = 0x2F
	// VK_0 - VK_9 are the same as ASCII '0' - '9' (0x30 - 0x39)
	// VK_A - VK_Z are the same as ASCII 'A' - 'Z' (0x41 - 0x5A)
	VK_LWIN      = 0x5B
	VK_RWIN      = 0x5C
	VK_APPS      = 0x5D
	VK_SLEEP     = 0x5F
	VK_NUMPAD0   = 0x60
	VK_NUMPAD1   = 0x61
	VK_NUMPAD2   = 0x62
	VK_NUMPAD3   = 0x63
	VK_NUMPAD4   = 0x64
	VK_NUMPAD5   = 0x65
	VK_NUMPAD6   = 0x66
	VK_NUMPAD7   = 0x67
	VK_NUMPAD8   = 0x68
	VK_NUMPAD9   = 0x69
	VK_MULTIPLY  = 0x6A
	VK_ADD       = 0x6B
	VK_SEPARATOR = 0x6C
	VK_SUBTRACT  = 0x6D
	VK_DECIMAL   = 0x6E
	VK_DIVIDE    = 0x6F
	VK_F1        = 0x70
	VK_F2        = 0x71
	VK_F3        = 0x72
	VK_F4        = 0x73
	VK_F5        = 0x74
	VK_F6        = 0x75
	VK_F7        = 0x76
	VK_F8        = 0x77
	VK_F9        = 0x78
	VK_F10       = 0x79
	VK_F11       = 0x7A
	VK_F12       = 0x7B
	VK_F13       = 0x7C
	VK_F14       = 0x7D
	VK_F15       = 0x7E
	VK_F16       = 0x7F
	VK_F17       = 0x80
	VK_F18       = 0x81
	VK_F19       = 0x82
	VK_F20       = 0x83
	VK_F21       = 0x84
	VK_F22       = 0x85
	VK_F23       = 0x86
	VK_F24       = 0x87
	VK_NUMLOCK   = 0x90
	VK_SCROLL    = 0x91
	// VK_L & VK_R - left and right Alt, Ctrl and Shift virtual keys.
	// Used only as parameters to GetAsyncKeyState() and GetKeyState().
	// No other API or message will distinguish left and right keys in this way.
	VK_LSHIFT   = 0xA0
	VK_RSHIFT   = 0xA1
	VK_LCONTROL = 0xA2
	VK_LCTRL    = 0xA2
	VK_RCONTROL = 0xA3
	VK_RCTRL    = 0xA3
	VK_LMENU    = 0xA4
	VK_RMENU    = 0xA5
)
