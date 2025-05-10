package repository

const (
	// Select Table
	TblEtalse = "etalases"

	// Selector Column
	SelectorTblEtalaseAll = "etalases.*"

	// List Query
	QueryDetailEtalase      = "SELECT * FROM products WHERE id IN (SELECT id_product FROM etalase_product WHERE id_etalase = ?) AND deleted_at IS NULL AND product_name LIKE ? ORDER BY product_name ASC"
	QuerySelectEtalaseList  = "(SELECT COUNT(*) FROM etalase_product ep JOIN products p ON ep.id_product = p.id WHERE ep.id_etalase = etalases.id AND p.deleted_at IS NULL) as total_product"
	QueryWhereSelectEtalase = "username = ? AND deleted_at IS NULL AND etalase_name LIKE ?"
	QueryWhereDeletedAtNull = "deleted_at IS NULL AND id = ?"

	// Selector Column
	ColumnSelectEtalaseID = "id_etalase = ?"
	ColumnSelectID        = "id = ? AND deleted_at IS NULL"
)
