package services

import (
	"net/http"
	"ziwex/db"
	"ziwex/dtos"
	"ziwex/models"
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
			txFinalCtx, txFinalCancel := utils.GetDatabaseContext()
			defer txFinalCancel()
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

func GetProductsSummery(d dtos.GetProductsSummery) types.Response {
	r := &jsonResponse.Response{}

	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()

	var rows pgx.Rows
	var err error
	skip := (d.Page - 1) * d.DataPerPage
	if d.CategoryId != nil {
		rows, err = db.Poll.Query(ctx, `--sql
			SELECT prod.*, jsonb_agg(DISTINCT jsonb_build_object('file_id' , f.file_id)) AS images FROM (
				SELECT p.id, p.url, p.title, p.price, p.main_image_index, COUNT(*) OVER() AS total_rows
				FROM products p
				INNER JOIN product_categories pc ON pc.product_id = p.id
				WHERE pc.category_id = $1
				ORDER BY id DESC OFFSET $2 LIMIT $3
			) prod
			LEFT JOIN product_images pi ON pi.product_id = prod.id
			LEFT JOIN files f ON f.id = pi.image_id
			GROUP BY prod.id, prod.url, prod.title, prod.price, prod.main_image_index, prod.total_rows;
		`, *d.CategoryId, skip, d.DataPerPage)
	} else {
		rows, err = db.Poll.Query(ctx, `--sql
			SELECT prod.*, jsonb_agg(DISTINCT jsonb_build_object('file_id' , f.file_id)) AS images FROM (
				SELECT id, url, title, price, main_image_index, COUNT(*) OVER() AS total_rows
				FROM products ORDER BY id DESC OFFSET $1 LIMIT $2
			) prod
			LEFT JOIN product_images pi ON pi.product_id = prod.id
			LEFT JOIN files f ON f.id = pi.image_id
			GROUP BY prod.id, prod.url, prod.title, prod.price, prod.main_image_index, prod.total_rows;
		`, skip, d.DataPerPage)
	}
	if err != nil {
		r.Error(err)
		return r
	}

	var totalRows int
	products := make([]models.Product, 0)
	for rows.Next() {
		p := models.Product{}
		rowErr := rows.Scan(&p.Id, &p.Url, &p.Title, &p.Price, &p.MainImageIndex, &totalRows, &p.Images)
		if rowErr != nil {
			r.Error(rowErr)
			return r
		}
		products = append(products, p)
	}

	r.Write(http.StatusOK, jsonResponse.Json{
		"message":  "ok",
		"products": products,
	})
	return r
}

func GetProductData(d dtos.GetProductData) types.Response {
	r := &jsonResponse.Response{}

	ctx, cancel := utils.GetDatabaseContext()
	defer cancel()

	p := models.Product{}
	err := db.Poll.QueryRow(ctx, `--sql
		SELECT
			prod.*,
			jsonb_agg(DISTINCT jsonb_build_object(
				'id', cat.id,
				'title', cat.title
			)) AS categories,
			jsonb_agg(DISTINCT jsonb_build_object(
				'file_id', pf.file_id
			)) AS images,
			jsonb_agg(DISTINCT jsonb_build_object(
				'id', prec.id,
				'images', prec.images
			)) AS product_recommend_products,
			jsonb_agg(DISTINCT jsonb_build_object(
				'id', prel.id,
				'images', prel.images
			)) AS product_related_products
		FROM products prod

		INNER JOIN product_categories pcat ON pcat.product_id = prod.id
		INNER JOIN categories cat ON cat.id = pcat.category_id

		INNER JOIN product_images pi ON pi.product_id = prod.id
		INNER JOIN files pf ON pf.id = pi.image_id

		INNER JOIN product_recommend_products precp ON precp.product_id = prod.id
		CROSS JOIN LATERAL (
			SELECT
				p.id,
				jsonb_agg(jsonb_build_object('file_id', f.file_id)) AS images
			FROM products p
			INNER JOIN product_images i ON i.product_id = p.id
			INNER JOIN files f ON f.id = i.image_id
			WHERE p.id = precp.recommend_product_id
			GROUP BY p.id
		) AS prec

		INNER JOIN product_related_products prelp ON prelp.product_id = prod.id
		CROSS JOIN LATERAL (
			SELECT
				p.id,
				jsonb_agg(jsonb_build_object('file_id', f.file_id)) AS images
			FROM products p
			INNER JOIN product_images i ON i.product_id = p.id
			INNER JOIN files f ON f.id = i.image_id
			WHERE p.id = prelp.related_product_id
			GROUP BY p.id
		) AS prel

		WHERE url = $1
		GROUP BY prod.id;
	`, d.ProductUrl).Scan(&p.Id, &p.Url, &p.Title, &p.Description, &p.Price, &p.Options,
		&p.DescriptionKeyValue, &p.MainImageIndex, &p.Categories, &p.Images, &p.RecommendProducts, &p.RelatedProducts)

	if err != nil {
		r.Error(err)
		return r
	}
	r.Write(http.StatusOK, &p)

	return r
}
