import api from '@/lib/api.ts';

export type Collection = {
    id: string
    name: string
    createdAt: string
    updatedAt: string
}

export type GalleryCount = {
    count: number
}

export async function getCollections(): Promise<Array<Collection>> {
    const response = await api.get('/api/v1/collections');
    return response.data;
}

export async function getCollection(collectionId: string): Promise<Collection> {
    const response = await api.get(`/api/v1/collections/${collectionId}`);
    return response.data;
}

export async function getGalleryCount(collectionId: string): Promise<GalleryCount> {
    const response = await api.get(`/api/v1/collections/${collectionId}/galleryCount`);
    return response.data;
}

export async function createCollection(name: string): Promise<Collection> {
    try {
        const response = await api.post('/api/v1/collections', { name });
        return response.data;
    } catch {
        throw new Error('Failed to create collection');
    }
}

export async function updateCollection(id: string, name: string): Promise<Collection> {
    try {
        const response = await api.put(`/api/v1/collections/${id}`, { name });
        return response.data;
    } catch {
        throw new Error('Failed to update collection');
    }
}

export async function deleteCollection(id: string): Promise<void> {
    try {
        await api.delete(`/api/v1/collections/${id}`);
    } catch {
        throw new Error('Failed to delete collection');
    }
}