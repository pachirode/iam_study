package posixsignal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pachirode/iam_study/pkg/shutdown"
)

const Name = "PosixSignalManager"

type PosixSignalManager struct {
	signals []os.Signal
}

func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{
		signals: sig,
	}
}

func (posixSignalManagerP *PosixSignalManager) GetName() string {
	return Name
}

func (posixSignalManagerP *PosixSignalManager) Start(gsi shutdown.GracefulShutdownInterface) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, posixSignalManagerP.signals...)

		<-c
		gsi.StartShutdown(posixSignalManagerP)
	}()

	return nil
}

func (posixSignalManagerP *PosixSignalManager) ShutdownStart() error {
	return nil
}

func (posixSignalManagerP *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)

	return nil
}
