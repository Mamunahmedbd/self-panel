package types

type (
	BlogPage struct {
		Title      string
		Subtitle   string
		Posts      []BlogPost
		Categories []string
	}

	BlogPost struct {
		ID       int
		Title    string
		Summary  string
		Content  string
		ImageURL string
		Category string
		Author   string
		Date     string
		ReadTime string
		Featured bool
	}
)
