package service

import (
	"context"

	"github.com/hanzy-dev/niskala/apps/api/internal/repository"
)

type MembershipService struct {
	membershipRepository *repository.MembershipRepository
}

func NewMembershipService(membershipRepository *repository.MembershipRepository) *MembershipService {
	return &MembershipService{
		membershipRepository: membershipRepository,
	}
}

func (s *MembershipService) IsAdmin(ctx context.Context, userID string) (bool, error) {
	role, exists, err := s.membershipRepository.GetRoleByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	return role == "admin", nil
}
