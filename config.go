package drum

var (
	tryStep = 5
	maxStep = 300
)

type DrumConfig func()

func SetTryStep(step int) DrumConfig {
	return func() {
		tryStep = step
	}
}

func SetMaxStep(step int) DrumConfig {
	return func() {
		maxStep = step
	}
}
