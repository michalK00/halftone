import { Button } from "@/components/ui/button";
import {
    Dialog, DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle, DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { Gallery, updateGallery } from "../../../api/galleries";
import { useToast } from "@/hooks/use-toast.ts";
import { Pencil } from "lucide-react";

type GalleryDialogProps = {
    gallery: Gallery;
};

export function EditGalleryDialog({ gallery }: GalleryDialogProps) {
    const [open, setOpen] = useState(false);
    const [name, setName] = useState(gallery.name);
    const { toast } = useToast();
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: async (name: string) => {
            return updateGallery(gallery.id, { name });
        },
        onSuccess: (updatedGallery) => {
            queryClient.setQueryData(['gallery', gallery.id], updatedGallery);

            toast({
                title: "Success",
                description: `Gallery updated successfully`,
            });
            setOpen(false);
        },
        onError: () => {
            toast({
                title: "Error",
                description: `Failed to update gallery`,
                variant: "destructive",
            });
        }
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!name.trim()) return;
        mutation.mutate(name);
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button variant="outline" size="icon">
                    <Pencil className="w-4 h-4"/>
                </Button>
            </DialogTrigger>
            <DialogContent>
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>
                            Edit Gallery
                        </DialogTitle>
                        <DialogDescription>
                            Update the gallery name.
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
                        <DialogClose asChild>
                            <Button variant="outline" disabled={mutation.isPending}>
                                Cancel
                            </Button>
                        </DialogClose>
                        <Button
                            type="submit"
                            disabled={mutation.isPending || !name.trim()}
                        >
                            {mutation.isPending ? (
                                <>
                                    <LoadingSpinner className="w-4 h-4 mr-2" />
                                    Updating...
                                </>
                            ) : (
                                <>Update</>
                            )}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}