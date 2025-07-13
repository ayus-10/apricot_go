package auth

import (
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var (
	googleClientId     = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleCallbackUrl  = os.Getenv("GOOGLE_CALLBACK_URL")
	sessionKey         = os.Getenv("SESSION_KEY")
	maxAge             = func() int {
		v, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
		if err != nil {
			return 86400 * 30
		}
		return v
	}()
	isProd = os.Getenv("ENVIRONMENT") == "production"
)

func NewAuth() {
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.MaxAge(maxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, googleCallbackUrl),
	)
}
