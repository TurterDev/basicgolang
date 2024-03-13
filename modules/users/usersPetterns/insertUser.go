package usersPetterns

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TurterDev/basicgolang/modules/users"
	"github.com/jmoiron/sqlx"
)

type IInertUser interface {
	Customer() (IInertUser, error)
	Admin() (IInertUser, error)
	Result() (*users.UserPassport, error)
}

type userReq struct {
	id  string
	req *users.UserRegisterReq
	db  *sqlx.DB
}

// สร้าง private struc ซึ่งการทำแบบนี้ที่ *userReq จะทำให้สามารถเข้าถึงเข้า .id , .req, .db ที่ userReq โดยตรงได้เลย
type customer struct {
	*userReq
}

type admin struct {
	*userReq
}

// โรงงารใหญ่ InsertUser
func InsertUser(db *sqlx.DB, req *users.UserRegisterReq, isAdmin bool) IInertUser {
	if isAdmin {
		return newAdmin(db, req)
	}
	return newCustomer(db, req)
}

// การสร้าง newCustomer กับ newAdmin คือการสร้างโรงงานเล็กจะสร้างมีโรงการใหญ่ควบคุมคือ InsertUser
func newCustomer(db *sqlx.DB, req *users.UserRegisterReq) IInertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}

}
func newAdmin(db *sqlx.DB, req *users.UserRegisterReq) IInertUser {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func (f *userReq) Customer() (IInertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2, $3, 1)
	RETURNING "id";`

	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}

	return f, nil
}

func (f *userReq) Admin() (IInertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2, $3, 2)
	RETURNING "id";`

	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}
	
	return f, nil
}

func (f *userReq) Result() (*users.UserPassport, error) {
	query := `
	SELECT
		json_build_object(
			'user',"t",
			'token', NULL
		)
	FROM (
		SELECT
		 	"u"."id",
		 	"u"."email",
		 	"u"."username",
		 	"u"."role_id",
		FROM "user" "u"
		WHERE "u"."id" = $1
	) AS "t"`

	data := make([]byte, 0)
	if err := f.db.Get(&data, query, f.id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	user := new(users.UserPassport)
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}

	return user, nil
}
