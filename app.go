package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
	svc *MailService
}

// NewApp creates a new App application struct
func NewApp(svc *MailService) *App {
	return &App{svc: svc}
}

// startup is called when the app starts. The context is saved so we can call the
// runtime methods, and the mail service is started (restore session + sync).
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.svc != nil {
		a.svc.Start(ctx)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
