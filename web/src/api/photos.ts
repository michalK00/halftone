import api from '@/lib/api';
import axios from 'axios';
import { getMimeType } from '@/lib/utils';

export interface Photo {
    id: string;
    url: string;
    originalFilename: string;
}

export interface UploadRequestBody {
    originalFilename: string;
}

export interface UploadResponseBody {
    id: string;
    originalFilename: string;
    presignedPostRequest: {
        URL: string;
        Values: Record<string, string>;
    };
}

/**
 * Get all photos for a gallery
 */
export async function getPhotos(galleryId: string): Promise<Photo[]> {
    try {
        const response = await api.get(`/api/v1/galleries/${galleryId}/photos`);
        return response.data;
    } catch (error) {
        throw new Error('Failed to fetch photos');
    }
}

/**
 * Get pre-signed URLs for uploading photos
 */
export async function getUploadUrls(galleryId: string, body: UploadRequestBody[]): Promise<UploadResponseBody[]> {
    try {
        const response = await api.post(`/api/v1/galleries/${galleryId}/photos`, body);
        return response.data;
    } catch (error) {
        throw new Error('Failed to get upload URLs');
    }
}

/**
 * Upload a file to AWS S3 using pre-signed URL
 * Note: We use axios directly here instead of our API client
 * because this request goes to AWS S3, not our backend API
 */
export async function uploadToAWS(file: File, uploadData: UploadResponseBody): Promise<void> {
    try {
        const formData = new FormData();

        // Add all fields from the pre-signed request
        Object.entries(uploadData.presignedPostRequest.Values).forEach(([key, value]) => {
            formData.append(key, value);
        });

        // Add content type and file
        const contentType = await getMimeType(file);
        formData.append('Content-Type', contentType);
        formData.append('file', file);

        // Use axios directly for the S3 upload
        await axios.post(uploadData.presignedPostRequest.URL, formData, {
            // Don't add auth headers for S3 uploads
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });
    } catch (error) {
        throw new Error('Failed to upload to AWS');
    }
}

/**
 * Confirm that a photo was successfully uploaded
 */
export async function confirmUpload(photoId: string): Promise<void> {
    try {
        await api.put(`/api/v1/photos/${photoId}/confirm`);
    } catch (error) {
        throw new Error('Failed to confirm upload');
    }
}

/**
 * Delete a photo
 */
export async function deletePhoto(photoId: string): Promise<void> {
    try {
        await api.delete(`/api/v1/photos/${photoId}`);
    } catch (error) {
        throw new Error('Failed to delete photo');
    }
}