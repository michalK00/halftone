import {getMimeType} from "@/lib/utils.ts";

export interface Photo {
    id: string
    url: string;
    originalFilename: string;
}

export interface UploadRequestBody {
    originalFilename: string,
}

export interface UploadResponseBody {
    id: string,
    originalFilename: string,
    presignedPostRequest: {
        URL: string,
        Values: Record<string, string>,
    }
}

export async function getPhotos(galleryId: string): Promise<Photo[]> {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/galleries/${galleryId}/photos`);
    if (!response.ok) {
        throw new Error('Failed to fetch galleries');
    }
    return response.json();
}

export async function getUploadUrls(galleryId: string, body: UploadRequestBody[]): Promise<UploadResponseBody[]> {
    const response= await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/galleries/${galleryId}/photos`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body)
    })

    if (!response.ok) {
        throw new Error('Failed to get upload URL')
    }
    return response.json();
}

export async function uploadToAWS(file: File, uploadData: UploadResponseBody) {
    const formData = new FormData()

    Object.entries(uploadData.presignedPostRequest.Values).forEach(([key, value]) => {
        formData.append(key, value)
    })
    formData.append("Content-Type", await getMimeType(file))
    formData.append("file", file)

    const response = await fetch(uploadData.presignedPostRequest.URL, {
        method: 'POST',
        body: formData,
    });
    if (!response.ok) {
        throw new Error('Failed to upload to AWS');
    }
}

export async function confirmUpload(photoId: string) {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/photos/${photoId}/confirm`, {
        method: 'PUT',
    });

    if (!response.ok) {
        throw new Error('Failed to confirm upload');
    }
};

export async function deletePhoto(photoId: string) {
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/photos/${photoId}`, {
        method: 'DELETE',
    })
    if (!response.ok) {
        throw new Error('Failed to delete photo')
    }

}