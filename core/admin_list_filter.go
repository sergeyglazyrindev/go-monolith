package core

import (
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ListFilter struct {
	Title             string
	URLFilteringParam string
	OptionsToShow     []*FieldChoice
	FetchOptions      func(m interface{}) []*FieldChoice
	CustomFilterQs    func(afo IAdminFilterObjects, filterString string)
	Template          string
	Ordering          int
}

func (lf *ListFilter) FilterQs(afo IAdminFilterObjects, filterString string) {
	if lf.CustomFilterQs != nil {
		lf.CustomFilterQs(afo, filterString)
	} else {
		afo.FilterQs(filterString)
	}
}

func (lf *ListFilter) IsItActive(fullURL *url.URL) bool {
	return strings.Contains(fullURL.String(), lf.URLFilteringParam)
}

func (lf *ListFilter) GetURLToClearFilter(fullURL *url.URL) string {
	clonedURL := CloneNetURL(fullURL)
	qs := clonedURL.Query()
	qs.Del(lf.URLFilteringParam)
	clonedURL.RawQuery = qs.Encode()
	return clonedURL.String()
}

func (lf *ListFilter) IsThatOptionActive(option *FieldChoice, fullURL *url.URL) bool {
	qs := fullURL.Query()
	value := qs.Get(lf.URLFilteringParam)
	if value != "" {
		optionValue := TransformValueForListDisplay(option.Value, true)
		if optionValue == value {
			return true
		}
	}
	return false
}

func (lf *ListFilter) GetURLForOption(option *FieldChoice, fullURL *url.URL) string {
	clonedURL := CloneNetURL(fullURL)
	qs := clonedURL.Query()
	qs.Set(lf.URLFilteringParam, TransformValueForListDisplay(option.Value, true))
	clonedURL.RawQuery = qs.Encode()
	return clonedURL.String()
}

type ListFilterRegistry struct {
	ListFilter []*ListFilter
}

type ListFilterList []*ListFilter

func (apl ListFilterList) Len() int { return len(apl) }
func (apl ListFilterList) Less(i, j int) bool {
	return apl[i].Ordering < apl[j].Ordering
}
func (apl ListFilterList) Swap(i, j int) { apl[i], apl[j] = apl[j], apl[i] }

func (lfr *ListFilterRegistry) Iterate() <-chan *ListFilter {
	chnl := make(chan *ListFilter)
	go func() {
		lfList := make(ListFilterList, 0)
		defer close(chnl)
		for _, lF := range lfr.ListFilter {
			lfList = append(lfList, lF)
		}
		sort.Slice(lfList, func(i int, j int) bool {
			return lfList[i].Ordering < lfList[j].Ordering
		})
		for _, lf := range lfList {
			chnl <- lf
		}
	}()
	return chnl
}

func (lfr *ListFilterRegistry) IsEmpty() bool {
	return !(len(lfr.ListFilter) > 0)
}

func (lfr *ListFilterRegistry) Add(lf *ListFilter) {
	lfr.ListFilter = append(lfr.ListFilter, lf)
}

type DisplayFilterOption struct {
	FilterField string
	FilterValue string
	DisplayAs   string
}

type FilterOption struct {
	FieldName    string
	FetchOptions func(afo IAdminFilterObjects) []*DisplayFilterOption
}

type FilterOptionsRegistry struct {
	FilterOption []*FilterOption
}

func (for1 *FilterOptionsRegistry) AddFilterOption(fo *FilterOption) {
	for1.FilterOption = append(for1.FilterOption, fo)
}

func (for1 *FilterOptionsRegistry) GetAll() <-chan *FilterOption {
	chnl := make(chan *FilterOption)
	go func() {
		defer close(chnl)
		for _, fo := range for1.FilterOption {
			chnl <- fo
		}
	}()
	return chnl
}

func NewFilterOptionsRegistry() *FilterOptionsRegistry {
	return &FilterOptionsRegistry{FilterOption: make([]*FilterOption, 0)}
}

func NewFilterOption() *FilterOption {
	return &FilterOption{}
}

func FetchOptionsFromGormModelFromDateTimeField(afo IAdminFilterObjects, filterOptionField string) []*DisplayFilterOption {
	ret := make([]*DisplayFilterOption, 0)
	database := NewDatabaseInstance()
	defer database.Close()
	filterString := database.Adapter.GetStringToExtractYearFromField(filterOptionField)
	rows, _ := afo.GetInitialQuerySet().Select(filterString + " as year, count(*) as total").Group(filterString).Rows()
	var filterValue uint
	var filterCount uint
	for rows.Next() {
		rows.Scan(&filterValue, &filterCount)
		filterString := strconv.Itoa(int(filterValue))
		ret = append(ret, &DisplayFilterOption{
			FilterField: filterOptionField,
			FilterValue: filterString,
			DisplayAs:   filterString,
		})
	}
	if len(ret) < 2 {
		ret = make([]*DisplayFilterOption, 0)
		filterString := database.Adapter.GetStringToExtractMonthFromField(filterOptionField)
		rows, _ := afo.GetInitialQuerySet().Select(filterString + " as month, count(*) as total").Group(filterString).Rows()
		var filterValue uint
		var filterCount uint
		for rows.Next() {
			rows.Scan(&filterValue, &filterCount)
			filterString := strconv.Itoa(int(filterValue))
			filteredMonth, _ := strconv.Atoi(filterString)
			ret = append(ret, &DisplayFilterOption{
				FilterField: filterOptionField,
				FilterValue: filterString,
				DisplayAs:   time.Month(filteredMonth).String(),
			})
		}
	}
	return ret
}
