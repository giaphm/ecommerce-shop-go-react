package adapters

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
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

func NewFirestoreUserRepository(firestoreClient *firestore.Client, userFactory user.Factory) *FirestoreUserRepository {
	if firestoreClient == nil {
		panic("missing firestoreClient")
	}
	// if productFactory.IsZero() {
	// 	panic("missing productFactory")
	// }

	return &FirestoreUserRepository{firestoreClient, userFactory}
}

func (f FirestoreUserRepository) GetUser(ctx context.Context, userUuid string) (*query.User, error) {
	userDoc, err := f.UserDocumentRef(userUuid).Get(ctx)

	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	if err != nil && status.Code(err) == codes.NotFound {
		return nil, err
	}

	var user *UserModel = &UserModel{}
	err = userDoc.DataTo(&user)
	if err != nil {
		return nil, err
	}

	return f.userModelToUserQuery(user), nil
}

func (f FirestoreUserRepository) GetUsers(ctx context.Context) ([]*query.User, error) {
	userIters, err := f.userDocuments(ctx)
	if err != nil {
		return nil, err
	}
	defer userIters.Stop()

	var users []*UserModel
	var user *UserModel

	for {
		user = &UserModel{}
		userIter, err := userIters.Next()
		fmt.Println("userIter", userIter)
		if err == iterator.Done {
			fmt.Println("Done userIters")
			break
		}
		if err != nil {
			return nil, err
		}
		fmt.Println("userIter.Data()", userIter.Data())

		if err := userIter.DataTo(user); err != nil {
			fmt.Println("user", user)
			return nil, err
		}

		users = append(users, user)
	}

	fmt.Println("users", users)

	return f.userModelsToUserQueries(users), nil
}

func (f FirestoreUserRepository) WithdrawBalance(ctx context.Context, userUuid string, amountChange float32) error {
	return f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var user *UserModel = &UserModel{}

		userDoc, err := tx.Get(f.UserDocumentRef(userUuid))
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			// user = &UserModel{
			// 	Balance: 0,
			// }
			// user.Balance = 0
			return err
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
		var user *UserModel = &UserModel{}

		userSnapshot, err := tx.Get(f.UserDocumentRef(userUuid))
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			return err
		}

		if err := userSnapshot.DataTo(user); err != nil {
			return err
		}

		user.Balance += amountChange
		if user.Balance < 0 {
			return errors.New("balance cannot be smaller than 0")
		}

		return tx.Set(userSnapshot.Ref, user)
	})
}

// const lastIPField = "LastIP"

func (f FirestoreUserRepository) UpdateLastIP(ctx context.Context, userUuid string, lastIP string) error {
	if lastIP == "" {
		return errors.New("last ip is empty")
	}

	return f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		var user *UserModel = &UserModel{}

		userDocRef := f.UserDocumentRef(userUuid)

		userSnapshot, err := tx.Get(userDocRef)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			return err
		}

		if err := userSnapshot.DataTo(user); err != nil {
			return err
		}

		user.LastIP = lastIP

		return tx.Set(userSnapshot.Ref, user)
	})
	return nil
}

func (f FirestoreUserRepository) SignIn(ctx context.Context, email string, password string) error {
	query := f.usersCollection().Query.Where("Email", "==", email).Limit(1)
	userDocIter := query.Documents(ctx)

	// Only get the first document
	doc, err := userDocIter.Next()
	if err != nil {
		return err
	}

	var userModel *UserModel = &UserModel{}
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

func (f FirestoreUserRepository) SignUp(ctx context.Context, uuid string, displayName string, email string, hashedPassword []byte, role string, lastIP string) error {

	var newBalance float32 = 0.0

	newUserDomain, err := f.userFactory.NewUser(uuid, displayName, email, hashedPassword, newBalance, role, lastIP)
	if err != nil {
		return err
	}

	newUserModel := f.userDomainToUserModel(newUserDomain)

	newDoc := f.usersCollection().Doc(newUserDomain.GetUuid())
	_, err = newDoc.Create(ctx, newUserModel)
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

func (f FirestoreUserRepository) userDocuments(ctx context.Context) (*firestore.DocumentIterator, error) {
	return f.usersCollection().Documents(ctx), nil //.GetAll()
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

func (f FirestoreUserRepository) userModelsToUserQueries(um []*UserModel) []*query.User {

	var users []*query.User
	var user *query.User

	for _, u := range um {
		user = f.userModelToUserQuery(u)
		users = append(users, user)
	}

	return users
}

func (f FirestoreUserRepository) userDomainToUserModel(user user.IUser) *UserModel {
	return &UserModel{
		Uuid:           user.GetUuid(),
		DisplayName:    user.GetDisplayName(),
		Email:          user.GetEmail(),
		HashedPassword: string(user.GetHashedPassword()),
		Balance:        user.GetBalance(),
		Role:           user.GetRole(),
		LastIP:         user.GetLastIP(),
	}
}
