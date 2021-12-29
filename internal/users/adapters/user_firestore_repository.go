package main

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserModel struct {
	uuid           string
	DisplayName    string
	Email          string
	HashedPassword string
	Balance        int
	Role           string
	LastIP         string
}

type FirestoreUserRepository struct {
	firestoreClient *firestore.Client
	userFactory     user.Factory
}

func (d FirestoreUserRepository) usersCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("users")
}

func (d FirestoreUserRepository) UserDocumentRef(userID string) *firestore.DocumentRef {
	return d.usersCollection().Doc(userID)
}

func (d FirestoreUserRepository) GetUser(ctx context.Context, userID string) (UserModel, error) {
	doc, err := d.UserDocumentRef(userID).Get(ctx)

	if err != nil && status.Code(err) != codes.NotFound {
		return UserModel{}, err
	}
	if err != nil && status.Code(err) == codes.NotFound {
		return UserModel{
			Balance: 0,
		}, nil
	}

	var user UserModel
	err = doc.DataTo(&user)
	if err != nil {
		return UserModel{}, err
	}

	return user, nil
}

func (d FirestoreUserRepository) WithdrawBalance(ctx context.Context, userID string, amountChange int) error {
	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var user UserModel

		userDoc, err := tx.Get(d.UserDocumentRef(userID))
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			user = UserModel{
				Balance: 0,
			}
		} else {
			if err := userDoc.DataTo(&user); err != nil {
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

func (d FirestoreUserRepository) DepositBalance(ctx context.Context, userID string, amountChange int) error {
	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var user UserModel

		userDoc, err := tx.Get(d.UserDocumentRef(userID))
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}
		if err != nil && status.Code(err) == codes.NotFound {
			user = UserModel{
				Balance: 0,
			}
		} else {
			if err := userDoc.DataTo(&user); err != nil {
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

func (d FirestoreUserRepository) UpdateLastIP(ctx context.Context, userID string, lastIP string) error {
	updates := []firestore.Update{
		{
			Path:  lastIPField,
			Value: lastIP,
		},
	}

	docRef := d.firestoreClient.UserDocumentRef(userID)

	_, err := docRef.Update(ctx, updates)
	userNotExist := status.Code(err) == codes.NotFound

	if userNotExist {
		_, err := docRef.Set(ctx, map[string]string{lastIPField: lastIP})
		return err
	}

	return err
}
