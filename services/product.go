package services

import (
	"fmt"
	"net/http"
	"ziwex/db"
	"ziwex/dtos"
	"ziwex/types"
	"ziwex/types/jsonResponse"
	"ziwex/utils"

	"github.com/jackc/pgx/v5"
)

func CreateProduct(d dtos.CreateProduct) types.Response {
	r := &jsonResponse.Response{}

	txCtx, txCancel := utils.GetDatabaseContext()
	defer txCancel()

	txFinalCtx, txFinalCancel := utils.GetDatabaseContext()
	defer txFinalCancel()

	tx, txErr := db.Poll.Begin(txCtx)
	if txErr != nil {
		r.Error(txErr)
		return r
	}

	productCtx, productCancel := utils.GetDatabaseContext()
	defer productCancel()

	var productId int
	productErr := tx.QueryRow(productCtx, `--sql
		INSERT INTO products (url, title, description, price, options, description_key_value, main_image_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;
	`, d.Url, d.Title, d.Description, d.Price, d.Options, d.DescriptionKeyValue, d.MainImageIndex).Scan(&productId)
	if productErr != nil {
		_ = tx.Rollback(txFinalCtx)
		r.Error(productErr)
		return r
	}

	for _, img := range d.Images {
		imageCtx, imageCancel := utils.GetDatabaseContext()
		defer imageCancel()

		imageErr := tx.QueryRow(imageCtx, `
			INSERT INTO product_images (product_id, image_id) VALUES ($1, $2);
		`, productId, img).Scan()

		if imageErr != nil && imageErr != pgx.ErrNoRows {
			_ = tx.Rollback(txFinalCtx)
			r.Write(http.StatusBadRequest, jsonResponse.Json{
				"message":  "image not found",
				"category": img,
			})
			return r
		}
	}

	if len(d.Categories) > 0 {
		for _, cat := range d.Categories {
			categoriesCtx, categoriesCancel := utils.GetDatabaseContext()
			defer categoriesCancel()
			categoriesErr := tx.QueryRow(categoriesCtx, `--sql
				INSERT INTO product_categories(product_id, category_id)
				VALUES ($1, $2)
			`, productId, cat).Scan()

			if categoriesErr != nil && categoriesErr != pgx.ErrNoRows {
				_ = tx.Rollback(txFinalCtx)
				r.Write(http.StatusBadRequest, jsonResponse.Json{
					"message":  "category not found",
					"category": cat,
				})
				return r
			}
		}
	}

	if len(d.RelatedProducts) > 0 {
		for _, rp := range d.RelatedProducts {
			relatedProductsCtx, relatedProductsCancel := utils.GetDatabaseContext()
			defer relatedProductsCancel()
			relatedProductsErr := tx.QueryRow(relatedProductsCtx, `--sql
				INSERT INTO product_related_products(product_id, related_product_id)
				VALUES ($1, $2)
			`, productId, rp).Scan()

			fmt.Println(relatedProductsErr)
			if relatedProductsErr != nil && relatedProductsErr != pgx.ErrNoRows {
				_ = tx.Rollback(txFinalCtx)
				r.Write(http.StatusBadRequest, jsonResponse.Json{
					"message": "related product not found",
					"product": rp,
				})
				return r
			}
		}
	}

	if len(d.RecommendProducts) > 0 {
		for _, rp := range d.RecommendProducts {
			recommendProductsCtx, recommendProductsCancel := utils.GetDatabaseContext()
			defer recommendProductsCancel()
			recommendProductsErr := tx.QueryRow(recommendProductsCtx, `--sql
				INSERT INTO product_recommend_products(product_id, recommend_product_id)
				VALUES ($1, $2)
			`, productId, rp).Scan()

			if recommendProductsErr != nil && recommendProductsErr != pgx.ErrNoRows {
				_ = tx.Rollback(txFinalCtx)
				r.Write(http.StatusBadRequest, jsonResponse.Json{
					"message": "recommend product not found",
					"product": rp,
				})
				return r
			}
		}
	}

	txFinalErr := tx.Commit(txFinalCtx)
	if txFinalErr != nil {
		rollbackCtx, rollbackCancel := utils.GetDatabaseContext()
		defer rollbackCancel()
		_ = tx.Rollback(rollbackCtx)
		r.Error(txFinalErr)
		return r
	}
	r.Write(http.StatusCreated, jsonResponse.Json{
		"message": "ok",
	})

	return r
}
