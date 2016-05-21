package persistence

import (
	"github.com/zlowram/godan/model"
)

type PersistenceManager interface {
	SaveBanner(banner model.Banner)
	QueryBanners(f model.Filters) ([]model.Banner, error)
	Close()
}
