import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
    getCollections,
    getCollection,
    createCollection,
    updateCollection,
    deleteCollection,
    getGalleryCount,
    type Collection,
} from '@/api';

/**
 * Hook to fetch all collections
 */
export function useCollections() {
    return useQuery({
        queryKey: ['collections'],
        queryFn: () => getCollections(),
    });
}

/**
 * Hook to fetch a single collection by ID
 */
export function useCollection(collectionId: string) {
    return useQuery({
        queryKey: ['collection', collectionId],
        queryFn: () => getCollection(collectionId),
        enabled: !!collectionId,
    });
}

/**
 * Hook to fetch gallery count for a collection
 */
export function useGalleryCount(collectionId: string) {
    return useQuery(getGalleryCountOptions(collectionId));
}

export function getGalleryCountOptions(collectionId: string) {
    return {
        queryKey: ['collection', collectionId, 'galleryCount'],
        queryFn: () => getGalleryCount(collectionId),
        enabled: !!collectionId,
    };
}
/**
 * Hook to create a new collection
 */
export function useCreateCollection() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (name: string) => createCollection(name),
        onSuccess: (newCollection: Collection) => {
            // Update collections list
            queryClient.setQueryData<Collection[]>(['collections'], (oldData = []) => [
                ...oldData,
                newCollection,
            ]);

            // Invalidate collections query to refetch
            queryClient.invalidateQueries({
                queryKey: ['collections'],
            });
        },
    });
}

/**
 * Hook to update a collection
 */
export function useUpdateCollection() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, name }: { id: string; name: string }) => updateCollection(id, name),
        onSuccess: (updatedCollection: Collection) => {
            // Update collection in cache
            queryClient.setQueryData(['collection', updatedCollection.id], updatedCollection);

            // Update collections list
            queryClient.setQueryData<Collection[]>(['collections'], (oldData = []) =>
                oldData.map(collection =>
                    collection.id === updatedCollection.id ? updatedCollection : collection
                )
            );
        },
    });
}

/**
 * Hook to delete a collection
 */
export function useDeleteCollection() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (id: string) => deleteCollection(id),
        onSuccess: (_, id) => {
            // Remove collection from cache
            queryClient.removeQueries({
                queryKey: ['collection', id],
            });

            // Update collections list
            queryClient.setQueryData<Collection[]>(['collections'], (oldData = []) =>
                oldData.filter(collection => collection.id !== id)
            );

            // Invalidate collections query to refetch
            queryClient.invalidateQueries({
                queryKey: ['collections'],
            });
        },
    });
}
