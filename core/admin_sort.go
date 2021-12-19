package core

type ISortBy interface {
	Sort(afo IAdminFilterObjects, direction int)
	GetDirection() int
	SetSortCustomImplementation(func(afo IAdminFilterObjects, field *Field, direction int))
}

type SortBy struct {
	Direction                int // -1 descending order, 1 ascending order
	Field                    *Field
	CustomSortImplementation func(afo IAdminFilterObjects, field *Field, direction int)
}

func (sb *SortBy) Sort(afo IAdminFilterObjects, direction int) {
	if sb.CustomSortImplementation != nil {
		sb.CustomSortImplementation(afo, sb.Field, direction)
	} else {
		afo.SortBy(sb.Field, direction)
	}
}

func (sb *SortBy) SetSortCustomImplementation(customSortImplementation func(afo IAdminFilterObjects, field *Field, direction int)) {
	sb.CustomSortImplementation = customSortImplementation
}

func (sb *SortBy) GetDirection() int {
	return sb.Direction
}
