package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/frain-dev/convoy/datastore"
	pager "github.com/gobeam/mongo-go-pagination"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	innerDB *mongo.Database
	client  *mongo.Collection
	store   datastore.Store
}

func NewUserRepo(db *mongo.Database, store datastore.Store) datastore.UserRepository {
	return &userRepo{
		innerDB: db,
		client:  db.Collection(UserCollection),
		store:   store,
	}
}

func (u *userRepo) CreateUser(ctx context.Context, user *datastore.User) error {
	user.ID = primitive.NewObjectID()
	user.ResetPasswordToken = uuid.NewString()

	err := u.store.Save(ctx, user, nil)

	if mongo.IsDuplicateKeyError(err) {
		return datastore.ErrDuplicateEmail
	}

	return err
}

func (u *userRepo) FindUserByEmail(ctx context.Context, email string) (*datastore.User, error) {
	user := &datastore.User{}

	filter := bson.M{"email": email}

	err := u.store.FindOne(ctx, filter, nil, user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, datastore.ErrUserNotFound
	}

	return user, nil
}

func (u *userRepo) FindUserByID(ctx context.Context, id string) (*datastore.User, error) {
	user := &datastore.User{}

	err := u.store.FindByID(ctx, id, nil, user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, datastore.ErrUserNotFound
	}

	return user, nil
}

func (u *userRepo) LoadUsersPaged(ctx context.Context, pageable datastore.Pageable) ([]datastore.User, datastore.PaginationData, error) {
	var users []datastore.User

	filter := bson.M{"document_status": datastore.ActiveDocumentStatus}

	paginatedData, err := pager.New(u.client).Context(ctx).Limit(int64(pageable.PerPage)).Page(int64(pageable.Page)).Sort("created_at", -1).Filter(filter).Decode(&users).Find()
	if err != nil {
		return users, datastore.PaginationData{}, err
	}

	if users == nil {
		users = make([]datastore.User, 0)
	}

	return users, datastore.PaginationData(paginatedData.Pagination), nil
}

func (u *userRepo) UpdateUser(ctx context.Context, user *datastore.User) error {
	update := bson.D{
		primitive.E{Key: "first_name", Value: user.FirstName},
		primitive.E{Key: "last_name", Value: user.LastName},
		primitive.E{Key: "email", Value: user.Email},
		primitive.E{Key: "password", Value: user.Password},
		primitive.E{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())},
		primitive.E{Key: "reset_password_token", Value: user.ResetPasswordToken},
		primitive.E{Key: "reset_password_expires_at", Value: user.ResetPasswordExpiresAt},
	}

	err := u.store.UpdateByID(ctx, user.UID, update)

	if mongo.IsDuplicateKeyError(err) {
		return datastore.ErrDuplicateEmail
	}
	
	return err
}

func (u *userRepo) FindUserByToken(ctx context.Context, token string) (*datastore.User, error) {
	user := &datastore.User{}

	filter := bson.M{"reset_password_token": token}

	err := u.store.FindOne(ctx, filter, nil, user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, datastore.ErrUserNotFound
	}

	return user, nil
}
