package ds

const (
	MODE_Base = iota
	MODE_Beta
	MODE_Speciale
)

const (
	MODEL_Chat = iota
	MODEL_Reasoner
)

type Cli struct {
	key string
	url string
}

func NewDeepseek(key string, mode int) *Cli {
	var url string
	if mode == MODE_Base {
		url = "https://api.deepseek.com"
	}
	if mode == MODE_Beta {
		url = "https://api.deepseek.com/beta"
	}
	if mode == MODE_Speciale {
		url = "https://api.deepseek.com/v3.2_speciale_expires_on_20251215"
	}
	return &Cli{
		key: key,
		url: url,
	}
}
