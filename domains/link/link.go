package link

import "gorm.io/gorm"

// Link represents a saved URL
type Link struct {
	gorm.Model
	Path string
	URL  string
}

// Service provides methods for managing links
type Service struct {
	DB *gorm.DB
}

// NewService creates a new LinkService
func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

// Create saves a new link to the database
func (s *Service) Create(link *Link) error {
	return s.DB.Create(link).Error
}

// Delete removes a link from the database
func (s *Service) Delete(link *Link) error {
	return s.DB.Delete(link).Error
}

// GetByPath retrieves a link by its path
func (s *Service) GetByPath(path string) (*Link, error) {
	var link Link
	err := s.DB.Where("path = ?", path).First(&link).Error
	return &link, err
}

// GetAll retrieves all links from the database
func (s *Service) GetAll() ([]Link, error) {
	var links []Link
	err := s.DB.Find(&links).Error
	return links, err
}
