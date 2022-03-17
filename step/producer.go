package step

import (
	"github.com/sirupsen/logrus"
	"gos-auto-bot/coordinate"
	"time"
)

type Producer func() []Step

func CreateImperialCarnivalHelpSteps(l logrus.FieldLogger) Producer {
	return func() []Step {
		steps := createLaunchSteps(l)
		steps = append(steps,
			ClickStep(coordinate.NewScaled(500, 435), 750),
			ClickStep(coordinate.NewScaled(275, 740), 750),
			ClickStep(coordinate.NewScaled(455, 830), 750),
			ClickStep(coordinate.NewScaled(383, 870), 750),
			ClickStep(coordinate.NewScaled(430, 335), 750),
		)
		return steps
	}
}

func CreateUnionSteps(l logrus.FieldLogger) Producer {
	return func() []Step {
		steps := createLaunchSteps(l)
		steps = append(steps,
			ClickStep(coordinate.NewScaled(353, 464), 500),
			GestureStep(130, 820, 450, 820, 750, 2000),
			ClickStep(coordinate.NewScaled(378, 289), 500),
			ClickStep(coordinate.NewScaled(447, 622), 500),
			ClickStep(coordinate.NewScaled(30, 58), 500),
			ClickStep(coordinate.NewScaled(30, 58), 500),
		)
		return steps
	}
}

func CreateRankingsSteps(l logrus.FieldLogger) Producer {
	return func() []Step {
		steps := createLaunchSteps(l)
		steps = append(steps,
			GestureStep(130, 820, 450, 820, 750, 4000),
			ClickStep(coordinate.NewScaled(316, 324), 750),
			ClickStep(coordinate.NewScaled(272, 61), 1000),
			ClickStep(coordinate.NewScaled(282, 300), 2000),
			ClickStep(coordinate.NewScaled(450, 923), 750), // honor
			ClickStep(coordinate.NewScaled(210, 115), 750), // union
			ClickStep(coordinate.NewScaled(205, 116), 750), // off
			ClickStep(coordinate.NewScaled(450, 923), 750), // honor
			ClickStep(coordinate.NewScaled(340, 115), 750), // intimacy
			ClickStep(coordinate.NewScaled(346, 116), 750), // off
			ClickStep(coordinate.NewScaled(450, 923), 750), // honor
			ClickStep(coordinate.NewScaled(475, 115), 750), // pet
			ClickStep(coordinate.NewScaled(470, 116), 750), // off
			ClickStep(coordinate.NewScaled(450, 923), 750), // honor
			ClickStep(coordinate.NewScaled(272, 61), 750),
			ClickStep(coordinate.NewScaled(30, 58), 750),
			ClickStep(coordinate.NewScaled(282, 668), 2000),
			ClickStep(coordinate.NewScaled(450, 923), 750), // honor
			ClickStep(coordinate.NewScaled(210, 115), 750), // union
			ClickStep(coordinate.NewScaled(205, 116), 750), // off
			ClickStep(coordinate.NewScaled(450, 923), 750), // honor
			ClickStep(coordinate.NewScaled(340, 115), 750), // intimacy
			ClickStep(coordinate.NewScaled(346, 116), 750), // off
			ClickStep(coordinate.NewScaled(450, 923), 750), // honor
			ClickStep(coordinate.NewScaled(475, 115), 750), // pet
			ClickStep(coordinate.NewScaled(470, 116), 750), // off
			ClickStep(coordinate.NewScaled(450, 923), 750), // honor
			ClickStep(coordinate.NewScaled(272, 61), 750),
			ClickStep(coordinate.NewScaled(30, 58), 750),
			ClickStep(coordinate.NewScaled(30, 58), 500),
		)
		return steps
	}
}

func CreateHaremSteps(l logrus.FieldLogger) func(haremAmount int) Producer {
	return func(haremAmount int) Producer {
		return func() []Step {
			steps := createLaunchSteps(l)
			steps = append(steps,
				GestureStep(130, 820, 450, 820, 750, 2000),
				ClickStep(coordinate.NewScaled(535, 357), 750),
				ClickStep(coordinate.NewScaled(225, 645), 500),
			)
			for i := 0; i < haremAmount; i++ {
				steps = append(steps,
					ClickStep(coordinate.NewScaled(440, 918), 5000),
					ClickStep(coordinate.NewScaled(276, 86), 500),
				)
			}
			steps = append(steps,
				ClickStep(coordinate.NewScaled(30, 58), 500),
				ClickStep(coordinate.NewScaled(30, 58), 500),
			)
			return steps
		}
	}
}

