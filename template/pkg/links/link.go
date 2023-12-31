package links

import "<projectrepo>/pkg/util"

func Link(name string, args ...string) string {
	switch name {
	case "controls":
		return "/controls"
	case "settings":
		return "/controls/settings"
	case "article":
		return "/articles/" + args[0]
	case "form_signup":
		return "/form/signup"
	case "form_login":
		return "/login"
	case "confirm_signup":
		return "/confirm_signup/" + args[0]
	case "form_save_settings":
		return "/controls/form/save_settings"
	case "form_change_password":
		return "/controls/form/change_password"
	case "logout":
		return "/logout"
	}

	return ""
}

func AbsLink(name string, args ...string) string {
	return util.SiteRoot() + Link(name, args...)
}
