package delivery

import (
	"context"
	"github.com/google/uuid"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/model/treemodel"
	"github.com/linkbox-group/linkbox-server/organization/pkg/log"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/cError"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateOrganization implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) CreateOrganization(ctx context.Context, req *organization.CreateOrganizationRequest) (resp *organization.CreateOrganizationResponse, err error) {
	oid := uuid.New().String()
	oModel := model.Organization{
		BaseModel: model.BaseModel{
			ID: oid,
		},
		UserID: req.UserId,
		TreeModel: treemodel.TreeModel{
			Code:       oid,
			Name:       req.Name,
			ParentCode: *req.ParentCode,
		},
		Description: req.Description,
		SortOrder:   int(*req.SortOrder),
	}

	err = s.service.CreateOrganizationService(ctx, &oModel)
	if err != nil {
		logrus.Errorln(err)
		return &organization.CreateOrganizationResponse{
			Result: &organization.CreateOrganizationResponse_Error{
				Error: &cError.Error{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, err
	}
	return &organization.CreateOrganizationResponse{
		Result: &organization.CreateOrganizationResponse_Organization{
			Organization: oModel.Convert(),
		},
	}, nil
}

// GetOrganization implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) GetOrganization(ctx context.Context, req *organization.GetOrganizationRequest) (resp *organization.GetOrganizationResponse, err error) {
	oModel, err := s.service.GetOrganizationService(ctx, req.Id, req.UserId)
	if err != nil {
		resp = &organization.GetOrganizationResponse{
			Result: &organization.GetOrganizationResponse_Error{
				Error: &cError.Error{
					Code:    500,
					Message: err.Error(),
				},
			},
		}
		return
	}

	resp = &organization.GetOrganizationResponse{
		Result: &organization.GetOrganizationResponse_Organization{
			Organization: oModel.Convert(),
		},
	}
	return
}

// UpdateOrganization implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) UpdateOrganization(ctx context.Context, req *organization.UpdateOrganizationRequest) (resp *organization.UpdateOrganizationResponse, err error) {
	oModel := model.Organization{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
	}
	log.Log().Debug("id is" + oModel.ID)

	if req.Name != nil {
		oModel.Name = *req.Name
	}
	if req.Description != nil {
		oModel.Description = req.Description
	}
	if req.SortOrder != nil {
		oModel.SortOrder = int(*req.SortOrder)
	}

	updatedModel, err := s.service.UpdateOrganizationService(ctx, &oModel)
	if err != nil {
		resp = &organization.UpdateOrganizationResponse{
			Result: &organization.UpdateOrganizationResponse_Error{
				Error: &cError.Error{
					Code:    500,
					Message: err.Error(),
				},
			},
		}
		return
	}

	resp = &organization.UpdateOrganizationResponse{
		Result: &organization.UpdateOrganizationResponse_Organization{
			Organization: &organization.Organization{
				Id:          updatedModel.ID,
				Code:        updatedModel.Code,
				ParentCode:  updatedModel.ParentCode,
				ParentCodes: updatedModel.ParentCodes,
				TreeLeaf:    updatedModel.TreeLeaf,
				TreeLevel:   int32(updatedModel.TreeLevel),
				TreeNames:   updatedModel.TreeNames,
				Name:        updatedModel.Name,
				UserId:      updatedModel.UserID,
				Description: *updatedModel.Description,
				SortOrder:   int32(updatedModel.SortOrder),
				ItemsCount:  uint32(updatedModel.ItemsCount),
				CreatedAt:   timestamppb.New(updatedModel.CreatedAt),
				UpdatedAt:   timestamppb.New(updatedModel.UpdatedAt),
			},
		},
	}
	return
}

// DeleteOrganization implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) DeleteOrganization(ctx context.Context, req *organization.DeleteOrganizationRequest) (resp *organization.DeleteOrganizationResponse, err error) {
	err = s.service.DeleteOrganizationService(ctx, req.Id, req.UserId, req.Cascade)
	if err != nil {
		resp = &organization.DeleteOrganizationResponse{
			Result: &organization.DeleteOrganizationResponse_Error{
				Error: &cError.Error{
					Code:    500,
					Message: err.Error(),
				},
			},
		}
		return
	}

	resp = &organization.DeleteOrganizationResponse{
		Result: &organization.DeleteOrganizationResponse_Success{
			Success: true,
		},
	}
	return
}

