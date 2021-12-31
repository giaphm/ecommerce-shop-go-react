package ports

import (
	"net/http"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/command"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

func (h HttpServer) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	authUser, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	// host, _, err := net.SplitHostPort(r.RemoteAddr)
	// if err == nil {
	// 	err = h.db.UpdateLastIP(r.Context(), authUser.UUID, host)
	// 	if err != nil {
	// 		httperr.InternalError("internal-server-error", err, w, r)
	// 		return
	// 	}
	// }

	cmd := command.UpdateLastIP{
		UserUuid:   authUser.UUID,
		RemoteAddr: r.RemoteAddr,
	}

	if err := h.app.Commands.UpdateLastIp.Handle(r.Context(), cmd); err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}

	// user, err := h.db.GetUser(r.Context(), authUser.UUID)
	// if err != nil {
	// 	httperr.InternalError("cannot-get-user", err, w, r)
	// 	return
	// }

	user, err := h.app.Queries.CurrentUser.Handle(
		r.Context(),
		authUser.UUID,
	)
	if err != nil {
		httperr.InternalError("cannot-get-current-user", err, w, r)
	}

	userResponse := User{
		Email:       authUser.Email,
		DisplayName: authUser.DisplayName,
		Balance:     user.Balance,
		Role:        authUser.Role,
	}

	render.Respond(w, r, userResponse)
}

func (h HttpServer) SignInUser(w http.ResponseWriter, r *http.Request) {

	var user *User
	if err := render.Decode(r, user); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, 12)
	// if err != nil {
	// 	httperr.RespondWithSlugError(err, w, r)
	// 	return
	// }

	cmd := command.SignIn{
		Email:   user.Email,
		Pasword: user.Password,
	}

	if err := h.app.Commands.SignIn.Handle(r.Context(), cmd); err != nil {
		httperr.InternalError("cannot-sign-in-user", err, w, r)
		return
	}

	w.Header().Set("content-location", "users/sign-in-user/"+cmd.Uuid)
	w.WriteHeader(http.StatusCreated)
}

func (h HttpServer) SignUpUser(w http.ResponseWriter, r *http.Request) {

	var user *User
	if err := render.Decode(r, user); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(user.Password, 12)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := command.SignUpUser{
		Uuid:          uuid.New().String(),
		DisplayName:   user.DisplayName,
		Email:         user.Email,
		HashedPasword: hashedPassword,
		Role:          user.Role,
		LastIP:        "unknown", //
	}

	if err := h.app.Commands.SignUp.Handle(r.Context(), cmd); err != nil {
		httperr.InternalError("cannot-sign-up-user", err, w, r)
		return
	}

	w.Header().Set("content-location", "users/sign-up-user/"+cmd.Uuid)
	w.WriteHeader(http.StatusCreated)
}
