package ports

import (
	"fmt"
	"net"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	firebaseAuth "firebase.google.com/go/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/query"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
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

	userQueryModel, err := h.app.Queries.User.Handle(r.Context(), cmd.Email)
	if err != nil {
		httperr.InternalError("cannot-get-current-user-in-sign-in-user", err, w, r)
		return
	}

	currentUser := userQueryModelToResponse(userQueryModel)
	render.Respond(w, r, currentUser)
	// return uuid to test
}

func signupUserInFirebaseAuth(w http.ResponseWriter, r *http.Request, user *UserSignUp) string {
	// sign up in firebase auth
	var opts []option.ClientOption
	if file := os.Getenv("SERVICE_ACCOUNT_FILE"); file != "" {
		opts = append(opts, option.WithCredentialsFile(file))
	}

	config := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT")}
	firebaseApp, err := firebase.NewApp(r.Context(), config, opts...)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}

	authClient, err := firebaseApp.Auth(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}

	userToCreate := (&firebaseAuth.UserToCreate{}).
		Email(user.Email).
		Password(user.Password).
		DisplayName(user.DisplayName)

	createdUser, err := authClient.CreateUser(r.Context(), userToCreate)
	if err != nil && firebaseAuth.IsEmailAlreadyExists(err) {
		existingUser, err := authClient.GetUserByEmail(r.Context(), user.Email)
		if err != nil {
			httperr.RespondWithSlugError(errors.Wrap(err, "unable to get created user"), w, r)
		}
		userUuid := existingUser.UID
		return userUuid
	} else {
		if err != nil {
			httperr.RespondWithSlugError(err, w, r)
		}

		err = authClient.SetCustomUserClaims(r.Context(), createdUser.UID, map[string]interface{}{
			"role": user.Role,
		})
		if err != nil {
			httperr.RespondWithSlugError(err, w, r)
		}
		userUuid := createdUser.UID
		return userUuid
	}
}

func (h HttpServer) SignUp(w http.ResponseWriter, r *http.Request) {
	var user *UserSignUp = &UserSignUp{}
	if err := render.Decode(r, user); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	// sign up in firestore
	fmt.Println("signupUserInFirebaseAuth: starting")
	userUuid := signupUserInFirebaseAuth(w, r, user)
	fmt.Println("signupUserInFirebaseAuth: successfuly")

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
		Uuid:          userUuid,
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

func (h HttpServer) UpdateUserInformation(w http.ResponseWriter, r *http.Request) {
	authUser, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var updatedUserInformation *UpdatedUserInformation = &UpdatedUserInformation{}
	if err := render.Decode(r, updatedUserInformation); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := command.UpdateUserInformation{
		Uuid:        authUser.UUID,
		DisplayName: updatedUserInformation.DisplayName,
		Email:       updatedUserInformation.Email,
	}

	if err := h.app.Commands.UpdateUserInformation.Handle(r.Context(), cmd); err != nil {
		httperr.InternalError("cannot-update-user-information", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	authUser, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var updatedUserPassword *UpdatedUserPassword = &UpdatedUserPassword{}
	if err := render.Decode(r, updatedUserPassword); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUserPassword.NewPassword), 12)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := command.UpdateUserPassword{
		Uuid:              authUser.UUID,
		NewHashedPassword: hashedNewPassword,
	}

	if err := h.app.Commands.UpdateUserPassword.Handle(r.Context(), cmd); err != nil {
		httperr.InternalError("cannot-update-user-password", err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func userQueryModelToResponse(model *query.User) *User {
	return &User{
		Balance:     model.Balance,
		DisplayName: model.DisplayName,
		Email:       model.Email,
		Role:        model.Role,
		Uuid:        model.Uuid,
	}
}

func userQueryModelsToResponse(models []*query.User) []*User {
	var users []*User
	for _, u := range models {

		users = append(users, userQueryModelToResponse(u))
	}

	return users
}
