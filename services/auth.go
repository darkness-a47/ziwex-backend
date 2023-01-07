package services

import (
	"net/http"
	"ziwex/db"
	"ziwex/dtos"
	"ziwex/models"
	"ziwex/types"
	"ziwex/utils"

	"github.com/jackc/pgx/v5"
)

func AuthAdminRegister(user dtos.AuthAdminRegister) types.Response {
	res := types.Response{}

	//check exists
	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()

	admin := models.AuthAdmin{}
	err := db.Poll.QueryRow(ctx, `--sql
		SELECT email, username FROM admins WHERE email = $1 OR username = $2;
	`, user.Email, user.Username).Scan(&admin.Email, &admin.Username)
	if err != nil && err != pgx.ErrNoRows {
		res.Error(err)
		return res
	}

	if admin.Username != "" || admin.Email != "" {
		e := make([]string, 0)
		if admin.Email == user.Email {
			e = append(e, "Email exists")
		}
		if admin.Username == user.Username {
			e = append(e, "Username exists")
		}
		res.Write(http.StatusConflict, types.JsonR{
			"message": e,
		})
		return res
	}

	//insert
	var err3 error
	user.Password, err3 = utils.GenerateHashedPassword(user.Password)
	if err3 != nil {
		res.Error(err3)
		return res
	}

	ctx2, cancel2 := utils.GetDatabaseContext()
	defer cancel2()
	err2 := db.Poll.QueryRow(ctx2, `--sql
		INSERT INTO admins (firstname, lastname, email, password, username) VALUES ($1, $2, $3,	$4,	$5);
	`, user.Firstname, user.Lastname, user.Email, user.Password, user.Username).Scan()

	if err2 != nil && err2 != pgx.ErrNoRows {
		res.Error(err2)
		return res
	}

	res.Write(http.StatusCreated, types.JsonR{
		"message": "User Created",
	})

	return res
}

func AuthAdminLogin(user dtos.AuthAdminLogin) types.Response {
	res := types.Response{}

	userDb := models.AuthAdmin{}
	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()
	err := db.Poll.QueryRow(ctx, `--sql
		SELECT password FROM admins WHERE username = $1;
	`, user.Username).Scan(&userDb.Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			res.Write(http.StatusUnauthorized, types.JsonR{
				"message": "unauthorized",
			})
			return res
		}
		res.Error(err)
		return res
	}

	if !utils.CompareHashPassword(userDb.Password, user.Password) {
		res.Write(http.StatusUnauthorized, types.JsonR{
			"message": "unauthorized",
		})
		return res
	}

	token, err2 := utils.JwtSignToken(user.Username, "admin")
	if err2 != nil {
		res.Error(err)
		return res
	}

	res.Write(http.StatusCreated, types.JsonR{
		"message": "Login Successful",
		"token":   token,
	})
	return res
}
