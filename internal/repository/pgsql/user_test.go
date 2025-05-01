package pgsql_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Fadil-Tao/gopher-pay/internal/model"
	"github.com/Fadil-Tao/gopher-pay/internal/repository/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestShouldRegister(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error initializing sqlmock : %s", err)
	}
	defer db.Close() 

	userRepo := &pgsql.UserRepo{db}
	
	mockUser := &model.User{
		Email: "dimas@gmail.com",
		Name: "dimas",
		Password: "randompasswordhuzzah",
		Salt:  "sosaltyhuhu",
		Phone: "0876364736",
	}
	
	mock.ExpectBegin()
	
	mock.ExpectPrepare("insert into users\\(email,name,password,salt,phone\\) values").
	ExpectQuery().
	WithArgs(mockUser.Email,mockUser.Name,mockUser.Password,mockUser.Salt,mockUser.Phone).
	WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	createFnCalled := false
	mockCreatefn := func (userId int) error{
		createFnCalled = true
		assert.Equal(t, 1, userId)
		return nil
	}
	err = userRepo.Register(context.Background(),*mockUser, mockCreatefn)

	assert.NoError(t,err)
	assert.True(t,createFnCalled, "create wallet fn should be called")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestShouldGetUserByEmail(t *testing.T){
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected ")
	}
	defer db.Close()

	userRepo := pgsql.NewUserRepo(db)

	rows := sqlmock.NewRows([]string{"id", "email", "name","password","salt","phone", "updated_at", "created_at"}).
		AddRow(1, "dimas@gmail.com", "dimas", "passworDimas", "saltDimas","09876543121",time.Now(), time.Now())
	
	query := "select id,email,name,password,salt,phone,updated_at,created_at from users where email = $1"

	email := "dimas@gmail.com"

	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(email).
		WillReturnRows(rows)

	user, err :=  userRepo.GetByEmail(context.Background(),email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}