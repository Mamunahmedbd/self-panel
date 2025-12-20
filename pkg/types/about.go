package types

type (
	AboutData struct {
		SupportEmail string
		Title        string
		Subtitle     string
		Description  string
		Stats        []AboutStat
		Mission      string
		Vision       string
		Features     []AboutFeature
	}

	AboutStat struct {
		Value string
		Label string
	}

	AboutFeature struct {
		Title       string
		Description string
		Icon        string
	}
)
