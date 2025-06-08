import api from '@/lib/api';

export type Gallery = {
    id: string
    collectionId: string
    name: string
    createdAt: string
    updatedAt: string
    sharing: Sharing
}

export type Sharing = {
    sharingEnabled: boolean
    accessToken: string
    sharingExpiryDate: string
    sharingUrl: string
}

/**
 * Get all galleries for a collection
 */
export async function getGalleries(collectionId: string): Promise<Gallery[]> {
    try {
        const response = await api.get(`/api/v1/collections/${collectionId}/galleries`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to fetch galleries');
    }
}

/**
 * Get a specific gallery by ID
 */
export async function getGallery(galleryId: string): Promise<Gallery> {
    try {
        const response = await api.get(`/api/v1/galleries/${galleryId}`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to fetch gallery');
    }
}

/**
 * Create a new gallery in a collection
 */
export async function createGallery(collectionId: string, data: { name: string }): Promise<Gallery> {
    try {
        const response = await api.post(`/api/v1/collections/${collectionId}/galleries`, data);
        return response.data;
    } catch (error) {
        throw new Error('Failed to create gallery');
    }
}

/**
 * Update an existing gallery
 */
export async function updateGallery(galleryId: string, data: { name: string }): Promise<Gallery> {
    try {
        const response = await api.put(`/api/v1/galleries/${galleryId}`, data);
        return response.data;
    } catch (error) {
        throw new Error('Failed to update gallery');
    }
}

/**
 * Delete a gallery
 */
export async function deleteGallery(galleryId: string): Promise<void> {
    try {
        await api.delete(`/api/v1/galleries/${galleryId}`);
    } catch (error) {
        throw new Error('Failed to delete gallery');
    }
}