package services
import (
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repositories"
)
type WishlistService struct {
	wishlistRepo *repositories.WishlistRepository
}
func NewWishlistService(wishlistRepo *repositories.WishlistRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
	}
}
func (s *WishlistService) GetUserWishlist(userID string, page, limit int) ([]models.WishlistItemWithProduct, int, error) {
	offset := (page - 1) * limit
	items, err := s.wishlistRepo.GetUserWishlistItems(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.wishlistRepo.CountUserWishlistItems(userID)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
func (s *WishlistService) AddToWishlist(userID, productID string) error {
	exists, err := s.wishlistRepo.IsInWishlist(userID, productID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Already in wishlist, no error
	}
	return s.wishlistRepo.AddToWishlist(userID, productID)
}
func (s *WishlistService) RemoveFromWishlist(userID, productID string) error {
	return s.wishlistRepo.RemoveFromWishlist(userID, productID)
}
func (s *WishlistService) IsInWishlist(userID, productID string) (bool, error) {
	return s.wishlistRepo.IsInWishlist(userID, productID)
}
func (s *WishlistService) ClearWishlist(userID string) error {
	return s.wishlistRepo.ClearUserWishlist(userID)
}
