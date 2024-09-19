package paging

import "github.com/gogf/gf/v2/util/gconv"

type SimplePageOpsImpl[T any] struct {
	currentPage  int
	perPageItems int
	prevPage     int
	nextPage     int
	totalPages   int
	totalItems   int
	pageData     *[]T
}

func SimplePaging[T any](currentPage int, perPageItems int) *SimplePageOpsImpl[T] {
	if currentPage < 1 {
		currentPage = 1
	}
	if perPageItems < 1 {
		perPageItems = 10
	}
	return &SimplePageOpsImpl[T]{currentPage: currentPage, perPageItems: perPageItems}
}

func (paging *SimplePageOpsImpl[T]) Offset() int {
	return (paging.currentPage - 1) * paging.Limit()
}

func (paging *SimplePageOpsImpl[T]) Limit() int {
	if paging.perPageItems > 1 {
		return paging.perPageItems
	}
	return 1
}

func (paging *SimplePageOpsImpl[T]) PageData(totalItems int, pageData *[]T) *SimplePageOpsImpl[T] {
	paging.totalItems = totalItems
	paging.pageData = pageData
	return paging
}

func (paging *SimplePageOpsImpl[T]) PageModels(totalItems int, models *[]any) *SimplePageOpsImpl[T] {
	paging.totalItems = totalItems
	data := new([]T)
	gconv.Scan(models, data)
	paging.pageData = data
	return paging
}

func (paging *SimplePageOpsImpl[T]) reCalculate() {
	if paging.totalItems%paging.perPageItems > 0 {
		paging.totalPages = (paging.totalItems / paging.perPageItems) + 1
	} else {
		paging.totalPages = paging.totalItems / paging.perPageItems
	}
	paging.prevPage = doMin(doMin(doMax(paging.currentPage-1, 0), paging.currentPage), paging.totalPages)
	paging.nextPage = doMin(paging.currentPage+1, paging.totalPages)
}

func (paging *SimplePageOpsImpl[T]) Fetch() *PageEntityRes[T] {
	if paging.totalItems == 0 {
		data := make([]T, 0)
		return &PageEntityRes[T]{
			CurrentPage:  0,
			PrevPage:     0,
			NextPage:     0,
			TotalPages:   0,
			PerPageItems: paging.perPageItems,
			TotalItems:   0,
			PageData:     &data,
		}
	}
	paging.reCalculate()
	return &PageEntityRes[T]{
		CurrentPage:  paging.currentPage,
		PrevPage:     paging.prevPage,
		NextPage:     paging.nextPage,
		TotalPages:   paging.totalPages,
		PerPageItems: paging.perPageItems,
		TotalItems:   paging.totalItems,
		PageData:     paging.pageData,
	}
}
