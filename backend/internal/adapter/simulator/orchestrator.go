package simulator

import (
	"context"
	"log"
	"secure-image-service/internal/domain"
)

type BuildOrchestrator interface {
	TriggerBuild(ctx context.Context, event *domain.BuildEvent) error
}

type MockBuildOrchestrator struct{}

func NewMockBuildOrchestrator() BuildOrchestrator {
	return &MockBuildOrchestrator{}
}

func (m *MockBuildOrchestrator) TriggerBuild(ctx context.Context, event *domain.BuildEvent) error {
	log.Printf("SIMULATOR: Triggering build for event ID %d (tenant: %s)", event.ID, event.TenantNamespace)
	// In a real system, this would call a CI/CD system (e.g., Jenkins, GitLab CI).
	// For this MVP, we just log the action. The build completion is simulated
	// by a separate API call to POST /v1/builds/{build_id}/complete.
	return nil
}

