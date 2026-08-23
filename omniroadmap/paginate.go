package omniroadmap

// paginate accumulates every page of a cursor-paginated productboard-go
// list call. fetch is called with the empty string for the first page,
// then with each successive NextCursor, until it returns "" (no next
// page).
func paginate[T any](fetch func(cursor string) (batch []T, nextCursor string, err error)) ([]T, error) {
	var all []T
	cursor := ""
	for {
		batch, next, err := fetch(cursor)
		if err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if next == "" {
			return all, nil
		}
		cursor = next
	}
}
