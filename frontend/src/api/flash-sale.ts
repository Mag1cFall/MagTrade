import http from '@/utils/http'
import type { ApiResponse, FlashSaleListResponse, FlashSaleDetail, RushResponse } from '@/types'

// 获取秒杀活动列表
export const getFlashSales = (page = 1, pageSize = 20, status?: number) => {
  const params: any = { page, page_size: pageSize }
  if (status !== undefined) params.status = status
  return http.get<any, ApiResponse<FlashSaleListResponse>>('/flash-sales', { params })
}

// 获取活动详情
export const getFlashSaleDetail = (id: number) => {
  return http.get<any, ApiResponse<FlashSaleDetail>>(`/flash-sales/${id}`)
}

// 获取实时库存
export const getFlashSaleStock = (id: number) => {
  return http.get<any, ApiResponse<{ stock: number }>>(`/flash-sales/${id}/stock`)
}

// 抢购
export const rushFlashSale = (id: number, quantity: number = 1) => {
  return http.post<any, ApiResponse<RushResponse>>(`/flash-sales/${id}/rush`, { quantity })
}