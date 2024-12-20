package storage

const (
	insertProductsQuery = "INSERT INTO prices (id, name, category, price, create_date) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING"
	getAllProductsQuery = "SELECT id, name, category, price, create_date FROM prices ORDER BY id"

	baseFilteredQuery = "SELECT id, name, category, price, create_date FROM prices"

	getStatsQuery = `
		SELECT 
			COUNT(*) as total_items,
			COUNT(DISTINCT category) as total_categories,
			COALESCE(SUM(price), 0) as total_price
		FROM prices;
	`

	fullDayEndSuffix = " 23:59:59"
)
