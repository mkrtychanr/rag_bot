package actioncontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	adddocumentscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/documents/add_document_screen"
	requestscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/request_screen"
	selectorscreen "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
)

var (
	ErrUnknownScreenType = errors.New("unknown screen type")
	ErrCorruptedData     = errors.New("corrupted operator data")
)

type controllerFunc = func(ctx context.Context, screen screen.Screen, data model.OperatorData) (screen.Screen, error)

type actionController struct {
	controllers map[screen.ScreenType]controllerFunc
}

func NewActionController() *actionController {
	return &actionController{
		controllers: map[screen.ScreenType]controllerFunc{
			screen.DefaultScreen:                 defaultScreenController,
			screen.RequestScreen:                 requestScreenController,
			screen.AddDocumentScreen:             addDocumentScreenController,
			screen.ChangeDocumentsScreen:         changeDocumentScreenController,
			screen.PostE2ESelectorScreen:         postE2EselectorScreenController,
			screen.AddDocumentIntoGroupScreen:    addDocumentIntoGroupController,
			screen.DeleteDocumentFromGroupScreen: addDocumentIntoGroupController,
			screen.DeleteDocumentScreen:          deleteDocumentScreenController,
			screen.GroupAccessScreen:             defaultScreenController,
			screen.PickGroupScreen:               pickGroupScreenController,
			screen.GroupDocumentsScreen:          defaultScreenController,
			screen.GroupUsersScreen:              defaultScreenController,
		},
	}
}

func (ac *actionController) GetScreenController(screenType screen.ScreenType) (controllerFunc, error) {
	f, ok := ac.controllers[screenType]
	if !ok {
		return nil, ErrUnknownScreenType
	}

	return f, nil
}

func defaultScreenController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	var v model.MenuOption
	if err := json.Unmarshal(payload.CallbackData, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	m := map[string]any{
		"option":  v,
		"user_id": payload.UserID,
	}

	newSc, err := sc.Next(ctx, m)
	if err != nil {
		return nil, fmt.Errorf("failed to get next screen: %w", err)
	}

	return newSc, nil
}

func requestScreenController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	if payload.Text != nil {
		m := map[string]any{
			"perform": requestscreen.PerformModel{
				UserID:  payload.UserID,
				Request: *payload.Text,
			},
		}
		sc, err := sc.Perform(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("failed to perform request screen: %w", err)
		}

		return sc, nil
	}
	if payload.CallbackData != nil {
		return defaultScreenController(ctx, sc, payload)
	}

	return nil, ErrCorruptedData
}

func addDocumentScreenController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	if payload.Text != nil && payload.DocumentID != nil {
		m := map[string]any{
			"perform": adddocumentscreen.PerformModel{
				Name:       *payload.Text,
				UserID:     payload.UserID,
				DocumentID: *payload.DocumentID,
			},
		}
		sc, err := sc.Perform(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("failed to perform add document screen: %w", err)
		}

		return sc, nil
	}

	if payload.CallbackData != nil {
		return defaultScreenController(ctx, sc, payload)
	}

	return nil, ErrCorruptedData
}

type lightWeightOption struct {
	Option int `json:"option"`
}

func changeDocumentScreenController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	if payload.CallbackData != nil {
		var v lightWeightOption
		if err := json.Unmarshal(payload.CallbackData, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		var op selectorscreen.SelectorOption

		if v.Option < 0 {
			op.Option = v.Option
		} else {
			if err := json.Unmarshal(payload.CallbackData, &op); err != nil {
				return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
			}
		}

		m := map[string]any{
			"selector_option": op,
			"user_id":         payload.UserID,
			"paper_id":        op.Payload.ID,
		}

		newSc, err := sc.Perform(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("failed to perform end to end selector: %w", err)
		}

		return newSc, nil
	}

	return nil, ErrCorruptedData
}

func postE2EselectorScreenController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	if payload.CallbackData != nil {
		var v model.MenuOption
		if err := json.Unmarshal(payload.CallbackData, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		m := sc.ExtractPayload()
		m["option"] = v
		m["user_id"] = payload.UserID

		newSc, err := sc.Next(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("failed to get next screen: %w", err)
		}

		return newSc, nil
	}

	return nil, ErrCorruptedData
}

func addDocumentIntoGroupController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	if payload.CallbackData != nil {
		var v lightWeightOption
		if err := json.Unmarshal(payload.CallbackData, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		var op selectorscreen.SelectorOption

		if v.Option < 0 {
			op.Option = v.Option
		} else {
			if err := json.Unmarshal(payload.CallbackData, &op); err != nil {
				return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
			}
		}

		m := sc.ExtractPayload()
		m["selector_option"] = op
		m["group_id"] = op.Payload.ID

		sc, err := sc.Perform(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("failed to perform screen: %w", err)
		}

		return sc, nil
	}

	return nil, ErrCorruptedData
}

func deleteDocumentScreenController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	if payload.CallbackData != nil {
		var v lightWeightOption
		if err := json.Unmarshal(payload.CallbackData, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		var op selectorscreen.SelectorOption

		if v.Option < 0 {
			op.Option = v.Option
		} else {
			if err := json.Unmarshal(payload.CallbackData, &op); err != nil {
				return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
			}
		}

		m := sc.ExtractPayload()
		m["selector_option"] = op
		m["paper_id"] = op.Payload.ID

		sc, err := sc.Perform(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("failed to perform screen: %w", err)
		}

		return sc, nil
	}

	return nil, ErrCorruptedData
}

func pickGroupScreenController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	if payload.CallbackData != nil {
		var v lightWeightOption
		if err := json.Unmarshal(payload.CallbackData, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		var op selectorscreen.SelectorOption

		if v.Option < 0 {
			op.Option = v.Option
		} else {
			if err := json.Unmarshal(payload.CallbackData, &op); err != nil {
				return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
			}
		}

		m := sc.ExtractPayload()
		m["selector_option"] = op
		m["group_id"] = op.Payload.ID

		sc, err := sc.Perform(ctx, m)
		if err != nil {
			return nil, fmt.Errorf("failed to perform screen: %w", err)
		}

		return sc, nil
	}

	return nil, ErrCorruptedData
}
