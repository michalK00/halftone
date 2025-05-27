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
import {Button} from "@/components/ui/button";
import {LoadingSpinner} from "@/components/ui/loading-spinner";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {Trash2} from "lucide-react";
import {useState} from "react";
import {Collection, deleteCollection} from "../../../api/collections";
import {useToast} from "@/hooks/use-toast.ts";

type DeleteCollectionDialogProps = {
    collection: Collection;
    onDelete: () => void;
};

export function DeleteCollectionDialog({ collection, onDelete }: DeleteCollectionDialogProps) {
    const [isOpen, setIsOpen] = useState(false);
    const { toast } = useToast();
    const queryClient = useQueryClient();

    const deleteMutation = useMutation({
        mutationFn: () => deleteCollection(collection.id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['collections'] });
            setIsOpen(false);
            toast({
                title: "Success",
                description: "Collection deleted successfully",
            });
            onDelete();
        },
        onError: () => {
            toast({
                title: "Error",
                description: "Failed to delete collection",
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
                <Trash2/>
            </Button>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                    <AlertDialogDescription>
                        This will permanently delete the collection "{collection.name}" and all its galleries and photos.
                        This action cannot be undone.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel disabled={deleteMutation.isPending}>Cancel</AlertDialogCancel>
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