func CreateDivinationSteps(l logrus.FieldLogger) Producer {
	return func() []Step {
		steps := createLaunchSteps(l)
		steps = append(steps,
			GestureStep(400, 820, 240, 820, 750, 2000),
			ClickStep(coordinate.NewScaled(215, 568), 500),
			ClickStep(coordinate.NewScaled(270, 618), 500),
			ClickStep(coordinate.NewScaled(270, 618), 500),
			ClickStep(coordinate.NewScaled(30, 58), 500),
		)
		return steps
	}
}

func CreateCouncilSteps(l logrus.FieldLogger) func(vip int, envoys int) Producer {
	return func(vip int, envoys int) Producer {
		return func() []Step {
			steps := createLaunchSteps(l)
			if vip < 1 {
				steps = append(steps,
					GestureStep(400, 820, 240, 820, 750, 2000),
					ClickStep(coordinate.NewScaled(247, 376), 500),
				)
				initial := envoys
				if envoys > 5 {
					initial = 5
				}
				for i := 0; i < initial; i++ {
					steps = append(steps,
						ClickStep(coordinate.NewScaled(60+(i*100), 675), 1000),
						ClickStep(coordinate.NewScaled(275, 345), 1000),
						ClickStep(coordinate.NewScaled(275, 918), 1000),
						ClickStep(coordinate.NewScaled(275, 918), 1000),
						ClickStep(coordinate.NewScaled(275, 893), 1000),
						ClickStep(coordinate.NewScaled(275, 893), 1000),
					)
				}
				if envoys > 5 {
					steps = append(steps, GestureStep(510, 720, 0, 720, 750, 2000))
					next := envoys - 5
					for i := 0; i < next; i++ {
						steps = append(steps,
							ClickStep(coordinate.NewScaled(290+(i*100), 675), 1000),
							ClickStep(coordinate.NewScaled(275, 345), 1000),
							ClickStep(coordinate.NewScaled(275, 918), 1000),
							ClickStep(coordinate.NewScaled(275, 918), 1000),
							ClickStep(coordinate.NewScaled(275, 893), 1000),
							ClickStep(coordinate.NewScaled(275, 893), 1000),
						)
					}
				}
				steps = append(steps,
					ClickStep(coordinate.NewScaled(30, 58), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
				)
			} else {
				steps = append(steps,
					GestureStep(400, 820, 240, 820, 750, 2000),
					ClickStep(coordinate.NewScaled(247, 376), 500),
					ClickStep(coordinate.NewScaled(260, 668), 500),
					ClickStep(coordinate.NewScaled(260, 308), 500),
					ClickStep(coordinate.NewScaled(275, 918), 500),
					ClickStep(coordinate.NewScaled(275, 918), 500),
					ClickStep(coordinate.NewScaled(275, 893), 500),
					ClickStep(coordinate.NewScaled(275, 893), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
				)
			}
			return steps
		}
	}
}

func CreateCoalitionSteps(l logrus.FieldLogger) Producer {
	return func() []Step {
		steps := createLaunchSteps(l)
		steps = append(steps,
			GestureStep(400, 820, 240, 820, 750, 2000),
			ClickStep(coordinate.NewScaled(429, 246), 500),
			ClickStep(coordinate.NewScaled(94, 415), 500),
			ClickStep(coordinate.NewScaled(41, 714), 500),
			ClickStep(coordinate.NewScaled(270, 700), 500),
			ClickStep(coordinate.NewScaled(270, 779), 500),
			ClickStep(coordinate.NewScaled(270, 779), 500),
			ClickStep(coordinate.NewScaled(30, 58), 500),
			ClickStep(coordinate.NewScaled(30, 58), 500),
			ClickStep(coordinate.NewScaled(30, 58), 500),
		)
		return steps
	}
}

func CreateAdSteps(l logrus.FieldLogger) func(adsOffset int, novice bool, sevenDays bool) Producer {
	return func(adsOffset int, novice bool, sevenDays bool) Producer {
		return func() []Step {
			steps := createLaunchSteps(l)
			base := 179
			offset := 72
			index := adsOffset - 1
			if novice {
				index += 1
			}
			if sevenDays {
				index += 1
			}
			y := base + (index * offset)
			steps = append(steps, ClickStep(coordinate.NewScaled(40, y), 5000))

			for i := 1; i < 5; i++ {
				steps = append(steps,
					ClickStep(coordinate.NewScaled(277, 660), 35000),
					BackStep(1000),
					ClickStep(coordinate.NewScaled(277, 660), 5000),
					ClickStep(coordinate.NewScaled(277, 256), 500),
				)
			}
			steps = append(steps, ClickStep(coordinate.NewScaled(514, 264), 500))
			return steps
		}
	}
}

func CreateAcademySteps(l logrus.FieldLogger) func(academySeats int) Producer {
	return func(academySeats int) Producer {
		return func() []Step {
			steps := createLaunchSteps(l)
			if academySeats < 5 {
				steps = append(steps,
					GestureStep(400, 820, 240, 820, 750, 2000),
					PrintStep(l)("Moved right."),
					ClickStep(coordinate.NewScaled(403, 483), 500),
					ClickStep(coordinate.NewScaled(270, 60), 500),
					ClickStep(coordinate.NewScaled(135, 185), 750),
					ClickStep(coordinate.NewScaled(105, 275), 750),
					ClickStep(coordinate.NewScaled(265, 935), 750),
					ClickStep(coordinate.NewScaled(400, 185), 750),
					ClickStep(coordinate.NewScaled(105, 275), 750),
					ClickStep(coordinate.NewScaled(265, 935), 750),
					ClickStep(coordinate.NewScaled(135, 400), 750),
					ClickStep(coordinate.NewScaled(105, 275), 750),
					ClickStep(coordinate.NewScaled(265, 935), 750),
					ClickStep(coordinate.NewScaled(400, 400), 750),
					ClickStep(coordinate.NewScaled(105, 275), 750),
					ClickStep(coordinate.NewScaled(265, 935), 750),
					ClickStep(coordinate.NewScaled(30, 58), 500),
				)
			} else {
				steps = append(steps,
					GestureStep(400, 820, 240, 820, 250, 1500),
					ClickStep(coordinate.NewScaled(403, 483), 500),
					ClickStep(coordinate.NewScaled(270, 60), 500),
					ClickStep(coordinate.NewScaled(360, 918), 500),
					ClickStep(coordinate.NewScaled(270, 60), 500),
					ClickStep(coordinate.NewScaled(360, 918), 2000),
					ClickStep(coordinate.NewScaled(500, 868), 500),
					ClickStep(coordinate.NewScaled(380, 800), 500),
					ClickStep(coordinate.NewScaled(270, 60), 500),
					ClickStep(coordinate.NewScaled(514, 138), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
				)
			}
			return steps
		}
	}
}

func CreateLevySteps(l logrus.FieldLogger) func(vip int, levyTotal int) Producer {
	return func(vip int, levyTotal int) Producer {
		return func() []Step {
			steps := createLaunchSteps(l)
			if vip < 1 {
				steps = append(steps,
					ClickStep(coordinate.NewScaled(220, 228), 500),
					ClickStep(coordinate.NewScaled(94, 571), 500),
					ClickStep(coordinate.NewScaled(265, 918), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
					ClickStep(coordinate.NewScaled(422, 571), 500),
				)
				for i := 1; i < levyTotal; i++ {
					steps = append(steps, ClickStep(coordinate.NewScaled(375, 768), 750))
				}
				steps = append(steps,
					ClickStep(coordinate.NewScaled(30, 58), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
				)
			} else {
				steps = append(steps,
					ClickStep(coordinate.NewScaled(220, 228), 500),
					ClickStep(coordinate.NewScaled(94, 571), 500),
					ClickStep(coordinate.NewScaled(265, 918), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
					ClickStep(coordinate.NewScaled(422, 571), 500),
					ClickStep(coordinate.NewScaled(375, 768), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
					ClickStep(coordinate.NewScaled(30, 58), 500),
				)
			}
			return steps
		}
	}
}

func CreateDailyCheckInSteps(l logrus.FieldLogger) Producer {
	return func() []Step {
		steps := createLaunchSteps(l)
		steps = append(steps,
			ClickStep(coordinate.NewScaled(355, 76), 750),
			ClickStep(coordinate.NewScaled(469, 123), 750),
			ClickColorStep(l)(coordinate.NewScaled(20, 305), coordinate.NewScaled(525, 550), GreenMatcher(), 500),
			ClickStep(coordinate.NewScaled(30, 58), 500),
		)
		return steps
	}
}

func createLaunchSteps(l logrus.FieldLogger) []Step {
	steps := make([]Step, 0)
	steps = append(steps,
		VerifyPixelStep(l)(coordinate.NewScaled(271, 361), 245, 133, 73, initMD(time.Duration(1)*time.Second), killAndRestartEmulator(l)),
		ClickStep(coordinate.NewScaled(269, 354), 500), // export import
		ClickStep(coordinate.NewScaled(147, 340), 500), // import
		ClickStep(coordinate.NewScaled(100, 183), 500), // top left
		ClickStep(coordinate.NewScaled(432, 543), 500),
		LaunchStep(),
		VerifyPixelStep(l)(coordinate.NewScaled(523, 945), 186, 185, 185, sleep(time.Duration(1)*time.Second), launchGame),
		ClickStep(coordinate.NewScaled(260, 850), 5000),
		VerifyPixelStep(l)(coordinate.NewScaled(287, 495), 125, 36, 5, back, nil),
		PrintStep(l)("Finished launching game."),
	)
	return steps
}
