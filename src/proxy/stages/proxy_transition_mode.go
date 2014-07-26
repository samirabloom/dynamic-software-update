package stages

type TransitionMode int64

const (
	InstantMode TransitionMode = 1
	SessionMode TransitionMode = 2
)

var ModesCodeToMode = map[string]TransitionMode {
	"SESSION": SessionMode,
	"INSTANT": InstantMode,
}

var ModesModeToCode = map[TransitionMode]string {
	SessionMode: "SESSION",
	InstantMode: "INSTANT",
}


