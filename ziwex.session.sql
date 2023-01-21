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
WHERE precp.product_id = 6
GROUP BY p.id;
