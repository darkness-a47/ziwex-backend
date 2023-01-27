package services

import (
	"fmt"
	"net/http"
	"strconv"
	"ziwex/cache"
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

	txCtx, txCancel := utils.GetPgContext()
	defer txCancel()

	txFinalCtx, txFinalCancel := utils.GetPgContext()
	defer txFinalCancel()

	tx, txErr := db.Pg.Begin(txCtx)
	if txErr != nil {
		r.Error(txErr)
		return r
	}

	productCtx, productCancel := utils.GetPgContext()
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
		imageCtx, imageCancel := utils.GetPgContext()
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
			categoriesCtx, categoriesCancel := utils.GetPgContext()
			defer categoriesCancel()
			categoriesErr := tx.QueryRow(categoriesCtx, `--sql
				INSERT INTO product_categories(product_id, category_id)
				VALUES ($1, $2)
			`, productId, cat).Scan()
			txFinalCtx, txFinalCancel := utils.GetPgContext()
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
			relatedProductsCtx, relatedProductsCancel := utils.GetPgContext()
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
			recommendProductsCtx, recommendProductsCancel := utils.GetPgContext()
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
		rollbackCtx, rollbackCancel := utils.GetPgContext()
		defer rollbackCancel()
		_ = tx.Rollback(rollbackCtx)
		return r
	}
	r.Write(http.StatusCreated, jsonResponse.Json{
		"message": "ok",
	})

	go func() {
		cache.InvalidateAll(&cache.Index{
			IndexType:    cache.ProductIndex,
			IndexSubType: cache.ProductSummerySubIndex,
		})
	}()

	return r
}

func GetProductsSummery(d dtos.GetProductsSummery) types.Response {
	r := &jsonResponse.Response{}

	ctx, cancel := utils.GetPgContext()
	defer cancel()

	var rows pgx.Rows
	var err error
	skip := (d.Page - 1) * d.DataPerPage
	if d.CategoryId != nil {
		rows, err = db.Pg.Query(ctx, `--sql
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
		rows, err = db.Pg.Query(ctx, `--sql
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

	jsonR := jsonResponse.Json{
		"message":  "ok",
		"products": products,
	}

	r.Write(http.StatusOK, jsonR)

	go func() {
		cid := 0
		if d.CategoryId != nil {
			cid = *d.CategoryId
		}
		index := fmt.Sprintf("%d,%d,%d", cid, d.Page, d.DataPerPage)
		cache.Store(d.RequestPath, "", &jsonR, &cache.Index{
			IndexType:    cache.ProductIndex,
			IndexSubType: cache.ProductSummerySubIndex,
			Index:        index,
		})
	}()

	return r
}

func GetProductData(d dtos.GetProductData) types.Response {
	r := &jsonResponse.Response{}

	ctx, cancel := utils.GetPgContext()
	defer cancel()

	p := models.Product{}
	err := db.Pg.QueryRow(ctx, `--sql
		SELECT prod.*,
			jsonb_agg(
				DISTINCT jsonb_build_object('id', cat.id, 'title', cat.title)
			) AS categories,
			jsonb_agg(
				DISTINCT jsonb_build_object('file_id', pf.file_id)
			) AS images
		FROM products prod
			LEFT JOIN product_categories pcat ON pcat.product_id = prod.id
			LEFT JOIN categories cat ON cat.id = pcat.category_id
			LEFT JOIN product_images pi ON pi.product_id = prod.id
			LEFT JOIN files pf ON pf.id = pi.image_id
		WHERE url = $1
		GROUP BY prod.id;
	`, d.ProductUrl).Scan(&p.Id, &p.Url, &p.Title, &p.Description, &p.Price, &p.Options,
		&p.DescriptionKeyValue, &p.MainImageIndex, &p.Categories, &p.Images)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.Write(http.StatusNotFound, jsonResponse.Json{
				"message": "product not fount",
			})
			return r
		}
		r.Error(err)
		return r
	}
	//related products

	relCtx, relCancel := utils.GetPgContext()
	defer relCancel()
	relRows, relErr := db.Pg.Query(relCtx, `--sql
		SELECT p.id,
		p.title,
		p.main_image_index,
		p.price,
		p.url,
		jsonb_agg(
			jsonb_build_object(
				'file_id',
				f.file_id
			)
		) AS images
		FROM product_recommend_products precp
			LEFT JOIN products p ON p.id = precp.recommend_product_id
			LEFT JOIN product_images pi ON pi.product_id = p.id
			LEFT JOIN files f ON f.id = pi.image_id
		WHERE precp.product_id = $1
		GROUP BY p.id;
	`, p.Id)

	if relErr != nil {
		r.Error(relErr)
		return r
	}

	for relRows.Next() {
		prel := models.Product{}
		err := relRows.Scan(&prel.Id, &prel.Title, &prel.MainImageIndex, &prel.Price, &prel.Url, &prel.Images)
		if err != nil {
			r.Error(err)
			return r
		}
		p.RelatedProducts = append(p.RelatedProducts, prel)
	}

	//recommend

	recCtx, recCancel := utils.GetPgContext()
	defer recCancel()
	recRows, recErr := db.Pg.Query(recCtx, `--sql
		SELECT p.id,
		p.title,
		p.main_image_index,
		p.price,
		p.url,
		jsonb_agg(
			jsonb_build_object(
				'file_id',
				f.file_id
			)
		) AS images
		FROM product_recommend_products precp
			LEFT JOIN products p ON p.id = precp.recommend_product_id
			LEFT JOIN product_images pi ON pi.product_id = p.id
			LEFT JOIN files f ON f.id = pi.image_id
		WHERE precp.product_id = $1
		GROUP BY p.id;
	`, p.Id)

	if recErr != nil {
		r.Error(recErr)
		return r
	}

	for recRows.Next() {
		prec := models.Product{}
		err := recRows.Scan(&prec.Id, &prec.Title, &prec.MainImageIndex, &prec.Price, &prec.Url, &prec.Images)
		if err != nil {
			r.Error(err)
			return r
		}
		p.RecommendProducts = append(p.RelatedProducts, prec)
	}

	if err != nil {
		r.Error(err)
		return r
	}
	r.Write(http.StatusOK, &p)

	go func() {
		cache.Store(d.RequestPath, "", &p, &cache.Index{
			IndexType:    cache.ProductIndex,
			IndexSubType: cache.ProductDataSubIndex,
			Index:        strconv.Itoa(p.Id),
		})
	}()

	return r
}
