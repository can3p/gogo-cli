package web

import (
	"context"
	"regexp"
	"strings"

	disposable "github.com/can3p/anti-disposable-email"
	"<projectrepo>/pkg/admin"
	"<projectrepo>/pkg/model/core"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var EmailRE *regexp.Regexp = regexp.MustCompile(`(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`)
var TestEmailRE *regexp.Regexp = regexp.MustCompile(`<testemailhead>(\+[^@]+)?@<testemailtail>`)
var AttributionRE *regexp.Regexp = regexp.MustCompile(`[a-z_]+`)

func EmailOKToSignup(ctx context.Context, db boil.ContextExecutor, address string) (string, bool) {
	if !EmailRE.MatchString(address) {
		return "Invalid email", false
	}

	if core.Users(core.UserWhere.Email.EQ(address)).ExistsP(ctx, db) {
		return "Email is already used in the system", false
	}

	if strings.Contains(address, "+") && !TestEmailRE.MatchString(address) {
		return "Plus sign is not allowed in the emails", false
	}

	parsedEmail, _ := disposable.ParseEmail(address)

	if parsedEmail.Disposable {
		go admin.NotifyThrowAwayEmailSignupAttempt(address)

		return "Email domain is not allowed, please reach out to us via the support form", false
	}

	return "", true
}
