package main

import (
	"context"
	"fmt"
	"time"

	"github.com/0xcafed00d/joystick"
	"github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        time.RFC3339Nano,
		QuoteEmptyFields:       true,
		PadLevelText:           true,
		DisableLevelTruncation: true,
	})
}

func main() {
	log.Info("starting program...")
	js, err := joystick.Open(0)
	if err != nil {
		panic(err)
	}
	defer js.Close()

	log.Infof("Joystick Name: %s", js.Name())
	log.Infof("   Axis Count: %d", js.AxisCount())
	log.Infof(" Button Count: %d", js.ButtonCount())

	ctx := signals.SetupSignalHandler()

	if err := displayControllerOutput(ctx, js); err != nil {
		log.Error("error reading controller input", err)
	}

}

func displayControllerOutput(ctx context.Context, js joystick.Joystick) error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 40)

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Ch == 'q' {
					return nil
				}
			}
			if ev.Type == termbox.EventResize {
				termbox.Flush()
			}

		case <-ticker.C:

			printAt(1, 0, "-- Press 'q' to Exit --")
			printAt(1, 1, fmt.Sprintf("Joystick Name: %s", js.Name()))
			printAt(1, 2, fmt.Sprintf("   Axis Count: %d", js.AxisCount()))
			printAt(1, 3, fmt.Sprintf(" Button Count: %d", js.ButtonCount()))
			readJoystick(js)
			termbox.Flush()
		}
	}
}

func printAt(x, y int, s string) {
	for _, r := range s {
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}

func readJoystick(js joystick.Joystick) {
	jinfo, err := js.Read()

	if err != nil {
		printAt(1, 5, "Error: "+err.Error())
		return
	}

	printAt(1, 5, "Buttons:")
	for button := 0; button < js.ButtonCount(); button++ {
		if jinfo.Buttons&(1<<uint32(button)) != 0 {
			printAt(10+button, 5, "X")
		} else {
			printAt(10+button, 5, ".")
		}
	}

	printAt(1, 20, fmt.Sprintf("button data: %d", 1<<uint32(jinfo.Buttons)))

	for axis := 0; axis < js.AxisCount(); axis++ {
		printAt(1, axis+7, fmt.Sprintf("Axis %2d Value: %7d", axis, jinfo.AxisData[axis]))
	}
}
