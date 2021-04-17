package reza

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	win "github.com/lxn/win"
)

var (
	CLSID_FileSaveDialog  = ole.NewGUID("{C0B4E2F3-BA21-4773-8DBA-335EC946EB8B}")
	IID_IFileDialog       = ole.NewGUID("{42F85136-DB7E-439C-85F1-E4075D135FC8}")
	IID_IFileDialogEvents = ole.NewGUID("{973510db-7d7f-452b-8975-74a85828d354}")
)

type COMDLG_FILTERSPEC struct {
	pszName *uint16
	pszSpec *uint16
}

type IFileSaveDialog struct {
	ole.IUnknown
}

type IFileSaveDialogVtbl struct {
	ole.IUnknownVtbl
	Show                   uintptr
	SetFileTypes           uintptr
	SetFileTypeIndex       uintptr
	GetFileTypeIndex       uintptr
	Advise                 uintptr
	Unadvise               uintptr
	SetOptions             uintptr
	GetOptions             uintptr
	SetDefaultFolder       uintptr
	SetFolder              uintptr
	GetFolder              uintptr
	GetCurrentSelection    uintptr
	SetFileName            uintptr
	GetFileName            uintptr
	SetTitle               uintptr
	SetOkButtonLabel       uintptr
	SetFileNameLabel       uintptr
	GetResult              uintptr
	AddPlace               uintptr
	SetDefaultExtension    uintptr
	Close                  uintptr
	SetClientGuid          uintptr
	ClearClientData        uintptr
	SetFilter              uintptr
	SetSaveAsItem          uintptr
	SetProperties          uintptr
	SetCollectedProperties uintptr
	GetProperties          uintptr
	ApplyProperties        uintptr
}

type IFileDialogEventsVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	// IFileDialogEvents methods
	OnFileOk          uintptr
	OnFolderChange    uintptr
	OnFolderChanging  uintptr
	OnHelp            uintptr
	OnSelectionChange uintptr
	OnShareViolation  uintptr
	OnTypeChange      uintptr
	OnOverwrite       uintptr
}

type IFileDialogEvents struct {
	vtbl IFileDialogEventsVtbl
}

type IFileDialog struct {
}

type IShellItem struct {
}

type FDE_SHAREVIOLATION_RESPONSE struct {
}

type FDE_OVERWRITE_RESPONSE struct {
}

func (v *IFileDialogEvents) Initialize() {
	v.vtbl.QueryInterface = syscall.NewCallback(v.QueryInterface)
	v.vtbl.AddRef = syscall.NewCallback(v.AddRef)
	v.vtbl.Release = syscall.NewCallback(v.Release)
	v.vtbl.OnFileOk = syscall.NewCallback(v.OnFileOk)
	v.vtbl.OnFolderChange = syscall.NewCallback(v.OnFolderChange)
	v.vtbl.OnFolderChanging = syscall.NewCallback(v.OnFolderChanging)
	v.vtbl.OnHelp = syscall.NewCallback(v.OnHelp)
	v.vtbl.OnSelectionChange = syscall.NewCallback(v.OnSelectionChange)
	v.vtbl.OnShareViolation = syscall.NewCallback(v.OnShareViolation)
	v.vtbl.OnTypeChange = syscall.NewCallback(v.OnTypeChange)
	v.vtbl.OnOverwrite = syscall.NewCallback(v.OnOverwrite)
}

func (v *IFileDialogEvents) QueryInterface(iid *ole.GUID) uintptr {
	log.Println("QueryInterface")
	return 0
}

func (v *IFileDialogEvents) AddRef() uintptr {
	log.Println("AddRef")
	return 0
}

func (v *IFileDialogEvents) Release() uintptr {
	log.Println("Release")
	return 0
}

func (v *IFileDialogEvents) OnFileOk(pfd *IFileDialog) uintptr {
	return 0
}

