package gonable

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/neteng-tools/cliPrompt"
)

type NewPage struct {
    Page *rod.Page
    ShortWait time.Duration
    MedWait time.Duration
    LongWait time.Duration
}

func (n *NewPage) Fill_Defaults() *NewPage {
	if n.ShortWait == 0 {
		n.ShortWait = 300 * time.Millisecond
	}
	if n.MedWait == 0 {
		n.MedWait = 1000 * time.Millisecond
	}
	if n.LongWait == 0 {
		n.LongWait = 3000 * time.Millisecond
	}
	return n
}

func (n *NewPage) Connect(url string) *NewPage {
	n.Page = rod.New().MustConnect().MustPage(url)
	n.Fill_Defaults()
	return n
}

func (n *NewPage) Login(username string, password string) {
    const (
        usernameField string = "#userNameId"
        passwordField string = "#passwordFieldId"
        twoFactorField string = "#twoFactorAuthenticationPasscode"

        nextButton string = "#nextButton_label > span"
        loginButton string = "#loginButton_label"
        submitButton string = "#submitButtonId_label"
		errBox string = "#errBoxId"
		backButton string = "#backButton"
    )
	loginFailed := false
    // Log in to the site.
	for {

		if username == "" || loginFailed {
			username = prompt.Credentials("Username: ")
		}
		
		n.Page.MustElement(usernameField).MustInput(username)
		fmt.Println("Clicking Next")
		n.Page.MustElement(nextButton).MustClick()
		n.Page.MustWaitStable()

		if password == "" || loginFailed {
			password = prompt.Credentials("Password: ")
		}
		n.Page.MustElement(passwordField).MustInput(password)
		n.Page.MustElement(loginButton).MustClick()
		n.Page.MustWaitStable()
		//check for error box after logging in to make sure we can move on.
		err := rod.Try(func() {
			n.Page.Timeout(n.ShortWait).MustElement(errBox)
		})
		if err == nil{
			fmt.Println("Username or Password not accepted")
			n.Page.MustElement(backButton).MustClick()
			loginFailed = true
			continue
		} 
		break
	}
	for {
		userCode := prompt.Credentials("Enter 2FA Code: ")
		n.Page.MustElement(twoFactorField).MustSelectAllText().MustInput(userCode)
		n.Page.MustElement(submitButton).MustClick()
		n.Page.MustWaitStable()
		time.Sleep(n.LongWait)
		if n.Page.MustInfo().Title != "Two-Step Verification" {
			break
		}
	}
	fmt.Println("Logging in...")

}

func (n *NewPage) AllDevicesPage() {
    const allDevicesButton string = "#VIEWS_pane > ul > li:nth-child(2)"
    //gets All Devices. ID includes a space which seems to be invalid so just manually navigating down the tree
	n.Page.MustWaitStable()
    n.Page.MustElement(allDevicesButton).MustClick()
    n.Page.MustWaitStable()
    time.Sleep(n.MedWait)
}

func (n *NewPage) BulkEdit() *NewPage {
	bulkElement, err := n.Page.Element("#bulkEditDevicesTable")
	if err != nil {
		panic(err)
	}
	fmt.Print(bulkElement.Text())
	return n
}

//Edit() > GetDeviceName(). Gets device name from Device Properties page. Panics on multiselect page as that's not currently supported.
//Great for verifying you're in the right place and allows you compare new and old device names before and after making changes.
func (n *NewPage) GetDeviceName() (string, bool) {
	const ( 
		deviceNameLoc string = "#deviceHeaderId > div.xtndDetailedHeaderOverview > div > span.xtndDetailedHeaderTitle"
	)
	n.Page.MustWaitStable()
	time.Sleep(n.MedWait)
	err := rod.Try(func() {
			n.Page.Timeout(10 * time.Second).MustElement(deviceNameLoc)
	})
	n.Page.MustWaitStable()
	if err != nil{
		return "", false
	} 
	deviceName, _ := n.Page.MustElement(deviceNameLoc).Text()
	return deviceName, true
}
func (n *NewPage) Search(searchString string) {
    const (
        searchBox string = "#lanDeviceIndex_searchBox"
        applyFilter string = "#lanDeviceIndex_applyFilter_label"
    )
	n.Page.MustWaitIdle()
	n.Page.MustWaitStable()
    n.Page.MustElement(searchBox).MustSelectAllText().MustInput(searchString)
	n.Page.MustWaitStable()
    n.Page.MustElement(applyFilter).MustClick()
    n.Page.MustWaitStable()
}

