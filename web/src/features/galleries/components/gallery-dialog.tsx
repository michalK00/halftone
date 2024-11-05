import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {useEffect, useState} from "react";
import {createGallery, Gallery, updateGallery} from "../api/galleries";
import {useToast} from "@/hooks/use-toast.ts";

type GalleryDialogProps = {
    mode: 'create' | 'edit';
    collectionId?: string;
    gallery?: Gallery;
    open: boolean;
    onOpenChange: (open: boolean) => void;
};

export function GalleryDialog({ mode, collectionId, gallery, open, onOpenChange }: GalleryDialogProps) {
    const [name, setName] = useState('');
    const { toast } = useToast();
    const queryClient = useQueryClient();

    useEffect(() => {
        if (open && mode === 'edit' && gallery) {
            setName(gallery.name);
        } else if (!open) {
            setName('');
        }
    }, [open, mode, gallery]);

    const mutation = useMutation({
        mutationFn: async (name: string) => {
            if (mode === 'create') {
                return createGallery(collectionId!, { name });
            } else {
                return updateGallery(gallery!.id, { name });
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['galleries', collectionId] });
            onOpenChange(false);
            setName('');
            toast({
                title: "Success",
                description: `Gallery ${mode === 'create' ? 'created' : 'updated'} successfully`,
            });
        },
        onError: () => {
            toast({
                title: "Error",
                description: `Failed to ${mode} gallery`,
                variant: "destructive",
            });
        },
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!name.trim()) return;
        mutation.mutate(name);
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>
                            {mode === 'create' ? 'Create New Gallery' : 'Edit Gallery'}
                        </DialogTitle>
                        <DialogDescription>
                            {mode === 'create'
                                ? 'Enter a name for your new gallery.'
                                : 'Update the gallery name.'}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="py-4">
                        <Input
                            placeholder="Gallery name"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            disabled={mutation.isPending}
                        />
                    </div>
                    <DialogFooter>
                        <Button
                            variant="outline"
                            onClick={() => onOpenChange(false)}
                            type="button"
                            disabled={mutation.isPending}
                        >
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            disabled={mutation.isPending || !name.trim()}
                        >
                            {mutation.isPending ? (
                                <>
                                    <LoadingSpinner className="w-4 h-4 mr-2" />
                                    {mode === 'create' ? 'Creating...' : 'Updating...'}
                                </>
                            ) : (
                                mode === 'create' ? 'Create' : 'Update'
                            )}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}