package forms

import (
	"context"
	"fmt"
	"<projectrepo>/pkg/model/core"
	"<projectrepo>/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type SettingsGeneralFormInput struct {
	Timezone string `form:"timezone"`
}

type SettingsGeneralForm struct {
	*FormBase[SettingsGeneralFormInput]
	User *core.User
}

func SettingsGeneralFormNew(u *core.User) Form {
	var form Form = &SettingsGeneralForm{
		FormBase: &FormBase[SettingsGeneralFormInput]{
			Name:                "settings_general",
			FormTemplate:        "form--settings-general.html",
			KeepValuesAfterSave: true,
			Input:               &SettingsGeneralFormInput{},
			ExtraTemplateData: map[string]interface{}{
				"User": u,
			},
		},
		User: u,
	}

	return form
}

func (f *SettingsGeneralForm) Validate(c *gin.Context, db boil.ContextExecutor) error {
	if f.Input.Timezone == "" {
		f.AddError("timezone", "timezone is required")
		return ErrValidationFailed
	}

	found := false
	for _, tz := range util.TimeZones {
		if tz == f.Input.Timezone {
			found = true
			break
		}
	}

	if !found {
		f.AddError("timezone", fmt.Sprintf("Cannot find the timezone [%s]", f.Input.Timezone))
		return ErrValidationFailed
	}

	return nil
}

func (f *SettingsGeneralForm) Save(c context.Context, exec boil.ContextExecutor) error {
	f.User.Timezone = f.Input.Timezone

	if _, err := f.User.Update(c, exec, boil.Whitelist(
		core.UserColumns.Timezone,
		core.UserColumns.UpdatedAt,
	)); err != nil {
		return errors.Wrapf(err, "failed to save to the db")
	}

	return f.FormBase.Save(c, exec)
}
