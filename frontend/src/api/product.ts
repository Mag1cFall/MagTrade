import http from '@/utils/http'
import type { ApiResponse, ProductListResponse } from '@/types'

export const getProducts = (page = 1, pageSize = 100) => {
  return http.get<any, ApiResponse<ProductListResponse>>('/products', {
    params: { page, page_size: pageSize }
  })
}