//only works on the AllDevices Page. Selects all devices listed. 
//If you searched for a device or set a filter it'll select that one device or group.
func (n *NewPage) SelectAll() *NewPage{
    const (
        selectAllBox string = "#lanDeviceIndexGrid-header > tr > th.dgrid-cell.dgrid-column-0.dgrid-selector-wrapper > div > input[type=checkbox]"
    )
	time.Sleep(n.ShortWait)
    n.Page.MustElement(selectAllBox).MustClick()
    n.Page.MustWaitStable()
    time.Sleep(n.ShortWait)
    return n
}
//AllDevicesPage() > SelectAll() > Edit()
func (n *NewPage) Edit() *NewPage {
	const (
		editButton string = "#lanDeviceIndex > div > span:nth-child(2)"
	)
	n.Page.MustWaitStable()
    n.Page.MustElement(editButton).MustClick()    
    n.Page.MustWaitStable()
    time.Sleep(n.ShortWait)
    return n
}

func (n *NewPage) settings() *NewPage {
    const (
        settingsTab string = "#editLanDeviceTabContainerId_tablist_settingsTabId"
    )
	n.Page.MustWaitStable()
    n.Page.MustElement(settingsTab).MustClick() 
    n.Page.MustWaitStable()
	time.Sleep(n.ShortWait)
    return n
}

func (n *NewPage) settingsProperties() *NewPage {
    const (
        PropertiesTab string = "#settingsNestedTabContainerId_tablist_settingsPropertiesTabId"
    )
	n.Page.MustWaitStable()
    n.Page.MustElement(PropertiesTab).MustClick()    
    n.Page.MustWaitStable()
    return n
}
//Edit() > DeviceProps(). Goes into Settings and clicks Properties.
func (n *NewPage) DeviceProps() *NewPage {
    n.settings().settingsProperties()
    return n
}

func (n *NewPage) discoveredNameCheckBox() *rod.Element {
    const (
        discoveredNameCheckBox string = "#useDiscoveredNameID"
    )
	n.Page.MustWaitStable()
    return n.Page.MustElement(discoveredNameCheckBox)

}

func (n *NewPage) osNameSelect() *NewPage {
    const (
        osNameDropDown = "#widget_deviceOsId > div.dijitReset.dijitRight.dijitButtonNode.dijitArrowButton.dijitDownArrowButton.dijitArrowButtonContainer > input"
		osNameDropOption = "#deviceOsId_popup0"
    )
	n.Page.MustWaitStable()
    n.Page.MustElement(osNameDropDown).MustClick()
	n.Page.MustElement(osNameDropOption).MustClick()
	n.Page.MustWaitStable()

    return n
}

func (n *NewPage) InputOsName(osName string) *NewPage {
    const (
        osNameField = "#deviceOsId"
    )
	n.Page.MustWaitStable()
    n.Page.MustElement(osNameField).MustSelectAllText().MustWaitStable().MustInput(osName)
    time.Sleep(n.ShortWait)
	return n
}
//grabs the ID for the checkbox on the Properties page and makes sure it's unchecked. It also checks if it's currently checked
func (n *NewPage) uncheckUseDiscovered() *NewPage {
	n.Page.MustWaitStable()
    NameBoxChecked, err := n.discoveredNameCheckBox().Property("checked")
    if err != nil {
        panic(err)
    }
    if NameBoxChecked.Bool() {
        n.discoveredNameCheckBox().MustClick()
        time.Sleep(n.ShortWait)
    }
    return n
}

func (n *NewPage) inputDeviceName(name string) *NewPage {
    const (
        deviceNameField = "#deviceNameId"
    )
	n.Page.MustWaitStable()
    n.Page.MustElement(deviceNameField).MustSelectAllText().MustInput(name) 
    time.Sleep(n.ShortWait)
    return n
}
//First checks to make sure the "Use Discovered Name" is unchecked so the Given name box isn't grayed out. 
//After that it selects all the existing characters and replaces with the provided string.
func (n *NewPage) ChangeDeviceName(name string) *NewPage {
    n.uncheckUseDiscovered().inputDeviceName(name)
	n.Page.MustWaitStable()
    return n
}
//Hits save button on device settings page. Saving takes 3-10 seconds so there's a mandatory 4 second delay after hitting the button.
//It might still be worth it to add a longer manual delay to avoid issues.
func (n *NewPage) SaveChanges() *NewPage {
    const (
        saveButton string = "#saveButtonId_label"
    )
    n.Page.MustElement(saveButton).MustClick()
	time.Sleep(n.ShortWait)
    n.Page.MustWaitStable()
	time.Sleep(n.LongWait)
    return n
}
//hits the cancel button the device edit page. General used after saving.
func (n *NewPage) DevicePageCancel() *NewPage {
    const (
        cancelButton string = "#cancelButtonId_label"
    )
    n.Page.MustElement(cancelButton).MustClick()
    n.Page.MustWaitStable()
    time.Sleep(n.MedWait)
    return n
}
//hits the cancel button on the multidevice edit page. 
//Usually the result of an error since we don't handle that page.
func (n *NewPage) MultiDevicePageCancel() *NewPage {
	const (
	multiDeviceCancel string = "#xtnd_form_CancelButton_0_label"
	)
	n.Page.MustElement(multiDeviceCancel).MustClick()
    n.Page.MustWaitStable()
    time.Sleep(n.MedWait)
	return n
}
