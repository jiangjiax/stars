package template

type Pagination struct {
	CurrentPage int
	PageSize    int
	TotalPosts  int
	TotalPages  int
	BaseURL     string
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
}

func NewPagination(currentPage, pageSize, totalPosts int, baseURL string) *Pagination {
	totalPages := (totalPosts + pageSize - 1) / pageSize

	p := &Pagination{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		TotalPosts:  totalPosts,
		TotalPages:  totalPages,
		BaseURL:     baseURL,
		HasPrev:     currentPage > 1,
		HasNext:     currentPage < totalPages,
	}

	if p.HasPrev {
		p.PrevPage = currentPage - 1
	}
	if p.HasNext {
		p.NextPage = currentPage + 1
	}

	return p
}
