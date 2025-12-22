export interface User {
  id: number
  username: string
  email: string
  role: 'user' | 'admin'
  status: number
  email_verified: boolean
  created_at: string
}

export interface Product {
  id: number
  name: string
  description: string
  original_price: number
  image_url: string
  status: number
}

export interface FlashSale {
  id: number
  product_id: number
  product?: Product
  flash_price: number
  total_stock: number
  available_stock: number
  per_user_limit: number
  start_time: string
  end_time: string
  status: number // 0-Pending 1-Active 2-Finished
}

export interface FlashSaleDetail {
  flash_sale: FlashSale
  current_stock: number
  server_time: string
}

export interface Order {
  id: number
  order_no: string
  user_id: number
  flash_sale_id: number
  flash_sale?: FlashSale
  amount: number
  quantity: number
  status: number // 0-Pending 1-Paid 2-Cancelled 3-Refunded
  created_at: string
  paid_at?: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: User
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface ProductListResponse {
  products: Product[]
  total: number
  page: number
  page_size: number
}

export interface FlashSaleListResponse {
  flash_sales: FlashSale[]
  total: number
  page: number
  page_size: number
}

export interface OrderListResponse {
  orders: Order[]
  total: number
  page: number
  page_size: number
}

export interface RushResponse {
  success: boolean
  ticket: string
  position?: number
  message: string
  order_no?: string
}

export interface ChatMessage {
  role: 'user' | 'assistant' | 'system'
  content: string
}

export interface ChatResponse {
  session_id: string
  response: string
  related_data?: any
}