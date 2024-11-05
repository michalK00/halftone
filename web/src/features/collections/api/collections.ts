export type Collection = {
    id: string
    name: string
    createdAt: string
    updatedAt: string
}

export async function getCollections(): Promise<Array<Collection>> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/collections`);
    return response.json();
}

export async function getCollection(collectionId: string): Promise<Collection> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/collections/${collectionId}`)
    return response.json()
}

export async function getGalleryCount(collectionId: string): Promise<number> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/collections/${collectionId}/galleryCount`);
    return response.json();
}

export async function createCollection(name: string): Promise<Collection> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/collections`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name }),
    });
    if (!response.ok) {
        throw new Error('Failed to create collection');
    }
    return response.json();
}

export async function updateCollection(id: string, name: string): Promise<Collection> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/collections/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name }),
    });
    if (!response.ok) {
        throw new Error('Failed to update collection');
    }
    return response.json();
}

export async function deleteCollection(id: string): Promise<void> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/collections/${id}`, {
        method: 'DELETE',
    });
    if (!response.ok) {
        throw new Error('Failed to delete collection');
    }
}

