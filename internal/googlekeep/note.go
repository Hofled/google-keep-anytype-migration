package googlekeep

type Note struct {
	Color                   string        `json:"color"`
	IsTrashed               bool          `json:"isTrashed"`
	IsPinned                bool          `json:"isPinned"`
	IsArchived              bool          `json:"isArchived"`
	Annotations             []Annotations `json:"annotations,omitempty"`
	TextContent             string        `json:"textContent,omitempty"`
	Title                   string        `json:"title"`
	UserEditedTimestampUsec uint64        `json:"userEditedTimestampUsec"`
	CreatedTimestampUsec    uint64        `json:"createdTimestampUsec"`
	TextContentHtml         string        `json:"textContentHtml"`
	Labels                  []Label       `json:"labels,omitempty"`
	ListContent             []ListContent `json:"listContent,omitempty"`
	Attachments             []Attachment  `json:"attachments,omitempty"`
	Sharees                 []Sharee      `json:"sharees,omitempty"`
	Tasks                   []Task        `json:"tasks,omitempty"`
}

type Annotations struct {
	Description string `json:"description"`
	Source      string `json:"source"`
	Title       string `json:"title"`
	Url         string `json:"url"`
}

type Label struct {
	Name string `json:"name"`
}

type ListContent struct {
	TextHtml  string `json:"textHtml"`
	Text      string `json:"text"`
	IsChecked bool   `json:"isChecked"`
}

type Attachment struct {
	FilePath string `json:"filePath"`
	MimeType string `json:"mimeType"`
}

type Sharee struct {
	IsOwner    bool   `json:"isOwner"`
	ShareeType string `json:"type"`
	Email      string `json:"email"`
}

type Task struct {
	Id string `json:"id"`
}
