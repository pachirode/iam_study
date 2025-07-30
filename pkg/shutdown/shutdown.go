package shutdown

import (
	"sync"
)

type ShutdownCallback interface {
	OnShutdown(string) error
}

type ShutdownManager interface {
	GetName() string
	Start(gsi GracefulShutdownInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

type GracefulShutdownInterface interface {
	StartShutdown(shutdownManager ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallback ShutdownCallback)
}

type ErrorHandler interface {
	OnError(err error)
}

type ErrorFunc func(err error)

func (f ErrorFunc) OnError(err error) {
	f(err)
}

type ShutdownFunc func(string) error

func (f ShutdownFunc) OnShutdown(shutdownManager string) error {
	return f(shutdownManager)
}

type GracefulShutdown struct {
	callbacks    []ShutdownCallback
	managers     []ShutdownManager
	errorHandler ErrorHandler
}

func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]ShutdownCallback, 0, 10),
		managers:  make([]ShutdownManager, 0, 3),
	}
}

func (gracefulShutdownP *GracefulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallback) {
	gracefulShutdownP.callbacks = append(gracefulShutdownP.callbacks, shutdownCallback)
}

func (gracefulShudownP *GracefulShutdown) AddShutdownManager(shutdownManager ShutdownManager) {
	gracefulShudownP.managers = append(gracefulShudownP.managers, shutdownManager)
}

func (gracefulShudownP *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gracefulShudownP.errorHandler = errorHandler
}

func (gracefulShutdownP *GracefulShutdown) StartShutdown(shutdownManager ShutdownManager) {
	gracefulShutdownP.ReportError(shutdownManager.ShutdownStart())

	var wg sync.WaitGroup
	for _, shutdownCallback := range gracefulShutdownP.callbacks {
		wg.Add(1)
		go func(shutdownCallback ShutdownCallback) {
			defer wg.Done()

			gracefulShutdownP.ReportError(shutdownCallback.OnShutdown(shutdownManager.GetName()))
		}(shutdownCallback)
	}

	wg.Wait()

	gracefulShutdownP.ReportError(shutdownManager.ShutdownFinish())
}

func (gracefulShutdownP *GracefulShutdown) ReportError(err error) {
	if err != nil && gracefulShutdownP.errorHandler != nil {
		gracefulShutdownP.errorHandler.OnError(err)
	}
}
