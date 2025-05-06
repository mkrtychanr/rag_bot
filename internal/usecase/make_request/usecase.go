package makerequest

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/utils"
)

type useCase struct {
	rag        rag
	repository repository
}

func NewUseCase(rag rag, repo repository) *useCase {
	return &useCase{
		rag:        rag,
		repository: repo,
	}
}

func (u *useCase) MakeRequest(ctx context.Context, request string, tgID int64) (string, error) {
	userGroups, err := u.repository.GetUserGroups(ctx, tgID)
	if err != nil {
		return "", fmt.Errorf("failed to get user groups: %w", err)
	}

	ownedGroups, err := u.repository.GetUserGroupsOwnership(ctx, tgID)
	if err != nil {
		return "", fmt.Errorf("failed to get user groups ownership: %w", err)
	}

	for _, group := range ownedGroups {
		userGroups = append(userGroups, group.ID)
	}

	papers := make([]int64, 0, len(userGroups))

	for _, group := range userGroups {
		groupPapers, err := u.repository.GetGroupPapers(ctx, group)
		if err != nil {
			return "", fmt.Errorf("failed to get group papers: %w", err)
		}

		papers = append(papers, groupPapers...)
	}

	uniquePapers := utils.GetUniqueSlice(papers)

	response, err := u.rag.GetLLMResponse(ctx, request, uniquePapers)
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}

	return response, nil
}
