package pagination

type Page struct {
	Number int
	Size   int
}

type PageInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func (p *Page) GetLimits(total int) (int, int) {
	start := (p.Number - 1) * p.Size
	end := start + p.Size

	return start, end
}
