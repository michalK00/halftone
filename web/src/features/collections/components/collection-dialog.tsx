import {Button} from "@/components/ui/button";
import {Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog";
import {Input} from "@/components/ui/input";
import {LoadingSpinner} from "@/components/ui/loading-spinner";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useState} from "react";
import {Collection, createCollection, updateCollection} from "../api/collections";
import {useToast} from "@/hooks/use-toast.ts";

type CollectionDialogProps = {
    mode: 'create' | 'edit';
    collection?: Collection;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    trigger?: React.ReactNode;
};

export function CollectionDialog({ mode, collection, open, onOpenChange, trigger }: CollectionDialogProps) {
    const [name, setName] = useState(collection?.name ?? '');
    const { toast } = useToast();
    const queryClient = useQueryClient();

    const mutation = useMutation({
        mutationFn: async (name: string) => {
            if (mode === 'create') {
                return createCollection(name);
            } else {
                return updateCollection(collection!.id, name);
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['collections'] });
            onOpenChange(false);
            setName('');
            toast({
                title: "Success",
                description: `Collection ${mode === 'create' ? 'created' : 'updated'} successfully`,
            });
        },
        onError: () => {
            toast({
                title: "Error",
                description: `Failed to ${mode} collection`,
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
            {trigger && <DialogTrigger asChild>{trigger}</DialogTrigger>}
            <DialogContent>
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>
                            {mode === 'create' ? 'Create New Collection' : 'Edit Collection'}
                        </DialogTitle>
                        <DialogDescription>
                            {mode === 'create'
                                ? 'Enter a name for your new collection.'
                                : 'Update the collection name.'}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="py-4">
                        <Input
                            placeholder="Collection name"
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