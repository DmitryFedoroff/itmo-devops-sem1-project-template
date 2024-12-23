package storage

const (
	createTempTableQuery = `
		CREATE TEMP TABLE temp_products (
			id INT, 
			name VARCHAR(255), 
			category VARCHAR(255), 
			price DECIMAL(10, 2), 
			create_date TIMESTAMP WITH TIME ZONE
		);
	`

	insertIntoTempTableQuery = `
		INSERT INTO temp_products (id, name, category, price, create_date) 
		VALUES ($1, $2, $3, $4, $5);
	`

	upsertAndStatsQuery = `
		WITH duplicates AS (
			SELECT COUNT(*) AS duplicate_count
			FROM temp_products t
			JOIN prices p
			ON t.name = p.name 
			   AND t.category = p.category
			   AND t.price = p.price
			   AND t.create_date = p.create_date
		),
		inserted AS (
			INSERT INTO prices (id, name, category, price, create_date)
			SELECT id, name, category, price, create_date
			FROM temp_products
			ON CONFLICT DO NOTHING
			RETURNING *
		),
		stats AS (
			SELECT
				(SELECT COUNT(*) FROM inserted) AS total_items,
				(SELECT COUNT(DISTINCT category) FROM prices) AS total_categories,
				(SELECT COALESCE(SUM(price), 0) FROM prices) AS total_price
		)
		SELECT
			(SELECT duplicate_count FROM duplicates),
			(SELECT total_items FROM stats),
			(SELECT total_categories FROM stats),
			(SELECT total_price FROM stats);
	`

	getAllProductsQuery = "SELECT id, name, category, price, create_date FROM prices ORDER BY id"

	baseFilteredQuery = "SELECT id, name, category, price, create_date FROM prices"

	fullDayEndSuffix = " 23:59:59"
)
