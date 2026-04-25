package anytype

type PropertyLinkWithValue struct {
	Key         string   `json:"key"`
	ID          string   `json:"id"`
	Text        string   `json:"text,omitempty"`
	Number      float64  `json:"number,omitempty"`
	Checkbox    bool     `json:"checkbox,omitempty"`
	Date        int64    `json:"date,omitempty"`
	Select      string   `json:"select,omitempty"`
	MultiSelect []string `json:"multi_select,omitempty"`
	URL         string   `json:"url,omitempty"`
	Email       string   `json:"email,omitempty"`
	Phone       string   `json:"phone,omitempty"`
}

func NewTextProperty(key, value string) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key:  key,
		Text: value,
	}
}

func NewNumberProperty(key string, value float64) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key:    key,
		Number: value,
	}
}

func NewCheckboxProperty(key string, value bool) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key:      key,
		Checkbox: value,
	}
}

func NewDateProperty(key string, timestamp int64) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key:  key,
		Date: timestamp,
	}
}

func NewSelectProperty(key, tagID string) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key:    key,
		Select: tagID,
	}
}

func NewMultiSelectProperty(key string, tagIDs []string) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key:         key,
		MultiSelect: tagIDs,
	}
}

func NewURLProperty(key, url string) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key: key,
		URL: url,
	}
}

func NewEmailProperty(key, email string) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key:   key,
		Email: email,
	}
}

func NewPhoneProperty(key, phone string) PropertyLinkWithValue {
	return PropertyLinkWithValue{
		Key:   key,
		Phone: phone,
	}
}
