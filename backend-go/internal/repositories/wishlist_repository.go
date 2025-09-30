package repositories
import (
	"database/sql"
	"time"
	"ecommerce-backend/internal/models"
	"github.com/google/uuid"
)
type WishlistRepository struct {
	db *sql.DB
}
func NewWishlistRepository(db *sql.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}
func (r *WishlistRepository) AddToWishlist(userID, productID string) error {
	query := `
		INSERT INTO wishlist_items (id, user_id, product_id, created_at)
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, uuid.New().String(), userID, productID, time.Now())
	return err
}
func (r *WishlistRepository) RemoveFromWishlist(userID, productID string) error {
	query := `DELETE FROM wishlist_items WHERE user_id = $1 AND product_id = $2`
	_, err := r.db.Exec(query, userID, productID)
	return err
}
func (r *WishlistRepository) IsInWishlist(userID, productID string) (bool, error) {
	query := `SELECT COUNT(*) FROM wishlist_items WHERE user_id = $1 AND product_id = $2`
	var count int
	err := r.db.QueryRow(query, userID, productID).Scan(&count)
	return count > 0, err
}
func (r *WishlistRepository) GetUserWishlistItems(userID string, limit, offset int) ([]models.WishlistItemWithProduct, error) {
	query := `
		SELECT wi.id, wi.user_id, wi.product_id, wi.created_at,
		       p.id, p.name, p.description, p.price, p.images, p.category_id,
		       p.stock, p.featured, p.created_at, p.updated_at
		FROM wishlist_items wi
		JOIN products p ON wi.product_id = p.id
		WHERE wi.user_id = $1
		ORDER BY wi.created_at DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.WishlistItemWithProduct
	for rows.Next() {
		var item models.WishlistItemWithProduct
		var product models.Product
		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.CreatedAt,
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.Images, &product.CategoryID, &product.Stock,
			&product.Featured, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		item.Product = &models.ProductWithRating{
			Product: product,
		}
		items = append(items, item)
	}
	return items, nil
}
func (r *WishlistRepository) CountUserWishlistItems(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM wishlist_items WHERE user_id = $1`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}
func (r *WishlistRepository) ClearUserWishlist(userID string) error {
	query := `DELETE FROM wishlist_items WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}