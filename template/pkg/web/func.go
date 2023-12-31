package web

import (
	"context"

	"<projectrepo>/pkg/auth"
	"<projectrepo>/pkg/model/core"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Index(c context.Context, db boil.ContextExecutor, userData *auth.UserData) map[string]interface{} {
	systemSettings := core.SystemSettings().OneP(c, db)

	return map[string]interface{}{
		"Name": "Super cool <projectname>",
		"User": userData,
	}
}

func Controls(c context.Context, db boil.ContextExecutor, userData *auth.UserData) map[string]interface{} {
	return map[string]interface{}{
		"Name": "Controls",
		"User": userData,
	}
}

func Settings(c context.Context, db boil.ContextExecutor, userData *auth.UserData) map[string]interface{} {
	return map[string]interface{}{
		"Name": "Settings",
		"User": userData,
	}
}
