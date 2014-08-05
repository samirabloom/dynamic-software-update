package stages

type TransitionMode int64

const (
	InstantMode TransitionMode = 1
	SessionMode TransitionMode = 2
	GradualMode TransitionMode = 3
	ConcurrentMode TransitionMode = 4
)

var ModesCodeToMode = map[string]TransitionMode {
	"INSTANT": InstantMode,
	"SESSION": SessionMode,
	"GRADUAL": GradualMode,
	"CONCURRENT": ConcurrentMode,
}

var ModesModeToCode = map[TransitionMode]string {
	InstantMode: "INSTANT",
	SessionMode: "SESSION",
	GradualMode: "GRADUAL",
	ConcurrentMode: "CONCURRENT",
}


