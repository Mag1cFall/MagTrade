import http from '@/utils/http'
import type { ApiResponse, OrderListResponse, Order } from '@/types'

export const getOrders = (page = 1, pageSize = 20) => {
  return http.get<any, ApiResponse<OrderListResponse>>('/orders', {
    params: { page, page_size: pageSize }
  })
}

export const getOrderDetail = (orderNo: string) => {
  return http.get<any, ApiResponse<Order>>(`/orders/${orderNo}`)
}

export const payOrder = (orderNo: string) => {
  return http.post<any, ApiResponse<Order>>(`/orders/${orderNo}/pay`)
}

export const cancelOrder = (orderNo: string) => {
  return http.post<any, ApiResponse<Order>>(`/orders/${orderNo}/cancel`)
}