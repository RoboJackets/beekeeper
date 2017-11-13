package main

import (
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"smoker/backends"
	"strconv"
)

// * Ui Functions

// ** Helper Functions
func authors() []string {
	return []string{"Jay Kamat (jgkamat)"}
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
	}, "")

	// ** High Level Structure
	vbox := gtk.NewVBox(false, 1)

	menubar := gtk.NewMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

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

	// ** Scanner Entry Field

	entry := gtk.NewEntry()
	entry.SetText("123456")
	framebox1.Add(entry)

	// ** Output of Scanner Interface
	label2 := gtk.NewLabel("Please Scan An Item!")
	label2.ModifyFontEasy("DejaVu Sans 15")
	framebox1.Add(label2)

	entry.Connect("activate", func() {
		handleScan(entry.GetText(), label2, backend)
	})

	vsep := gtk.NewVSeparator()
	framebox2.PackStart(vsep, false, false, 0)

	combos := gtk.NewHBox(false, 1)
	searchBox := gtk.NewEntry()
	combos.Add(searchBox)

	// ** Update Button
	updateButton := gtk.NewButtonWithLabel("Update Dump")
	combos.Add(updateButton)

	framebox2.PackStart(combos, false, false, 0)

	// ** Dump of Items
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

	updateButton.Clicked(func() {
		updateListView(store, searchBox.GetText(), backend)
	})
	searchBox.Connect("activate", func() {
		updateListView(store, searchBox.GetText(), backend)
	})

	// ** Alt Menu
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
		dialog.SetLicense(
			`This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.`)
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

	// DUMMY ADD ITEMS
	generateDummyItems(backend)

	updateListView(store, searchBox.GetText(), backend)

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
			if count, err := strconv.ParseUint(countEntry.GetText(), 10, 32); err != nil {
				showError("'"+countEntry.GetText()+"' was not a valid count.", label)
			} else {
				component := backends.NewComponent(id, uint(count), nameEntry.GetText(), manEntry.GetText())
				if _, err := b.AddComponent(component); err != nil {
					showError("We were unable to add this component.", label)
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

func generateDummyItems(b backends.Backend) {
	c := backends.NewComponent("100", 124, "Capacitor", "DigiKey")
	b.AddComponent(c)
	c = backends.NewComponent("200", 342325, "Resistor", "DigiKey")
	b.AddComponent(c)
	c = backends.NewComponent("300", 10, "IC-1", "DigiKey")
	b.AddComponent(c)
	c = backends.NewComponent("400", 20, "IC-2", "DigiKey")
	b.AddComponent(c)
	c = backends.NewComponent("500", 3, "Soldering Iron", "MicroMark")
	b.AddComponent(c)
	c = backends.NewComponent("540", 3, "Soldering Iron Mark 2", "Marky Mark")
	b.AddComponent(c)
	c = backends.NewComponent("600", 1, "KSP", "Squad")
	b.AddComponent(c)
	c = backends.NewComponent("700", 9999999, "Magit", "Tarsius")
	b.AddComponent(c)
	c = backends.NewComponent("800", 30, "IC-3", "DigiKey")
	b.AddComponent(c)
	c = backends.NewComponent("900", 33, "IC-4", "DigiKey")
	b.AddComponent(c)
	c = backends.NewComponent("1000", 5000, "qutebrowser", "The-Compiler")
	b.AddComponent(c)
	c = backends.NewComponent("1100", 1, "powder toy", "everyone")
	b.AddComponent(c)
	c = backends.NewComponent("1200", 100, "Leaded Solder", "MicroMark")
	b.AddComponent(c)
	c = backends.NewComponent("1300", 0, "Lead Free Solder", "MicroMark")
	b.AddComponent(c)
	c = backends.NewComponent("1400", 214, "Batteries", "Duracell")
	b.AddComponent(c)
}

func updateListView(list *gtk.ListStore, searchCriteria string, b backends.Backend) {
	list.Clear()

	var c []backends.Component

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

func showError(text string, parent *gtk.Label) {
	messagedialog := gtk.NewMessageDialog(
		parent.GetTopLevelAsWindow(),
		gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		text)

	messagedialog.Response(func() {
		messagedialog.Destroy()
	})

	messagedialog.ShowAll()
	messagedialog.Run()
}
