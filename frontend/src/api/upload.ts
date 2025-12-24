import http from '@/utils/http'
import type { ApiResponse } from '@/types'

interface UploadResponse {
    url: string
}

export const uploadImage = (file: File): Promise<ApiResponse<UploadResponse>> => {
    const formData = new FormData()
    formData.append('file', file)
    return http.post('/admin/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
    })
}
