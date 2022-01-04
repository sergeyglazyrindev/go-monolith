package core

import "github.com/pilagod/gorm-cursor-paginator/v2/paginator"

type GORMPaginatedResponse struct {
	Count        int64
	NextLink     string
	PrevLink     string
}

func CreateAPIResultPaginator(
	cursor paginator.Cursor,
) *paginator.Paginator {
	opts := []paginator.Option{
		&paginator.Config{
			Rules: []paginator.Rule{{
				Key:     "ID",
			}},
			Limit: 12,
			Order: paginator.ASC,
		},
	}
	if cursor.After != nil {
		opts = append(opts, paginator.WithAfter(*cursor.After))
	}
	paginator1 := paginator.New(opts...)

	return paginator1
}
