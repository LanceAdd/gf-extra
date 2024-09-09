package paging

type DataPageOpsImpl[T any] struct {
	currentPage  int
	perPageItems int
	prevPage     int
	nextPage     int
	totalPages   int
	totalItems   int
	pageData     *[]T
}

func DataPaging[T any](pageData *[]T) *DataPageOpsImpl[T] {
	return &DataPageOpsImpl[T]{pageData: pageData, totalItems: len(*pageData), currentPage: 1, perPageItems: 10}
}

func (paging *DataPageOpsImpl[T]) CurrentPage(currentPage int) *DataPageOpsImpl[T] {
	if currentPage < 1 {
		paging.currentPage = 1
	} else {
		paging.currentPage = currentPage
	}
	return paging
}
func (paging *DataPageOpsImpl[T]) PerPageItems(perPageItems int) *DataPageOpsImpl[T] {
	if perPageItems < 1 {
		paging.perPageItems = 1
	} else {
		paging.perPageItems = perPageItems
	}
	return paging
}

func (paging *DataPageOpsImpl[T]) Offset() int {
	return (paging.currentPage - 1) * paging.Limit()
}

func (paging *DataPageOpsImpl[T]) Limit() int {
	if paging.perPageItems > 1 {
		return paging.perPageItems
	}
	return 1
}

func (paging *DataPageOpsImpl[T]) reCalculate() {
	if paging.totalItems%paging.perPageItems > 0 {
		paging.totalPages = (paging.totalItems / paging.perPageItems) + 1
	} else {
		paging.totalPages = paging.totalItems / paging.perPageItems
	}
	paging.prevPage = doMin(doMin(doMax(paging.currentPage-1, 0), paging.currentPage), paging.totalPages)
	paging.nextPage = doMin(paging.currentPage+1, paging.totalPages)
}

func (paging *DataPageOpsImpl[T]) Fetch() *PageEntityRes[T] {
	if paging.totalItems == 0 {
		return &PageEntityRes[T]{
			CurrentPage:  0,
			PrevPage:     0,
			NextPage:     0,
			TotalPages:   0,
			PerPageItems: paging.perPageItems,
			TotalItems:   0,
			PageData:     new([]T),
		}
	}
	paging.reCalculate()
	pageData := (*paging.pageData)[paging.Offset() : paging.Offset()+paging.Limit()]
	return &PageEntityRes[T]{
		CurrentPage:  paging.currentPage,
		PrevPage:     paging.prevPage,
		NextPage:     paging.nextPage,
		TotalPages:   paging.totalPages,
		PerPageItems: paging.perPageItems,
		TotalItems:   paging.totalItems,
		PageData:     &pageData,
	}
}
