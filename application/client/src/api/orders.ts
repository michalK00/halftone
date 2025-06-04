import api from '@/lib/api';

export type OrderStatus = 'pending' | 'completed';

export type Order = {
    id: string
    galleryId: string
    clientEmail: string
    comment: string
    photos: Photo[]
    status: OrderStatus
    createdAt: string
    updatedAt: string
}

type Photo = {
    photoId: string
}

export type UpdateOrderData = {
    status?: OrderStatus
    comment?: string
}

/**
 * Get all orders for galleries owned by the authenticated user
 */
export async function getOrders(): Promise<Order[]> {
    try {
        const response = await api.get('/api/v1/orders');
        return response.data;
    } catch (error) {
        throw new Error('Failed to fetch orders');
    }
}

/**
 * Get a specific order by ID
 */
export async function getOrder(orderId: string): Promise<Order> {
    try {
        const response = await api.get(`/api/v1/orders/${orderId}`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to fetch order');
    }
}

/**
 * Update an existing order
 */
export async function updateOrder(orderId: string, data: UpdateOrderData): Promise<Order> {
    try {
        const response = await api.put(`/api/v1/orders/${orderId}`, data);
        return response.data;
    } catch (error) {
        throw new Error('Failed to update order');
    }
}

/**
 * Delete an order
 */
export async function deleteOrder(orderId: string): Promise<void> {
    try {
        await api.delete(`/api/v1/orders/${orderId}`);
    } catch (error) {
        throw new Error('Failed to delete order');
    }
}