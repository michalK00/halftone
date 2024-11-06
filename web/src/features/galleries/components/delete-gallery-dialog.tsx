import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Trash2 } from "lucide-react";
import { useState } from "react";
import {deleteGallery, Gallery} from "../api/galleries";
import {useToast} from "@/hooks/use-toast.ts";

type DeleteGalleryDialogProps = {
    gallery: Gallery;
    onDelete?: () => void;
};

export function DeleteGalleryDialog({ gallery, onDelete }: DeleteGalleryDialogProps) {
    const [isOpen, setIsOpen] = useState(false);
    const { toast } = useToast();
    const queryClient = useQueryClient();

    const deleteMutation = useMutation({
        mutationFn: () => deleteGallery(gallery.id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['galleries', gallery.collectionId] });
            setIsOpen(false);
            toast({
                title: "Success",
                description: "Gallery deleted successfully",
            });
            if (onDelete) {
                onDelete()
            }
        },
        onError: () => {
            toast({
                title: "Error",
                description: "Failed to delete gallery",
                variant: "destructive",
            });
        },
    });

    return (
        <AlertDialog open={isOpen} onOpenChange={setIsOpen}>
            <Button
                variant="destructive"
                size="icon"
                onClick={() => setIsOpen(true)}
            >
                <Trash2 className="w-4 h-4" />
            </Button>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                    <AlertDialogDescription>
                        This will permanently delete the gallery "{gallery.name}" and all its photos.
                        This action cannot be undone.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel disabled={deleteMutation.isPending}>
                        Cancel
                    </AlertDialogCancel>
                    <AlertDialogAction
                        onClick={(e) => {
                            e.preventDefault();
                            deleteMutation.mutate();
                        }}
                        className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                        disabled={deleteMutation.isPending}
                    >
                        {deleteMutation.isPending ? (
                            <>
                                <LoadingSpinner className="w-4 h-4 mr-2" />
                                Deleting...
                            </>
                        ) : (
                            'Delete'
                        )}
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
}