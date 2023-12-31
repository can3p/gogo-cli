package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	_ "time/tzdata" // help go learn about timezones

	"<projectrepo>/pkg/admin"
	"<projectrepo>/pkg/auth"
	"<projectrepo>/pkg/forms"
	"<projectrepo>/pkg/forms/validation"
	"<projectrepo>/pkg/links"
	"<projectrepo>/pkg/mail"
	"<projectrepo>/pkg/markdown"
	"<projectrepo>/pkg/model/core"
	"<projectrepo>/pkg/pgsession"
	"<projectrepo>/pkg/util"
	"<projectrepo>/pkg/web"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres db driver
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var staticRoute = "/static"

// replace with cdn domain
var staticRouteCluster = "/static"

var requiredVars = []string{
	"DATABASE_URL",
	"SESSION_SALT",
	"MJ_APIKEY_PRIVATE",
	"MJ_APIKEY_PUBLIC",
	"SITE_ROOT",
}

func enforceEnvVars() {
	for _, v := range requiredVars {
		if _, ok := os.LookupEnv(v); !ok {
			panic(fmt.Sprintf("var %s is not set", v))
		}
	}
}

func main() {
	enforceEnvVars()

	// fly.io does not have sslmode enabled
	db := sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()

	sessionSalt := os.Getenv("SESSION_SALT")

	store := pgsession.NewStore(db, []byte(sessionSalt))

	// developer timezone only messes things up
	time.Local = time.UTC

	r := gin.Default()

	r.Use(sessions.Sessions("sess", store))
	r.Use(func(c *gin.Context) { auth.Auth(c, db) })

	if util.InCluster() {
		r.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
			userData := auth.GetUserData(c)
			user := userData.DBUser

			admin.NotifyPageFailure(c, err, user)
		}))
	} else {
		log.Println("Custom error reporter skipped")
	}

	html := flag.String("html", "client/html", "path to html templates")

	flag.Parse()

	r.SetFuncMap(funcmap())
	r.LoadHTMLGlob(fmt.Sprintf("%s/*.html", *html))

	r.GET("/", func(c *gin.Context) {
		userData := auth.GetUserData(c)

		if userData.IsLoggedIn {
			c.Redirect(http.StatusFound, "/controls")
			return
		}

		c.HTML(http.StatusOK, "index.html", web.Index(c, db, &userData))
	})

	//cache static forever
	if util.InCluster() {
		r.Group("/static", func(c *gin.Context) {
			c.Header("cache-control", "max-age=31536000, public")
			c.Next()
		}).Static("/", "dist")
	} else {
		r.Group("/static").Static("/", "dist")
	}

	r.GET("/articles/:id", func(c *gin.Context) {
		articleName := c.Param("id")

		fname := fmt.Sprintf("client/articles/%s.md", articleName)

		body, err := ioutil.ReadFile((fname))

		if err != nil {
			panic(err)
		}

		lines := util.SplitLines(string(body))

		title := lines[0]
		signupAttribution := lines[1]
		sbody := strings.TrimSpace(strings.Join(lines[2:], "\n"))

		userData := auth.GetUserData(c)
		c.HTML(http.StatusOK, "article.html", gin.H{
			"Name":        title,
			"Body":        sbody,
			"User":        userData,
			"Attribution": signupAttribution,
		})
	})

	r.POST("/form/signup", func(c *gin.Context) {
		var input struct {
			Email       string `form:"email" binding:"required"`
			Password    string `form:"password" binding:"required"`
			Attribution string `form:"attribution"`
		}

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"explanation": "Failed to process request", "err": err.Error()})
			return
		}

		systemSettings := core.SystemSettings().OneP(c, db)

		if !systemSettings.RegistrationOpen {
			c.Status(http.StatusForbidden)
			return
		}

		errors := map[string]string{}
		hasErrors := false

		email := strings.TrimSpace(strings.ToLower(input.Email))

		if email == "" {
			errors["email"] = "This email cannot be empty"
			hasErrors = true
		} else if reason, isOK := web.EmailOKToSignup(c, db, input.Email); !isOK {
			errors["email"] = reason
			hasErrors = true
		}

		if err := validation.ValidatePassword(input.Password); err != nil {
			errors["password"] = err.Error()
			hasErrors = true
		}

		if hasErrors {
			c.HTML(http.StatusOK, "form--signup.html", map[string]interface{}{
				"Errors": errors,
				"Input":  input,
			})
			return
		}

		attribution := input.Attribution

		// we're not enforcing a specific enum of attributions
		// since it's just additional work at the moment
		if !web.AttributionRE.MatchString(attribution) {
			attribution = "unknown"
		}

		if len(attribution) > 100 {
			attribution = attribution[0:100]
		}

		// we're getting the input from the form and not
		if len(attribution) > 100 {
			attribution = attribution[0:100]
		}

		user, err := auth.Signup(c, db, email, input.Password, attribution)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"explanation": "Failed to sign up", "err": err.Error()})
			return
		}

		if err := mail.ConfirmSignup(db, user); err != nil {
			panic(err)
		}

		c.HTML(http.StatusOK, "partial--signup-goto-email.html", map[string]interface{}{})
	})

	r.POST("/login", func(c *gin.Context) {
		var input struct {
			Email    string `form:"email" binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"explanation": "Failed to process request", "err": err.Error()})
			return
		}

		errors := map[string]string{}
		hasErrors := false

		if err := auth.Login(c, db, input.Email, input.Password); err != nil {
			errors["email"] = fmt.Sprintf("Failed to login: %v", err)
			hasErrors = true
		}

		if hasErrors {
			c.HTML(http.StatusOK, "form--login.html", map[string]interface{}{
				"Errors": errors,
				"Input":  input,
			})
			return
		}

		// needs for htmx redirect
		c.Header("HX-Redirect", "/controls/")
		c.Status(http.StatusOK)
	})

	r.GET("/logout", auth.Logout)

	controls := r.Group("/controls", auth.EnforceAuth)
	actions := controls.Group("/action")

	setupActions(actions, db)

	controls.GET("/", func(c *gin.Context) {
		userData := auth.GetUserData(c)

		c.HTML(http.StatusOK, "controls.html", web.Controls(c, db, &userData))
	})

	controls.GET("/settings", func(c *gin.Context) {
		userData := auth.GetUserData(c)

		c.HTML(http.StatusOK, "settings.html", web.Settings(c, db, &userData))
	})

	r.GET("/confirm_signup/:id", func(c *gin.Context) {
		id := c.Param("id")
		userData := auth.GetUserData(c)

		if userData.IsLoggedIn {
			c.Redirect(http.StatusFound, "/controls")
		}

		user, err := core.Users(
			core.UserWhere.EmailConfirmSeed.EQ(null.StringFrom(id)),
		).One(c, db)

		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			panic(err)
		}

		if !user.EmailConfirmedAt.Valid {
			user.EmailConfirmedAt = null.TimeFrom(time.Now())
			user.UpdateP(c, db, boil.Infer())

			go admin.NotifySignupConfirmed(user)
		}

		c.HTML(http.StatusOK, "signup_confirmed.html", map[string]interface{}{
			"User": userData,
		})
	})

	controls.POST("/form/save_settings", func(c *gin.Context) {
		userData := auth.GetUserData(c)
		dbUser := userData.DBUser

		form := forms.SettingsGeneralFormNew(dbUser)

		forms.DefaultHandler(c, db, form)
	})

	controls.POST("/form/change_password", func(c *gin.Context) {
		userData := auth.GetUserData(c)
		dbUser := userData.DBUser

		form := forms.ChangePasswordFormNew(dbUser)

		forms.DefaultHandler(c, db, form)
	})

	if err := r.Run(); err != nil {
		panic(err)
	}
}

func funcmap() template.FuncMap {
	return template.FuncMap{
		"static_asset": func() func(n string) string {
			manifest, err := ioutil.ReadFile("dist/manifest.json")

			if err != nil {
				panic(err)
			}

			files := map[string]string{}

			err = json.Unmarshal(manifest, &files)

			if err != nil {
				panic(err)
			}

			return func(n string) string {
				path, ok := files[n]

				if !ok {
					panic(fmt.Sprintf("asset [%s] is not defined", n))
				}

				prefix := staticRoute

				//if util.InCluster() {
				//prefix = staticRouteCluster
				//}

				return fmt.Sprintf("%s/%s", prefix, path)
			}
		}(),

		"link": links.Link,

		"abslink": links.AbsLink,

		"renderTimestamp": func(t time.Time, user *core.User) string {
			if user != nil {
				t = localizeTime(user, t)
			}

			return t.Format("Mon, 02 Jan 2006 15:04")
		},

		"toMap": func(args ...interface{}) map[string]interface{} {
			if len(args)%2 != 0 {
				panic("toMap got uneven number of arguments")
			}

			out := map[string]interface{}{}

			idx := 0

			for idx+1 < len(args) {
				k := args[idx].(string)

				out[k] = args[idx+1]
				idx += 2
			}

			return out
		},

		"markdown": func(s string) template.HTML {
			return markdown.ToTemplate(s)
		},

		"tzlist": func() []string {
			return util.TimeZones
		},
	}
}

func localizeTime(user *core.User, t time.Time) time.Time {
	l, err := time.LoadLocation(user.Timezone)

	if err != nil {
		log.Printf("failed to parse timezone setting: [%s] - %v", user.Timezone, err)
		return t
	}

	return t.In(l)
}
