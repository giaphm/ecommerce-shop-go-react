package ports

import (
	"fmt"
	"net"
	"net/http"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/query"
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

	if err := h.app.Commands.UpdateLastIP.Handle(r.Context(), cmd); err != nil {
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

	var userResponse *User

	userResponse = &User{
		Uuid:        authUser.UUID,
		Email:       authUser.Email,
		DisplayName: authUser.DisplayName,
		Balance:     user.Balance,
		Role:        authUser.Role,
	}

	render.Respond(w, r, userResponse)
}

func (h HttpServer) GetUsers(w http.ResponseWriter, r *http.Request) {
	_, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	userQueryModels, err := h.app.Queries.Users.Handle(r.Context())
	fmt.Println("userQueryModels", userQueryModels)
	if err != nil {
		fmt.Println("Errors userQueryModels, err := h.app.Queries.Users.Handle(r.Context())")
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	users := userQueryModelsToResponse(userQueryModels)
	render.Respond(w, r, users)
}

func userQueryModelsToResponse(models []*query.User) []*User {
	var users []*User
	for _, u := range models {

		users = append(users, &User{
			Balance:     u.Balance,
			DisplayName: u.DisplayName,
			Email:       u.Email,
			Role:        u.Role,
			Uuid:        u.Uuid,
		})
	}

	return users
}

func (h HttpServer) SignIn(w http.ResponseWriter, r *http.Request) {

	var user *UserSignIn = &UserSignIn{}
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
		Email:    user.Email,
		Password: user.Password,
	}

	if err := h.app.Commands.SignIn.Handle(r.Context(), cmd); err != nil {
		httperr.InternalError("cannot-sign-in-user", err, w, r)
		return
	}
	// return uuid to test
}

func (h HttpServer) SignUp(w http.ResponseWriter, r *http.Request) {

	var user *UserSignUp = &UserSignUp{}
	if err := render.Decode(r, user); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := command.SignUp{
		Uuid:          uuid.New().String(),
		DisplayName:   user.DisplayName,
		Email:         user.Email,
		HashedPasword: hashedPassword,
		Role:          user.Role,
		LastIP:        host, //
	}

	if err := h.app.Commands.SignUp.Handle(r.Context(), cmd); err != nil {
		httperr.InternalError("cannot-sign-up-user", err, w, r)
		return
	}

	w.Header().Set("content-location", "users/sign-up-user/"+cmd.Uuid)
	w.WriteHeader(http.StatusCreated)
}
