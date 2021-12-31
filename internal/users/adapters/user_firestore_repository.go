package adapters

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserModel struct {
	Uuid           string  `firestore:"Uuid"`
	DisplayName    string  `firestore:"DisplayName"`
	Email          string  `firestore:"Email"`
	HashedPassword string  `firestore:"HashedPassword"`
	Balance        float32 `firestore:"Balance"`
	Role           string  `firestore:"Role"`
	LastIP         string  `firestore:"LastIP"`
}

type FirestoreUserRepository struct {
	firestoreClient *firestore.Client
	userFactory     user.Factory
}

func (f FirestoreUserRepository) GetUser(ctx context.Context, userUuid string) (*query.User, error) {
	userDoc, err := f.UserDocumentRef(userUuid).Get(ctx)

	if err != nil && status.Code(err) != codes.NotFound {
		return &query.User{}, err
	}
	if err != nil && status.Code(err) == codes.NotFound {
		return &query.User{
			Balance: 0,
		}, nil
	}

	var user *UserModel
	err = userDoc.DataTo(&user)
	if err != nil {
		return &query.User{}, err
	}

	return f.userModelToUserQuery(user), nil
}

func (f FirestoreUserRepository) WithdrawBalance(ctx context.Context, userUuid string, amountChange float32) error {
	return f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var user *UserModel

		userDoc, err := tx.Get(f.UserDocumentRef(userUuid))
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			user = &UserModel{
				Balance: 0,
			}
		} else {
			if err := userDoc.DataTo(user); err != nil {
				return err
			}
		}

		user.Balance -= amountChange
		if user.Balance < 0 {
			return errors.New("balance cannot be smaller than 0")
		}

		return tx.Set(userDoc.Ref, user)
	})
}

func (f FirestoreUserRepository) DepositBalance(ctx context.Context, userUuid string, amountChange float32) error {
	return f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var user *UserModel

		userDoc, err := tx.Get(f.UserDocumentRef(userUuid))
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			user = &UserModel{
				Balance: 0,
			}
		} else {
			if err := userDoc.DataTo(user); err != nil {
				return err
			}
		}

		user.Balance += amountChange
		if user.Balance < 0 {
			return errors.New("balance cannot be smaller than 0")
		}

		return tx.Set(userDoc.Ref, user)
	})
}

const lastIPField = "LastIP"

func (f FirestoreUserRepository) UpdateLastIP(ctx context.Context, userUuid string, lastIP string) error {
	updates := []firestore.Update{
		{
			Path:  lastIPField,
			Value: lastIP,
		},
	}

	docRef := f.UserDocumentRef(userUuid)

	_, err := docRef.Update(ctx, updates)
	userNotExist := status.Code(err) == codes.NotFound

	if userNotExist {
		_, err := docRef.Set(ctx, map[string]string{lastIPField: lastIP})
		return err
	}

	return err
}

func (f FirestoreUserRepository) SignIn(ctx context.Context, email string, password string) error {
	query := f.usersCollection().Query.Where("Email", "==", email).limit(1)
	userDocIter := query.Documents(ctx)

	// Only get the first document
	doc, err := userDocIter.Next()
	if err != nil {
		return err
	}

	var userModel *UserModel
	if err := doc.DataTo(userModel); err != nil {
		return err
	}

	userQuery := f.userModelToUserQuery(userModel)

	// compare hashed password from http request and firestore
	if err := bcrypt.CompareHashAndPassword([]byte(userQuery.HashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (f FirestoreUserRepository) SignUp(ctx context.Context, uuid string, displayName string, email string, hashedPassword string, role string, lastIP string) error {

	newBalance := 0.0

	newUserDomain := f.userFactory.NewUser(uuid, displayName, email, hashedPassword, newBalance, role, lastIP)

	newUserModel := f.userDomainToUserModel(newUserDomain)

	newDoc := f.usersCollection().Doc(newUserDomain.GetUuid())
	_, err := newDoc.Create(ctx, newUserModel)
	if err != nil {
		return err
	}

	return nil
}

func (f FirestoreUserRepository) usersCollection() *firestore.CollectionRef {
	return f.firestoreClient.Collection("users")
}

func (f FirestoreUserRepository) UserDocumentRef(userUuid string) *firestore.DocumentRef {
	return f.usersCollection().Doc(userUuid)
}

func (f FirestoreUserRepository) userModelToUserQuery(userModel *UserModel) *query.User {
	return &query.User{
		Uuid:           userModel.Uuid,
		DisplayName:    userModel.DisplayName,
		Email:          userModel.Email,
		HashedPassword: userModel.HashedPassword,
		Balance:        userModel.Balance,
		Role:           userModel.Role,
		LastIP:         userModel.LastIP,
	}
}

func (f FirestoreUserRepository) userDomainToUserModel(user *user.User) *UserModel {
	return &UserModel{
		Uuid:           user.GetUuid(),
		DisplayName:    user.GetDisplayName(),
		Email:          user.GetEmail(),
		HashedPassword: user.GetHashedPassword(),
		Balance:        user.GetBalance(),
		Role:           user.GetRole(),
		LastIP:         user.GetLastIP(),
	}
}
