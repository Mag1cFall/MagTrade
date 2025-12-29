import http from '@/utils/http'
import type { ApiResponse, Product, FlashSale } from '@/types'

// 创建商品
export const createProduct = (data: any) => {
  return http.post<any, ApiResponse<Product>>('/admin/products', data)
}

// 创建秒杀活动
export const createFlashSale = (data: any) => {
  return http.post<any, ApiResponse<FlashSale>>('/admin/flash-sales', data)
}