func (v *IFileDialogEvents) OnFolderChange(pfd *IFileDialog) uintptr {
	return 0
}
func (v *IFileDialogEvents) OnFolderChanging(pfd *IFileDialog, psi *IShellItem) uintptr {
	return 0
}
func (v *IFileDialogEvents) OnHelp(pfd *IFileDialog) uintptr {
	return 0
}
func (v *IFileDialogEvents) OnSelectionChange(pfd *IFileDialog) uintptr {
	return 0
}
func (v *IFileDialogEvents) OnShareViolation(pfd *IFileDialog, psi *IShellItem, res *FDE_SHAREVIOLATION_RESPONSE) uintptr {
	return 0
}
func (v *IFileDialogEvents) OnTypeChange(pfd *IFileDialog) uintptr {
	return 0
}
func (v *IFileDialogEvents) OnOverwrite(pfd *IFileDialog, psi *IShellItem, res *FDE_OVERWRITE_RESPONSE) uintptr {
	return 0
}

func (v *IFileSaveDialog) VTable() *IFileSaveDialogVtbl {
	return (*IFileSaveDialogVtbl)(unsafe.Pointer(v.RawVTable))
}

func (self *IFileSaveDialog) Show(hwndOwner win.HWND) win.HRESULT {
	ret, _, _ := syscall.Syscall(
		self.VTable().Show,
		2,
		uintptr(unsafe.Pointer(self)),
		uintptr(hwndOwner), 0)
	return win.HRESULT(ret)
}

func (self *IFileSaveDialog) Advise(pfde *IFileDialogEvents, pdwCookie *uint32) win.HRESULT {
	ret, _, _ := syscall.Syscall(
		self.VTable().Advise,
		3,
		uintptr(unsafe.Pointer(self)),
		uintptr(unsafe.Pointer(pfde)),
		uintptr(unsafe.Pointer(pdwCookie)))
	return win.HRESULT(ret)
}

func (self *IFileSaveDialog) SetFileTypes(cFileTypes uint32, rgFilterSpec *COMDLG_FILTERSPEC) win.HRESULT {
	ret, _, _ := syscall.Syscall(
		self.VTable().SetFileTypes,
		3,
		uintptr(unsafe.Pointer(self)),
		uintptr(cFileTypes),
		uintptr(unsafe.Pointer(rgFilterSpec)))
	return win.HRESULT(ret)
}

func (self *IFileSaveDialog) SetFileName(name string) win.HRESULT {
	nameUTF16, _ := syscall.UTF16PtrFromString(name)
	ret, _, _ := syscall.Syscall(
		self.VTable().SetFileName,
		2,
		uintptr(unsafe.Pointer(self)),
		uintptr(unsafe.Pointer(nameUTF16)),
		0)
	return win.HRESULT(ret)
}

func Foo() uint32 {
	log.Println("Foo")
	return 0
}

func TestNewDialog(window Window) {
	//var dwCookie uint32 = 0
	//var fileEvent *IFileDialogEvents = &IFileDialogEvents{}
	//fileEvent.Initialize()

	unknown, _ := ole.CreateInstance(CLSID_FileSaveDialog, IID_IFileDialog)
	var fileSaveDialog = (*IFileSaveDialog)(unsafe.Pointer(unknown))

	rgFilterSpec := []COMDLG_FILTERSPEC{
		{syscall.StringToUTF16Ptr("JPEG"), syscall.StringToUTF16Ptr("*.jpg;*.jpeg")},
		{syscall.StringToUTF16Ptr("Bitmap Files"), syscall.StringToUTF16Ptr("*.bmp")},
		{syscall.StringToUTF16Ptr("PNG"), syscall.StringToUTF16Ptr("*.png")},
		{syscall.StringToUTF16Ptr("All Files"), syscall.StringToUTF16Ptr("*.*")},
	}

	//fileSaveDialog.Advise(fileEvent, &dwCookie)

	fileSaveDialog.SetFileName("Untitled.png")
	fileSaveDialog.SetFileTypes(4, &rgFilterSpec[0])

	fileSaveDialog.Show(window.GetHandle())

	fileSaveDialog.Release()

}
