package step

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
	"gos-auto-bot/coordinate"
	"gos-auto-bot/identity"
	"image"
	"image/color"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	LaunchUrl     = "https://trigger.macrodroid.com/%s/launch"
	LaunchMDUrl   = "https://trigger.macrodroid.com/%s/launchmd"
	BackUrl       = "https://trigger.macrodroid.com/%s/back"
	ClickUrl      = "https://trigger.macrodroid.com/%s/click?clickx=%d&clicky=%d"
	GestureUrl    = "https://trigger.macrodroid.com/%s/gesture?startX=%d&startY=%d&endX=%d&endY=%d&duration=%d"
	ScreenshotUrl = "https://trigger.macrodroid.com/%s/ss"
)

type Operator func(bot identity.Identity)

type Step struct {
	StepFunc Operator
	Delay    int
}

func GestureStep(startX int, startY int, endX int, endY int, duration int, delay int) Step {
	return Step{
		StepFunc: gesture(startX, startY, endX, endY, duration),
		Delay:    delay,
	}
}

func gesture(startX int, startY int, endX int, endY int, duration int) Operator {
	return func(bot identity.Identity) {
		simpleGet(fmt.Sprintf(GestureUrl, bot.DeviceId(), startX, startY, endX, endY, duration))
	}
}

func BackStep(delay int) Step {
	return Step{
		StepFunc: back,
		Delay:    delay,
	}
}

func back(bot identity.Identity) {
	simpleGet(fmt.Sprintf(BackUrl, bot.DeviceId()))
}

func ClickStep(c coordinate.Model, delay int) Step {
	return Step{
		StepFunc: click(c.X(), c.Y()),
		Delay:    delay,
	}
}

func PrintStep(l logrus.FieldLogger) func(message string) Step {
	return func(message string) Step {
		return Step{
			StepFunc: func(bot identity.Identity) {
				l.Infof(message)
			},
			Delay: 0,
		}
	}
}

func click(x int, y int) Operator {
	return func(bot identity.Identity) {
		simpleGet(fmt.Sprintf(ClickUrl, bot.DeviceId(), x, y))
	}
}

func LaunchStep() Step {
	return Step{
		StepFunc: launchGame,
		Delay:    0,
	}
}

func launchGame(bot identity.Identity) {
	simpleGet(fmt.Sprintf(LaunchUrl, bot.DeviceId()))
}

func launchMD(bot identity.Identity) {
	simpleGet(fmt.Sprintf(LaunchMDUrl, bot.DeviceId()))
}

func getImageBotImageDir(l logrus.FieldLogger) func(name string) string {
	botImagePath, ok := os.LookupEnv("BOT_IMAGE_PATH")
	if !ok {
		l.Fatalf("Unable to lookup BOT_IMAGE_PATH")
	}
	return func(name string) string {
		return fmt.Sprintf(botImagePath, name)
	}
}

func VerifyPixelStep(l logrus.FieldLogger) func(c coordinate.Model, r uint8, g uint8, b uint8, recoveryFunc Operator, reset Operator) Step {
	return func(c coordinate.Model, r uint8, g uint8, b uint8, recoveryFunc Operator, reset Operator) Step {
		return Step{
			StepFunc: func(bot identity.Identity) {
				dir := getImageBotImageDir(l)(bot.Name())
				err := clearImageDirectory(l)(bot.Name(), dir)
				if err != nil {
					return
				}

				repeatFuncForImage(l)(bot, dir, imageVerification(bot, c.X(), c.Y(), r, g, b, recoveryFunc), reset)
			},
			Delay: 0,
		}
	}
}

func repeatFuncForImage(l logrus.FieldLogger) func(bot identity.Identity, dir string, repeatFunc repeatableImageFunc, reset Operator) {
	return func(bot identity.Identity, dir string, repeatFunc repeatableImageFunc, reset Operator) {
		attempt := 0
		maxAttempt := 60
		verified := false
		for !verified {
			takeScreenshot(bot.DeviceId())

			tries := 0
			maxTries := 15
			for !verified && tries < maxTries {
				tries++

				files, _ := ioutil.ReadDir(dir)
				var newestFile string

				if len(files) > 0 {
					var newestTime int64 = 0
					for _, f := range files {
						fi, err := os.Stat(dir + f.Name())
						if err != nil {
							l.WithError(err).Warnf("Unable to read statistics for file %s.", f.Name())
							break
						}
						currTime := fi.ModTime().Unix()
						if currTime > newestTime {
							newestTime = currTime
							newestFile = f.Name()
						}
					}

					f, err := os.Open(dir + newestFile)
					if err != nil {
						l.WithError(err).Warnf("Unable to open image %s.", newestFile)
						_ = f.Close()
						continue
					}

					i, _, err := image.Decode(f)
					if err != nil {
						l.WithError(err).Warnf("Unable to decode image %s.", newestFile)
						_ = f.Close()
						_ = os.Remove(dir + newestFile)
						continue
					}

					verified = !repeatFunc(i)

					err = f.Close()
					if err != nil {
						l.WithError(err).Errorf("Unable to close image %s.", newestFile)
						continue
					}
					err = os.Remove(dir + newestFile)
					if err != nil {
						l.WithError(err).Errorf("Unable to delete image %s.", newestFile)
						continue
					}
				} else {
					time.Sleep(time.Duration(250) * time.Millisecond)
				}
			}
			attempt++

			if attempt >= maxAttempt {
				if reset != nil {
					reset(bot)
				}
				attempt = 0
			}
		}
	}
}

