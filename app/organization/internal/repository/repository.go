package repository

import (
	"context"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/model/treemodel"
	"github.com/linkbox-group/linkbox-server/organization/internal/acl"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.OrganizationRepositoryItf), new(*OrganizationRepository)), NewOrganizationRepository)

var _ acl.OrganizationRepositoryItf = &OrganizationRepository{}

type OrganizationRepository struct {
	db          *gorm.DB
	treeService *treemodel.TreeService
}

func (r *OrganizationRepository) GetDefaultOrgID(ctx context.Context, code string, userID string) (id string, err error) {
	org := model.Organization{}
	r.db.Where("code = ? AND user_id = ?", code, userID).Select("id").First(&org)
	return org.ID, err
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{
		db:          db,
		treeService: treemodel.NewTreeService(db),
	}
}
