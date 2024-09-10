package paging

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"
)

type PageEntityRes[T any] struct {
	CurrentPage  int  `json:"currentPage"`
	PrevPage     int  `json:"prevPage"`
	NextPage     int  `json:"nextPage"`
	TotalPages   int  `json:"totalPages"`
	PerPageItems int  `json:"perPageItems"`
	TotalItems   int  `json:"totalItems"`
	PageData     *[]T `json:"pageData"`
}

type PageOpsImpl[T any, V any, Z any] struct {
	currentPage  int
	prevPage     int
	nextPage     int
	totalPages   int
	perPageItems int
	totalItems   int
	pageData     *[]T
	params       *V
	models       *[]Z
	query        func(ctx context.Context, offset int, limit int, params *V) (*[]Z, *int, error)
	converter    func(ctx context.Context, source *[]Z, target *[]T)
	ctx          context.Context
}

func Paging[T any, V any, Z any]() *PageOpsImpl[T, V, Z] {
	return &PageOpsImpl[T, V, Z]{currentPage: 1, perPageItems: 10}
}

func (paging *PageOpsImpl[T, V, Z]) CurrentPage(currentPage int) *PageOpsImpl[T, V, Z] {
	if currentPage < 1 {
		paging.currentPage = 1
	} else {
		paging.currentPage = currentPage
	}
	return paging

}

func (paging *PageOpsImpl[T, V, Z]) PerPageItems(perPageItems int) *PageOpsImpl[T, V, Z] {
	if perPageItems < 1 {
		paging.perPageItems = 10
	} else {
		paging.perPageItems = perPageItems
	}
	return paging
}

func (paging *PageOpsImpl[T, V, Z]) Params(params *V) *PageOpsImpl[T, V, Z] {
	if params != nil {
		paging.params = params
	}
	return paging
}

func (paging *PageOpsImpl[T, V, Z]) Query(f func(ctx context.Context, offset int, limit int, params *V) (*[]Z, *int, error)) *PageOpsImpl[T, V, Z] {
	if f != nil {
		paging.query = f
	}
	return paging
}

func (paging *PageOpsImpl[T, V, Z]) Converter(f func(ctx context.Context, source *[]Z, target *[]T)) *PageOpsImpl[T, V, Z] {
	if f != nil {
		paging.converter = f
	}
	return paging
}

func (paging *PageOpsImpl[T, V, Z]) Ctx(ctx context.Context) *PageOpsImpl[T, V, Z] {
	if ctx != nil {
		paging.ctx = ctx
	}
	return paging
}

func (paging *PageOpsImpl[T, V, Z]) Fetch() (*PageEntityRes[T], error) {
	items, totalItems, err := paging.query(paging.ctx, paging.Offset(), paging.Limit(), paging.params)
	if err != nil {
		return nil, err
	}
	if *totalItems == 0 {
		data := make([]T, 0)
		return &PageEntityRes[T]{
			CurrentPage:  0,
			PrevPage:     0,
			NextPage:     0,
			TotalPages:   0,
			PerPageItems: paging.perPageItems,
			TotalItems:   0,
			PageData:     &data,
		}, nil
	}
	paging.reCalculate(*totalItems)
	paging.pageData = new([]T)
	if paging.converter != nil {
		paging.converter(paging.ctx, items, paging.pageData)
	} else {
		err = gconv.Scan(items, paging.pageData)
	}
	if err != nil {
		return nil, err
	}
	return &PageEntityRes[T]{
		CurrentPage:  paging.currentPage,
		PrevPage:     paging.prevPage,
		NextPage:     paging.nextPage,
		TotalPages:   paging.totalPages,
		PerPageItems: paging.perPageItems,
		TotalItems:   paging.totalItems,
		PageData:     paging.pageData,
	}, nil
}

func (paging *PageOpsImpl[T, V, Z]) Offset() int {
	return (paging.currentPage - 1) * paging.Limit()
}

func (paging *PageOpsImpl[T, V, Z]) Limit() int {
	if paging.perPageItems > 1 {
		return paging.perPageItems
	}
	return 1
}

func (paging *PageOpsImpl[T, V, Z]) reCalculate(totalItems int) {
	paging.totalItems = totalItems
	if paging.totalItems%paging.perPageItems > 0 {
		paging.totalPages = (paging.totalItems / paging.perPageItems) + 1
	} else {
		paging.totalPages = paging.totalItems / paging.perPageItems
	}
	paging.prevPage = doMin(doMin(doMax(paging.currentPage-1, 0), paging.currentPage), paging.totalPages)
	paging.nextPage = doMin(paging.currentPage+1, paging.totalPages)
}
func doMax(source int, target int) int {
	if source >= target {
		return source
	}
	return target
}

func doMin(source int, target int) int {
	if source <= target {
		return source
	}
	return target
}
