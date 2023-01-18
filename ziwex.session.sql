ALTER TABLE products 
ALTER COLUMN options SET DATA TYPE JSONB[] USING options::jsonb[],
ALTER COLUMN description_key_value SET DATA TYPE JSONB[] USING description_key_value::jsonb[]