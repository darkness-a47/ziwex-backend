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

WHERE url = 'necklace2'
GROUP BY prod.id;