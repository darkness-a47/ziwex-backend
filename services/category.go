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

func CreateCategory(d dtos.CreateCategory) types.Response {
	res := types.Response{}

	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()

	cat := models.Category{}
	err := db.Poll.QueryRow(ctx, `--sql
		INSERT INTO categories (title, image_url, description, parent_category_id, tags) VALUES ($1, $2, $3, $4, $5) RETURNING id;
		`, d.Title, d.ImageUrl, d.Description, d.ParentCategoryId, d.Tags).Scan(&cat.Id)
	if err != nil {
		res.Error(err)
		return res
	}

	res.Write(http.StatusCreated, types.JsonR{
		"message":     "category created",
		"category_id": cat.Id,
	})

	return res
}

func GetCategories(d dtos.GetCategories) types.Response {
	res := types.Response{}

	categories := make([]models.Category, 0)

	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()

	offset := (d.Page - 1) * d.DataPerPage
	var rows pgx.Rows
	var err error
	if d.ParentCategoryId != nil {
		rows, err = db.Poll.Query(ctx, `--sql
			SELECT id, title, image_url, description, parent_category_id, tags, COUNT(*) OVER() AS total_count 
			FROM categories WHERE parent_category_id = $1 OFFSET $2 LIMIT $3
		`, d.ParentCategoryId, offset, d.DataPerPage)
	} else {
		rows, err = db.Poll.Query(ctx, `--sql
		SELECT id, title, image_url, description, parent_category_id, tags,  COUNT(*) OVER() AS total_count 
		FROM categories WHERE parent_category_id IS NULL OFFSET $1 LIMIT $2
		`, offset, d.DataPerPage)
	}
	if err != nil {
		res.Error(err)
		return res
	}

	var totalRows int
	for rows.Next() {
		c := models.Category{}
		err := rows.Scan(&c.Id, &c.Title, &c.ImageUrl, &c.Description, &c.ParentCategoryId, &c.Tags, &totalRows)
		if err != nil {
			res.Error(err)
			return res
		}
		categories = append(categories, c)
	}

	res.Write(http.StatusOK, types.JsonR{
		"message":    "ok",
		"categories": categories,
		"total_rows": totalRows,
	})
	return res
}
