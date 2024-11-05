export type Gallery = {
    id: string
    collectionId: string
    name: string
    createdAt: string
    updatedAt: string
    sharingEnabled: boolean
    sharingExpiryDate: string
}

export async function getGalleries(collectionId: string): Promise<Gallery[]> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/collections/${collectionId}/galleries`);
    if (!response.ok) {
        throw new Error('Failed to fetch galleries');
    }
    return response.json();
}

export async function createGallery(collectionId: string, data: { name: string }): Promise<Gallery> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/collections/${collectionId}/galleries`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });
    if (!response.ok) {
        throw new Error('Failed to create gallery');
    }
    return response.json();
}

export async function updateGallery(galleryId: string, data: { name: string }): Promise<Gallery> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/galleries/${galleryId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });
    if (!response.ok) {
        throw new Error('Failed to update gallery');
    }
    return response.json();
}

export async function deleteGallery(galleryId: string): Promise<void> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/galleries/${galleryId}`, {
        method: 'DELETE',
    });
    if (!response.ok) {
        throw new Error('Failed to delete gallery');
    }
}