// GetUserOrganizations implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) GetUserOrganizations(ctx context.Context, req *organization.GetUserOrganizationsRequest) (resp *organization.GetUserOrganizationsResponse, err error) {
	organizations, err := s.service.GetUserOrganizationsService(ctx, req.UserId)
	logrus.Infoln(organizations)
	if err != nil {
		resp = &organization.GetUserOrganizationsResponse{
			Result: &organization.GetUserOrganizationsResponse_Error{
				Error: &cError.Error{
					Code:    500,
					Message: err.Error(),
				},
			},
		}
		return
	}

	orgList := make([]*organization.Organization, 0, len(organizations))
	for _, org := range organizations {
		orgList = append(orgList, &organization.Organization{
			Id:          org.ID,
			Code:        org.Code,
			ParentCode:  org.ParentCode,
			ParentCodes: org.ParentCodes,
			TreeLeaf:    org.TreeLeaf,
			TreeLevel:   int32(org.TreeLevel),
			TreeNames:   org.TreeNames,
			Name:        org.Name,
			UserId:      org.UserID,
			Description: *org.Description,

			SortOrder:  int32(org.SortOrder),
			ItemsCount: uint32(org.ItemsCount),
			CreatedAt:  timestamppb.New(org.CreatedAt),
			UpdatedAt:  timestamppb.New(org.UpdatedAt),
		})
	}

	resp = &organization.GetUserOrganizationsResponse{
		Result: &organization.GetUserOrganizationsResponse_Organizations{
			Organizations: &organization.OrganizationsData{
				Organizations: orgList,
			},
		},
	}
	return
}

// MoveOrganization implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) MoveOrganization(ctx context.Context, req *organization.MoveOrganizationRequest) (resp *organization.MoveOrganizationResponse, err error) {
	err = s.service.MoveOrganizationService(ctx, req.Id, req.UserId, req.NewParentCode)
	if err != nil {
		resp = &organization.MoveOrganizationResponse{
			Result: &organization.MoveOrganizationResponse_Error{
				Error: &cError.Error{
					Code:    500,
					Message: err.Error(),
				},
			},
		}
		return
	}

	resp = &organization.MoveOrganizationResponse{
		Result: &organization.MoveOrganizationResponse_Success{
			Success: true,
		},
	}
	return
}

// AddItemsToOrganization implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) AddItemsToOrganization(ctx context.Context, req *organization.AddItemsToOrganizationRequest) (resp *organization.AddItemsToOrganizationResponse, err error) {
	err = s.service.AddItemsToOrganizationService(ctx, req.OrganizationId, req.UserId, req.ItemIds)
	if err != nil {
		return &organization.AddItemsToOrganizationResponse{
			Result: &organization.AddItemsToOrganizationResponse_Error{
				Error: &cError.Error{
					Message: err.Error(),
				},
			}}, err
	}

	return &organization.AddItemsToOrganizationResponse{
		Result: &organization.AddItemsToOrganizationResponse_Success{
			Success: true}}, nil
}

// RemoveItemsFromOrganization implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) RemoveItemsFromOrganization(ctx context.Context, req *organization.RemoveItemsFromOrganizationRequest) (resp *organization.RemoveItemsFromOrganizationResponse, err error) {
	err = s.service.RemoveItemsFromOrganizationService(ctx, req.OrganizationId, req.UserId, req.ItemIds)
	if err != nil {
		return &organization.RemoveItemsFromOrganizationResponse{
			Result: &organization.RemoveItemsFromOrganizationResponse_Error{
				Error: &cError.Error{
					Message: err.Error(),
				},
			}}, err
	}

	return &organization.RemoveItemsFromOrganizationResponse{
		Result: &organization.RemoveItemsFromOrganizationResponse_Success{
			Success: true}}, nil
}

// ReorderOrganizationItems implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) ReorderOrganizationItems(ctx context.Context, req *organization.ReorderOrganizationItemsRequest) (resp *organization.ReorderOrganizationItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// ReorderOrganizations implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) ReorderOrganizations(ctx context.Context, req *organization.ReorderOrganizationsRequest) (resp *organization.ReorderOrganizationsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetOrganizationActivity implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) GetOrganizationActivity(ctx context.Context, req *organization.GetOrganizationActivityRequest) (resp *organization.GetOrganizationActivityResponse, err error) {
	// TODO: Your code here...
	return
}

// GetOrganizationTree implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) GetOrganizationTree(ctx context.Context, req *organization.GetOrganizationTreeRequest) (resp *organization.GetOrganizationTreeResponse, err error) {
	// TODO: Your code here...
	return
}

// GetOrganizationChildren implements the OrganizationDelivery interface.
func (s *OrganizationDelivery) GetOrganizationChildren(ctx context.Context, req *organization.GetOrganizationChildrenRequest) (resp *organization.GetOrganizationChildrenResponse, err error) {
	// TODO: Your code here...
	return
}
