import {Gallery} from "@/features/galleries/api/galleries.ts";

export interface ShareGalleryRequestBody {
    sharingExpiry: Date,
}

export interface ShareGalleryResponse {
    galleryId: string,
    accessToken: string,
    sharingExpiry: Date,
    shareUrl: string,
}

export async function shareGallery(galleryId: string, body: ShareGalleryRequestBody): Promise<ShareGalleryResponse> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/galleries/${galleryId}/sharing/share`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body)
    });
    if (!response.ok) {
        throw new Error('Failed to get upload URL')
    }
    return response.json();
}

export async function rescheduleGallerySharing(galleryId: string, body: ShareGalleryRequestBody): Promise<ShareGalleryResponse> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/galleries/${galleryId}/sharing/reschedule`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body)
    });
    if (!response.ok) {
        throw new Error('Failed to get upload URL')
    }
    return response.json();
}

export async function stopGallerySharing(galleryId: string): Promise<Gallery> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/galleries/${galleryId}/sharing/stop`, {
        method: 'PUT',
    });
    if (!response.ok) {
        throw new Error('Failed to get upload URL')
    }
    return response.json();
}