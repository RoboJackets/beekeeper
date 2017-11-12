package main

import (
	"fmt"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os/exec"
	"regexp"
	"smoker/backends"
	"sort"
	"strconv"
	"strings"
	// "smoker/backends"
)

// * Ui Functions

func uniq(strings []string) (ret []string) {
	return
}

func authors() []string {
	if b, err := exec.Command("git", "log").Output(); err == nil {
		lines := strings.Split(string(b), "\n")

		var a []string
		r := regexp.MustCompile(`^Author:\s*([^ <]+).*$`)
		for _, e := range lines {
			ms := r.FindStringSubmatch(e)
			if ms == nil {
				continue
			}
			a = append(a, ms[1])
		}
		sort.Strings(a)
		var p string
		lines = []string{}
		for _, e := range a {
			if p == e {
				continue
			}
			lines = append(lines, e)
			p = e
		}
		return lines
	}
	return []string{"The Dark Lord, Sauron"}
}

func startUi(backend backends.Backend, credManager backends.CredentialManager, cred backends.Credential) {
	var menuitem *gtk.MenuItem
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("Varroa")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("bye bye!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)

	//--------------------------------------------------------
	// GtkMenuBar
	//--------------------------------------------------------
	menubar := gtk.NewMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	//--------------------------------------------------------
	// GtkVPaned
	//--------------------------------------------------------
	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	//--------------------------------------------------------
	// GtkFrame
	//--------------------------------------------------------
	frame1 := gtk.NewFrame("Scan Interface")
	framebox1 := gtk.NewVBox(false, 1)
	frame1.Add(framebox1)

	frame2 := gtk.NewFrame("Dump Interface")
	framebox2 := gtk.NewVBox(false, 1)
	frame2.Add(framebox2)

	vpaned.Pack1(frame1, false, false)
	vpaned.Pack2(frame2, false, false)

	label := gtk.NewLabel("Varroa Inventory Client")
	label.ModifyFontEasy("DejaVu Sans 15")
	framebox1.PackStart(label, false, true, 0)

	//--------------------------------------------------------
	// GtkEntry
	//--------------------------------------------------------
	entry := gtk.NewEntry()
	entry.SetText("Scanner Interface")
	framebox1.Add(entry)

	// ** Output of Scanner Interface
	label2 := gtk.NewLabel("Please Scan An Item!")
	label2.ModifyFontEasy("DejaVu Sans 15")
	framebox1.Add(label2)

	entry.Connect("activate", func() {
		handleScan(entry.GetText(), label2, backend)
	})

	//--------------------------------------------------------
	// GtkVSeparator
	//--------------------------------------------------------
	vsep := gtk.NewVSeparator()
	framebox2.PackStart(vsep, false, false, 0)

	//--------------------------------------------------------
	// Search Box
	//--------------------------------------------------------
	combos := gtk.NewHBox(false, 1)
	searchBox := gtk.NewEntry()
	combos.Add(searchBox)

	// Update button
	updateButton := gtk.NewButtonWithLabel("Update Stuff")
	combos.Add(updateButton)

	framebox2.PackStart(combos, false, false, 0)

	//--------------------------------------------------------
	// Dump of items
	//--------------------------------------------------------

	swin := gtk.NewScrolledWindow(nil, nil)
	store := gtk.NewListStore(glib.G_TYPE_STRING,
		glib.G_TYPE_STRING,
		glib.G_TYPE_STRING,
		glib.G_TYPE_STRING,
		glib.G_TYPE_STRING)
	treeview := gtk.NewTreeView()
	swin.Add(treeview)
	treeview.SetModel(store)
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Id", gtk.NewCellRendererText(), "text", 0))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Bin", gtk.NewCellRendererText(), "text", 1))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Name", gtk.NewCellRendererText(), "text", 2))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Manufacturer", gtk.NewCellRendererText(), "text", 3))
	treeview.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Count", gtk.NewCellRendererText(), "text", 4))
	framebox2.Add(swin)

	updateButton.Clicked(func () {
		updateListView(store, searchBox.GetText(), backend)
	})


	//--------------------------------------------------------
	// GtkMenuItem
	//--------------------------------------------------------
	cascademenu := gtk.NewMenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("E_xit")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_View")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	checkmenuitem := gtk.NewCheckMenuItemWithMnemonic("_Disable")
	checkmenuitem.Connect("activate", func() {
		vpaned.SetSensitive(!checkmenuitem.GetActive())
	})
	submenu.Append(checkmenuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("_About")
	menuitem.Connect("activate", func() {
		dialog := gtk.NewAboutDialog()
		dialog.SetName("Varroa Info")
		dialog.SetProgramName("Varroa")
		dialog.SetAuthors(authors())
		dialog.SetLicense("The library is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License). The idea is that if you can use Go (and Gtk) in a project, you should also be able to use go-gtk.")
		dialog.SetWrapLicense(true)
		dialog.Run()
		dialog.Destroy()
	})
	submenu.Append(menuitem)

	//--------------------------------------------------------
	// GtkStatusbar
	//--------------------------------------------------------
	statusbar := gtk.NewStatusbar()
	context_id := statusbar.GetContextId("go-gtk")
	statusbar.Push(context_id, "(and (> qt gtk ) (terrible-p 'golang))")

	framebox2.PackStart(statusbar, false, false, 0)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.Add(vbox)
	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
}

// * Data Helper Functions

func handleScan(id string, label *gtk.Label, b backends.Backend) {

	if component, bin, err := b.LookupId(id); err != nil {
		messagedialog := gtk.NewMessageDialog(
			label.GetTopLevelAsWindow(),
			gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT,
			gtk.MESSAGE_QUESTION,
			gtk.BUTTONS_OK,
			"Info On: "+id)

		nameEntry := gtk.NewEntry()
		nameEntry.SetText("Part Name")
		manEntry := gtk.NewEntry()
		manEntry.SetText("Part Manufacturer")
		countEntry := gtk.NewEntry()
		countEntry.SetText("Part Count")

		messagedialog.GetVBox().Add(nameEntry)
		messagedialog.GetVBox().Add(manEntry)
		messagedialog.GetVBox().Add(countEntry)

		messagedialog.Response(func() {
			if idInt, err := strconv.ParseUint(countEntry.GetText(), 10, 32); err != nil {
				fmt.Println("COUNT WAS INVALID")
			} else {
				component := backends.NewComponent(id, uint(idInt), nameEntry.GetText(), manEntry.GetText())
				if _, err := b.AddComponent(component); err != nil {
					fmt.Println("ERROR ADDING COMPONENT")
				}
			}
			messagedialog.Destroy()
		})

		messagedialog.ShowAll()
		messagedialog.Run()
	} else {
		label.SetText("Bin: " + bin.GetName() + "\nCount: " + strconv.Itoa(int(component.GetCount())))
	}
}

func updateListView(list *gtk.ListStore, searchCriteria string, b backends.Backend) {
	list.Clear()

	var c []backends.Component;

	if len(searchCriteria) == 0 {
		// Dump everything
		c = b.GetAllComponents()
	} else {
		// Narrow on search results
		c = b.GeneralSearch(searchCriteria)
	}
	for _, v := range c {
		var iter gtk.TreeIter
		list.Append(&iter)
		list.Set(&iter, 0, v.GetId(), 1, v.GetBin(), 2, v.GetName(), 3, v.GetManufacturer(), 4, strconv.Itoa(int(v.GetCount())))
	}
}
