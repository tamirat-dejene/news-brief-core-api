package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
	"github.com/RealEskalate/G6-NewsBrief/internal/handler/http/dto"
)

// SeedAdminUsingUC creates an admin via UserUsecase.Register, then clears topics/subscriptions (nil)
// and ensures the role is admin. Safe to run multiple times.
func SeedAdminUsingUC(
	ctx context.Context,
	userUC contract.IUserUseCase,
	userRepo contract.IUserRepository,
	email, password string,
) error {
	// Check if user exists by email
	existing, err := userRepo.GetUserByEmail(ctx, email)
	if err == nil && existing != nil {
		// Normalize: topics/subscriptions nil and role admin
		existing.Preferences.Topics = nil
		existing.Preferences.SubscribedSources = nil
		existing.Role = entity.UserRoleAdmin
		existing.IsVerified = true
		existing.UpdatedAt = time.Now()

		if _, err := userRepo.UpdateUser(ctx, existing); err != nil {
			return fmt.Errorf("seeder: failed to normalize existing admin: %w", err)
		}
		return nil
	}

	// Create via usecase (adjust fields to match your dto.RegisterRequest)
	req := dto.RegisterRequest{
		Email:    email,
		Password: password,
		Fullname: "admin",
		Username: "admin",
	}

	if _, err := userUC.Register(ctx, req.Username, req.Email, req.Password, req.Fullname); err != nil {
		return fmt.Errorf("seeder: register failed: %w", err)
	}

	// Fetch newly created user and normalize
	u, err := userRepo.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return fmt.Errorf("seeder: failed to fetch newly created admin: %w", err)
	}

	u.Preferences.Topics = nil
	u.Preferences.SubscribedSources = nil
	u.Role = entity.UserRoleAdmin
	u.IsVerified = true
	u.UpdatedAt = time.Now()

	if _, err := userRepo.UpdateUser(ctx, u); err != nil {
		return fmt.Errorf("seeder: failed to finalize admin normalization: %w", err)
	}

	return nil
}
