package forms

import (
	"context"
	"database/sql"
	"<projectrepo>/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"net/http"
)

var ErrValidationFailed = errors.Errorf("Form validation failed")

type Form interface {
	ShouldBind(c *gin.Context) error
	RenderForm(c *gin.Context)
	AddError(fieldName string, message string)
	Save(c context.Context, exec boil.ContextExecutor) error
	TemplateData() map[string]interface{}

	// the only thing to implement in the child form
	Validate(c *gin.Context, exec boil.ContextExecutor) error
}

type FormBase[T any] struct {
	Name                 string
	FormTemplate         string
	FormSaved            bool
	KeepValuesAfterSave  bool
	FullPageReloadOnSave bool

	// we make no effort to sync the keys names between
	// errors and struct, we just assume it to be matching.
	// Same goes for the templates - it's completely left
	// to the developer
	Errors            map[string]string
	Input             *T
	ExtraTemplateData map[string]interface{}
}

func (f *FormBase[T]) ShouldBind(c *gin.Context) error {
	err := c.ShouldBind(&f.Input)

	if err != nil {
		return errors.Wrapf(err, "failed to bind input to the form [%s]", f.Name)
	}

	return nil
}

func (f *FormBase[T]) TemplateData() map[string]interface{} {
	data := map[string]interface{}{
		"Input":     f.Input,
		"Errors":    f.Errors,
		"FormSaved": f.FormSaved,
	}

	for k, v := range f.ExtraTemplateData {
		data[k] = v
	}

	return data
}

func (f *FormBase[T]) AddError(fieldName string, message string) {
	if f.Errors == nil {
		f.Errors = map[string]string{}
	}

	f.Errors[fieldName] = message
}

func (f *FormBase[T]) RenderForm(c *gin.Context) {
	if f.FormSaved && f.FullPageReloadOnSave {
		c.Header("HX-Refresh", "true")
		c.Status(http.StatusOK)
		return
	}

	c.HTML(http.StatusOK, f.FormTemplate, f.TemplateData())
}

func (f *FormBase[T]) Save(c context.Context, exec boil.ContextExecutor) error {
	f.FormSaved = true
	if !f.KeepValuesAfterSave {
		// reset input values to get a pristine form
		f.Input = new(T)
	}

	return nil
}

func DefaultHandler(c *gin.Context, db *sqlx.DB, form Form) {
	if err := form.ShouldBind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"explanation": "Failed to process request", "err": err.Error()})
		return
	}

	if err := form.Validate(c, db); err != nil {
		form.RenderForm(c)
		return
	}

	if err := util.Transact(db, func(tx *sql.Tx) error {
		return form.Save(c, tx)
	}); err != nil {
		log.Panicf("Failed to save the form: %v", err)
	}

	form.RenderForm(c)
}
