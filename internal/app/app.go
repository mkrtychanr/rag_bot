package app

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mkrtychanr/rag_bot/internal/bot"
	"github.com/mkrtychanr/rag_bot/internal/config"
	tgapi "github.com/mkrtychanr/rag_bot/internal/gateway/tg_api"
	"github.com/mkrtychanr/rag_bot/internal/repository"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	smoothoperator "github.com/mkrtychanr/rag_bot/internal/smooth_operator"
	screenchanger "github.com/mkrtychanr/rag_bot/internal/usecase/screen_changer"
	"golang.org/x/sync/errgroup"
)

type app struct {
	config config.Config
	api    *tgbotapi.BotAPI
	conn   *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg config.Config) (*app, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create new bot api: %w", err)
	}

	conn, err := pgxpool.New(ctx, cfg.Postrges.BuildPostgresConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to get new pgx pool: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return &app{
		config: cfg,
		api:    api,
		conn:   conn,
	}, nil
}

func newTree() screen.Screen {
	return nil
}

func (a *app) Run(ctx context.Context) error {
	repo := repository.NewRepository(a.conn)

	tgGateway := tgapi.NewTgAPI(a.api)

	screenChangerUseCase := screenchanger.NewUseCase(tgGateway, repo)

	reciver, err := bot.NewBot(a.api, a.config.Bot)
	if err != nil {
		return fmt.Errorf("failed to create reciver: %w", err)
	}

	operator := smoothoperator.NewOperator(screenChangerUseCase, newTree)

	ch := make(chan tgbotapi.Update, 20)

	defer close(ch)

	var eg errgroup.Group

	eg.Go(func() error {
		if err := reciver.Recive(ctx, ch); err != nil {
			return fmt.Errorf("reciver failed: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		if err := operator.Operate(ctx, ch); err != nil {
			return fmt.Errorf("operator failed: %w", err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("err group finished with error: %w", err)
	}

	return nil
}
