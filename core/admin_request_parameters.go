package core

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type AdminRequestParams struct {
	CreateSession     bool
	GenerateCSRFToken bool
	NeedAllLanguages  bool
	Paginator         *AdminRequestPaginator
	RequestURL        string
	Search            string
	Ordering          []string
}

func (arp *AdminRequestParams) GetOrdering() string {
	return strings.Join(arp.Ordering, ",")
}

func NewAdminRequestParams() *AdminRequestParams {
	return &AdminRequestParams{
		CreateSession:     true,
		GenerateCSRFToken: true,
		NeedAllLanguages:  false,
		Paginator:         &AdminRequestPaginator{},
	}
}

func NewAdminRequestParamsFromGinContext(ctx *gin.Context) *AdminRequestParams {
	ret := &AdminRequestParams{
		CreateSession:     true,
		GenerateCSRFToken: true,
		NeedAllLanguages:  false,
		Paginator:         &AdminRequestPaginator{},
	}
	if ctx == nil {
		return ret
	}
	if ctx.Query("perpage") != "" {
		perPage, _ := strconv.Atoi(ctx.Query("perpage"))
		ret.Paginator.PerPage = perPage
	} else {
		ret.Paginator.PerPage = CurrentConfig.D.GoMonolith.AdminPerPage
	}
	if ctx.Query("offset") != "" {
		offset, _ := strconv.Atoi(ctx.Query("offset"))
		ret.Paginator.Offset = offset
	}
	if ctx.Query("p") != "" {
		page, _ := strconv.Atoi(ctx.Query("p"))
		if page > 1 {
			ret.Paginator.Offset = (page - 1) * ret.Paginator.PerPage
		}
	}
	ret.RequestURL = ctx.Request.URL.String()
	// autocomplete search
	if ctx.Query("term") != "" {
		ret.Search = ctx.Query("term")
	} else {
		ret.Search = ctx.Query("search")
	}
	orderingParts := strings.Split(ctx.Query("initialOrder"), ",")
	currentOrder := ctx.Query("o")
	currentOrderNameWithoutDirection := currentOrder
	if strings.HasPrefix(currentOrderNameWithoutDirection, "-") {
		currentOrderNameWithoutDirection = currentOrderNameWithoutDirection[1:]
	}
	foundNewOrder := false
	for i, part := range orderingParts {
		if strings.HasPrefix(part, "-") {
			part = part[1:]
		}
		if part == currentOrderNameWithoutDirection {
			orderingParts[i] = currentOrder
			foundNewOrder = true
		}
	}
	if !foundNewOrder {
		orderingParts = append(orderingParts, currentOrder)
	}
	finalOrderingParts := make([]string, 0)
	for _, part := range orderingParts {
		if part != "" {
			finalOrderingParts = append(finalOrderingParts, part)
		}
	}
	ret.Ordering = finalOrderingParts
	return ret
}
