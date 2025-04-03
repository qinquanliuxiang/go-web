package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"qqlx/base/conf"
	"qqlx/base/constant"
	"qqlx/base/server"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

const (
	Version = "1.0.0"
)

// Application is the main struct of the application
type Application struct {
	name    string
	version string
	servers []server.ServerInterface
	signals []os.Signal
}

func NewApplication(e *gin.Engine) *Application {
	return newApp(
		withName(conf.GetProjectName()),
		withVersion(constant.ServerVersion),
		withServer(server.NewServer(e)),
	)
}

// Option application support option
type Option func(application *Application)

// newApp creates a new Application
func newApp(ops ...Option) *Application {
	app := &Application{}
	for _, op := range ops {
		op(app)
	}

	// default accept signals
	if len(app.signals) == 0 {
		app.signals = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	}
	return app
}

// withName application add name
func withName(name string) func(application *Application) {
	return func(application *Application) {
		application.name = name
	}
}

// withVersion application add version
func withVersion(version string) func(application *Application) {
	return func(application *Application) {
		application.version = version
	}
}

// withServer application add server
func withServer(servers ...server.ServerInterface) func(application *Application) {
	return func(application *Application) {
		application.servers = servers
	}
}

// WithSignals application add listen signals
func WithSignals(signals []os.Signal) func(application *Application) {
	return func(application *Application) {
		application.signals = signals
	}
}

// Run application run
func (app *Application) Run(ctx context.Context) error {
	if len(app.servers) == 0 {
		return nil
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, app.signals...)
	errCh := make(chan error, 1)

	for _, s := range app.servers {
		go func(srv server.ServerInterface) {
			if err := srv.Start(); err != nil {
				zap.S().Errorf("failed to start server, err: %s", err)
				errCh <- err
			}
		}(s)
	}

	select {
	case err := <-errCh:
		_ = app.Stop()
		return err
	case <-ctx.Done():
		return app.Stop()
	case <-quit:
		return app.Stop()
	}
}

// Stop application stop
func (app *Application) Stop() error {
	wg := sync.WaitGroup{}
	for _, s := range app.servers {
		wg.Add(1)
		go func(srv server.ServerInterface) {
			defer wg.Done()
			if err := srv.Shutdown(); err != nil {
				zap.S().Errorf("failed to stop server, err: %s", err)
			}
		}(s)
	}
	wg.Wait()
	return nil
}
