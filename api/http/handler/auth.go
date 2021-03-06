package handler

import (
	"github.com/akash-Click2cloud/dolphin/api"

	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	httperror "github.com/akash-Click2cloud/dolphin/api/http/error"
	"github.com/akash-Click2cloud/dolphin/api/http/security"
)

// AuthHandler represents an HTTP API handler for managing authentication.
type AuthHandler struct {
	*mux.Router
	Logger          *log.Logger
	authDisabled    bool
	UserService     dockm.UserService
	CryptoService   dockm.CryptoService
	JWTService      dockm.JWTService
	LDAPService     dockm.LDAPService
	SettingsService dockm.SettingsService
}

const (
	// ErrInvalidCredentialsFormat is an error raised when credentials format is not valid
	ErrInvalidCredentialsFormat = dockm.Error("Invalid credentials format")
	// ErrInvalidCredentials is an error raised when credentials for a user are invalid
	ErrInvalidCredentials = dockm.Error("Invalid credentials")
	// ErrAuthDisabled is an error raised when trying to access the authentication endpoints
	// when the server has been started with the --no-auth flag
	ErrAuthDisabled = dockm.Error("Authentication is disabled")
)

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(bouncer *security.RequestBouncer, authDisabled bool) *AuthHandler {
	h := &AuthHandler{
		Router:       mux.NewRouter(),
		Logger:       log.New(os.Stderr, "", log.LstdFlags),
		authDisabled: authDisabled,
	}
	h.Handle("/auth",
bouncer.PublicAccess(http.HandlerFunc(h.handlePostAuth))).Methods(http.MethodPost)

	return h
}

type (
	postAuthRequest struct {
		Username string `valid:"required"`
		Password string `valid:"required"`
	}

postAuthResponse struct {
		JWT string `json:"jwt"`
	}
)

func (handler *AuthHandler) handlePostAuth(w http.ResponseWriter, r *http.Request) {

	if handler.authDisabled {
		httperror.WriteErrorResponse(w, ErrAuthDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	var req postAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidCredentialsFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	var username = req.Username
	var password = req.Password

	u, err := handler.UserService.UserByUsername(username)
	if err == dockm.ErrUserNotFound {
		httperror.WriteErrorResponse(w, ErrInvalidCredentials, http.StatusBadRequest, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	settings, err := handler.SettingsService.Settings()
	if err != nil {
	      httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

  if settings.AuthenticationMethod == dockm.AuthenticationLDAP && u.ID != 1 {
		err = handler.LDAPService.AuthenticateUser(username, password, &settings.LDAPSettings)
		if err != nil {
			httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
			return
		}
	} else {
		err = handler.CryptoService.CompareHashAndData(u.Password, password)
		if err != nil {
			httperror.WriteErrorResponse(w, ErrInvalidCredentials, http.StatusUnprocessableEntity, handler.Logger)
			return
		}
	}

	tokenData := &dockm.TokenData{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
	
  token, err := handler.JWTService.GenerateToken(tokenData)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, &postAuthResponse{JWT: token}, handler.Logger)
}