type repeatableImageFunc func(i image.Image) bool

func imageVerification(bot identity.Identity, x int, y int, r uint8, g uint8, b uint8, recoveryFunc Operator) repeatableImageFunc {
	return func(i image.Image) bool {
		nrgb := i.At(x, y).(color.NRGBA)
		if r == nrgb.R && g == nrgb.G && b == nrgb.B {
			return false
		} else {
			recoveryFunc(bot)
			return true
		}
	}
}

func takeScreenshot(deviceId string) {
	simpleGet(fmt.Sprintf(ScreenshotUrl, deviceId))
}

func clearImageDirectory(l logrus.FieldLogger) func(botName string, dir string) error {
	return func(botName string, dir string) error {
		err := os.RemoveAll(dir)
		if err != nil {
			l.WithError(err).Errorf("Unable to clear the image directory of %s.", botName)
			return err
		}
		err = os.Mkdir(dir, os.ModeDir)
		if err != nil {
			l.WithError(err).Errorf("Unable to make the image directory of %s.", botName)
			return err
		}
		return nil
	}
}

func ClickColorStep(l logrus.FieldLogger) func(start coordinate.Model, end coordinate.Model, cmFunc ColorMatcher, delay int) Step {
	return func(start coordinate.Model, end coordinate.Model, cmFunc ColorMatcher, delay int) Step {
		return Step{
			StepFunc: func(bot identity.Identity) {
				dir := getImageBotImageDir(l)(bot.Name())
				err := clearImageDirectory(l)(bot.Name(), dir)
				if err != nil {
					return
				}

				repeatFuncForImage(l)(bot, dir, findColorAndClick(l)(bot, start.X(), start.Y(), end.X(), end.Y(), cmFunc), nil)
			},
			Delay: delay,
		}
	}
}

type ColorMatcher func(r uint8, g uint8, b uint8) bool

func GreenMatcher() ColorMatcher {
	return func(r uint8, g uint8, b uint8) bool {
		if uint16(r) > 100 && uint16(r) > uint16(g)*2 && uint16(r) > uint16(b)*2 {
			return false
		}
		if uint16(g) > 100 && uint16(g) > uint16(r)*2 && uint16(g) > uint16(b)*2 {
			return true
		}
		if uint16(b) > 100 && uint16(b) > uint16(g)*2 && uint16(b) > uint16(r)*2 {
			return false
		}
		return false
	}
}

func findColorAndClick(l logrus.FieldLogger) func(bot identity.Identity, startX int, startY int, endX int, endY int, cmFunc ColorMatcher) repeatableImageFunc {
	return func(bot identity.Identity, startX int, startY int, endX int, endY int, cmFunc ColorMatcher) repeatableImageFunc {
		return func(i image.Image) bool {
			x := startX
			y := startY

			for x <= endX && y <= endY {
				nrgb := i.At(x, y).(color.NRGBA)
				if cmFunc(nrgb.R, nrgb.G, nrgb.B) {
					click(x, y)(bot)
					return false
				}
				x += 1
				if x > endX {
					x = startX
					y += 1
				}
			}
			l.Errorf("Unable to locate free check in for bot. Not handling recovery step for now.")
			return false
		}
	}
}

func sleep(d time.Duration) Operator {
	return func(bot identity.Identity) {
		time.Sleep(d)
	}
}

func initMD(d time.Duration) Operator {
	return func(bot identity.Identity) {
		launchMD(bot)
		time.Sleep(d)
	}
}

func simpleGet(url string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	_, _ = http.DefaultClient.Do(req)
}

func killAndRestartEmulator(l logrus.FieldLogger) func(bot identity.Identity) {
	noxExe, ok := os.LookupEnv("NOX_EXE")
	if !ok {
		l.Fatalf("Unable to lookup NOX_EXE")
	}

	return func(bot identity.Identity) {
		l.Warnf("Something prevented the emulator for %s from starting appropriately. Killing the processes and restarting.", bot.Name())
		ps, err := process.Processes()
		if err != nil {
			return
		}
		for _, p := range ps {
			cmd, err := p.Cmdline()
			if err != nil {
				continue
			}
			if strings.Contains(cmd, bot.DeviceId()) {
				println(p.Pid)
				_ = p.Kill()
			}
		}

		time.Sleep(500 * time.Millisecond)

		launchCmd := exec.Command(noxExe, "-clone:"+bot.NoxId(), "-package:com.arlosoft.macrodroid")
		err = launchCmd.Start()
		if err != nil {
			l.WithError(err).Errorf("Unable to launch NOX instance for %s.", bot.Name())
		}
		time.Sleep(500 * time.Millisecond)
	}
}
