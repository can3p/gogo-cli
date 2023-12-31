package forms

import (
	"context"
	"<projectrepo>/pkg/forms/validation"
	"<projectrepo>/pkg/model/core"
	"<projectrepo>/pkg/pgsession"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ChangePasswordFormInput struct {
	OldPassword string `form:"old_password"`
	Password    string `form:"password"`
}

type ChangePasswordForm struct {
	*FormBase[ChangePasswordFormInput]
	User *core.User
}

func ChangePasswordFormNew(u *core.User) Form {
	var form Form = &ChangePasswordForm{
		FormBase: &FormBase[ChangePasswordFormInput]{
			Name:         "change_password",
			FormTemplate: "form--settings-change-password.html",
			Input:        &ChangePasswordFormInput{},
		},
		User: u,
	}

	return form
}

func (f *ChangePasswordForm) Validate(c *gin.Context, db boil.ContextExecutor) error {
	if f.Input.Password == "" {
		f.AddError("password", "password is required")
		return ErrValidationFailed
	}

	if f.Input.OldPassword == "" {
		f.AddError("old_password", "old password is required")
		return ErrValidationFailed
	}

	if err := validation.ValidatePassword(f.Input.Password); err != nil {
		f.AddError("password", err.Error())
		return ErrValidationFailed
	}

	h := pgsession.HashUserPwd(f.User.Email, f.Input.OldPassword)

	if h == "" {
		return errors.Errorf("Failed to has user password, cannot proceed")
	}

	if h != f.User.Pwdhash.String {
		f.AddError("old_password", "old password is not correct")
		return ErrValidationFailed
	}

	return nil
}

func (f *ChangePasswordForm) Save(c context.Context, exec boil.ContextExecutor) error {
	f.User.Pwdhash = null.StringFrom(pgsession.HashUserPwd(f.User.Email, f.Input.Password))

	if _, err := f.User.Update(c, exec, boil.Infer()); err != nil {
		return errors.Wrapf(err, "failed to save to the db")
	}

	return f.FormBase.Save(c, exec)
}
