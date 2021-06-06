package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"
)

const Checked = "checked"
const Timeout = 30 * time.Second

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	//ctx, cancel = chromedp.NewExecAllocator(ctx)
	//defer cancel()

	// create chrome instance
	ctx, _ = chromedp.NewContext(
		ctx,
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	opts, err := usage()
	if err != nil {
		logrus.Fatal(err)
	}

	password, _ := opts.String("--pass")
	if err := login(ctx, password); err != nil {
		logrus.Fatal(err)
	}

	enable, _ := opts.Bool("enable")
	disable, _ := opts.Bool("disable")

	if enable != disable {
		logrus.WithField("enable", enable).Info("changing wifi status")
		if err := setWlEnabled(ctx, enable); err != nil {
			logrus.WithError(err).Fatal("failed changing wifi status")
		}
	}
	isEnabled, err := isWlEnabled(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("cannot determine wireless enablement status")
	}

	logrus.WithField("enabled", isEnabled).Info("wifi enable status")
}

func usage() (docopt.Opts, error) {
	doc := `Enable or disable T8200M WiFi radio

Usage:
	wifi --pass=<password>
	wifi --pass=<password> [enable|disable]
    wifi -h | --help

Options:
    -h --help           This screen
    --pass=<password>   (Required) Admin Password
`
	return docopt.ParseDoc(doc)
}

func selectorById(id, tag string) string {
	if len(tag) == 0 {
		tag = "*"
	}
	return fmt.Sprintf(`%s#%s`, tag, id)
}

func login(ctx context.Context, password string) error {
	BaseUrl := "http://10.25.73.1"
	SetupUrl := BaseUrl + "/wirelesssetup_basic.html"

	userNameSel := selectorById("admin_user_name", "input")
	passwordSel := selectorById("admin_password", "input")
	loginBtnSel := selectorById("btn_login", "a")
	footerSel := selectorById("footer_homepage", "div")
	loginFormSel := selectorById("login_form", "div")
	var loginFormStyle string
	var ok bool

	err := chromedp.Run(ctx,
		chromedp.Navigate(BaseUrl),
		chromedp.WaitVisible(passwordSel),
		chromedp.SendKeys(userNameSel, "admin"),
		chromedp.SendKeys(passwordSel, password),
		chromedp.Click(loginBtnSel, chromedp.NodeVisible),
		chromedp.WaitVisible(footerSel),
		chromedp.AttributeValue(loginFormSel, "style", &loginFormStyle, &ok),
		chromedp.Navigate(SetupUrl),
	)
	if err != nil {
		return err
	}
	if !strings.Contains(loginFormStyle, "display: none") {
		return errors.New("login failed")
	}
	return nil
}

func wlRadioSel(enable bool) string {
	onOff := "off"
	if enable {
		onOff = "on"
	}
	return selectorById("id_wl_"+onOff, "input")
}

func isWlEnabled(ctx context.Context) (bool, error) {
	enableWlSel := wlRadioSel(true)
	disableWlSel := wlRadioSel(false)
	footerSel := selectorById("footer", "div")
	var isWlEnabled, isWlDisabled bool

	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(footerSel),
		chromedp.JavascriptAttribute(enableWlSel, Checked, &isWlEnabled),
		chromedp.JavascriptAttribute(disableWlSel, Checked, &isWlDisabled),
	); err != nil {
		return false, err
	}

	if isWlEnabled == isWlDisabled {
		return false, errors.New("enabled==disabled")
	}
	return isWlEnabled, nil
}

func setWlEnabled(ctx context.Context, enable bool) error {
	sel := wlRadioSel(enable)
	footerSel := selectorById("footer", "div")
	applyBtnSel := selectorById("btn_apply", "a")

	wlSettingSel := selectorById("id_wl_settings", "div")
	getWaiter := func(enable bool) chromedp.Action {
		f := chromedp.WaitNotVisible
		if enable {
			f = chromedp.WaitVisible
		}
		return f(wlSettingSel)
	}

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(footerSel),
		chromedp.Click(sel),
		getWaiter(enable),
		chromedp.Click(applyBtnSel),
	)
	return err
}
