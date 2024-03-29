package services

import (
	"net/http"
	"ziwex/db"
	"ziwex/dtos"
	"ziwex/models"
	"ziwex/types/jsonResponse"
	"ziwex/utils"

	"github.com/jackc/pgx/v5"
)

func AuthAdminRegister(user dtos.AuthAdminRegister) jsonResponse.Response {
	res := jsonResponse.Response{}

	//check exists
	ctx, cancel := utils.GetPgContext()
	defer cancel()

	admin := models.AuthAdmin{}
	err := db.Pg.QueryRow(ctx, `--sql
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
		res.Write(http.StatusConflict, jsonResponse.Json{
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

	ctx2, cancel2 := utils.GetPgContext()
	defer cancel2()
	err2 := db.Pg.QueryRow(ctx2, `--sql
		INSERT INTO admins (firstname, lastname, email, password, username) VALUES ($1, $2, $3,	$4,	$5);
	`, user.Firstname, user.Lastname, user.Email, user.Password, user.Username).Scan()

	if err2 != nil && err2 != pgx.ErrNoRows {
		res.Error(err2)
		return res
	}

	res.Write(http.StatusCreated, jsonResponse.Json{
		"message": "User Created",
	})

	return res
}

func AuthAdminLogin(user dtos.AuthAdminLogin) jsonResponse.Response {
	res := jsonResponse.Response{}

	userDb := models.AuthAdmin{}
	ctx, cancel := utils.GetPgContext()
	defer cancel()
	err := db.Pg.QueryRow(ctx, `--sql
		SELECT password FROM admins WHERE username = $1;
	`, user.Username).Scan(&userDb.Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			res.Write(http.StatusUnauthorized, jsonResponse.Json{
				"message": "unauthorized",
			})
			return res
		}
		res.Error(err)
		return res
	}

	if !utils.CompareHashPassword(userDb.Password, user.Password) {
		res.Write(http.StatusUnauthorized, jsonResponse.Json{
			"message": "unauthorized",
		})
		return res
	}

	token, err2 := utils.JwtSignToken(user.Username, "admin")
	if err2 != nil {
		res.Error(err)
		return res
	}

	res.Write(http.StatusCreated, jsonResponse.Json{
		"message": "Login Successful",
		"token":   token,
	})
	return res
}
