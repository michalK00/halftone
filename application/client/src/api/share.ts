import api from '@/lib/api';
import { Gallery } from '@/api/galleries';

export interface ShareGalleryRequestBody {
    sharingExpiry: Date;
}

export interface ShareGalleryResponse {
    galleryId: string;
    accessToken: string;
    sharingExpiry: Date;
    shareUrl: string;
}

/**
 * Enable sharing for a gallery
 */
export async function shareGallery(galleryId: string, body: ShareGalleryRequestBody): Promise<ShareGalleryResponse> {
    try {
        const response = await api.post(`/api/v1/galleries/${galleryId}/sharing/share`, body);
        return response.data;
    } catch (error) {
        throw new Error('Failed to share gallery');
    }
}

/**
 * Update the expiry date for a shared gallery
 */
export async function rescheduleGallerySharing(galleryId: string, body: ShareGalleryRequestBody): Promise<ShareGalleryResponse> {
    try {
        const response = await api.put(`/api/v1/galleries/${galleryId}/sharing/reschedule`, body);
        return response.data;
    } catch (error) {
        throw new Error('Failed to reschedule gallery sharing');
    }
}

/**
 * Disable sharing for a gallery
 */
export async function stopGallerySharing(galleryId: string): Promise<Gallery> {
    try {
        const response = await api.put(`/api/v1/galleries/${galleryId}/sharing/stop`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to stop gallery sharing');
    }
}