package core

type AdminBreadcrumb struct {
	Name     string
	URL      string
	IsActive bool
	Icon     string
}

type AdminBreadCrumbsRegistry struct {
	BreadCrumbs []*AdminBreadcrumb
}

func (abcr *AdminBreadCrumbsRegistry) AddBreadCrumb(breadcrumb *AdminBreadcrumb) {
	if len(abcr.BreadCrumbs) != 0 {
		abcr.BreadCrumbs[0].IsActive = false
	} else {
		breadcrumb.IsActive = true
	}
	abcr.BreadCrumbs = append(abcr.BreadCrumbs, breadcrumb)
}

func (abcr *AdminBreadCrumbsRegistry) GetAll() <-chan *AdminBreadcrumb {
	chnl := make(chan *AdminBreadcrumb)
	go func() {
		defer close(chnl)
		for _, adminBreadcrumb := range abcr.BreadCrumbs {
			chnl <- adminBreadcrumb
		}
	}()
	return chnl
}

func NewAdminBreadCrumbsRegistry() *AdminBreadCrumbsRegistry {
	ret := &AdminBreadCrumbsRegistry{BreadCrumbs: make([]*AdminBreadcrumb, 0)}
	return ret
}
