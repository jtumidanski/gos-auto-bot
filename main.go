package main

import (
	"context"
	"gos-auto-bot/discord"
	"gos-auto-bot/logger"
	"gos-auto-bot/orchestrator"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l := logger.CreateLogger("gos-auto-bot").WithField("thread", "main")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	instructions := make(chan orchestrator.Instruction, 1)
	go orchestrator.NewOrchestrator(l)(ctx, wg, instructions)
	go discord.NewWorker(l)(ctx, wg, instructions)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
