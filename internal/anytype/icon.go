package anytype

type IconFormat string

const (
	IconFormatEmoji IconFormat = "emoji"
	IconFormatFile  IconFormat = "file"
	IconFormatIcon  IconFormat = "icon"
)

type Icon struct {
	Emoji  string     `json:"emoji,omitempty"`
	Icon   string     `json:"icon,omitempty"`
	File   string     `json:"file,omitempty"`
	Format IconFormat `json:"format,omitempty"`
}
