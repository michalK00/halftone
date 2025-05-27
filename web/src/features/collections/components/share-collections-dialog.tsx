import {Button} from "@/components/ui/button";
import {Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog";
import {Collection} from "../../../api/collections";
import {Share} from "lucide-react";

type CollectionDialogProps = {
    collection?: Collection;
    onOpenChange: (open: boolean) => void;
};

export function ShareCollectionDialog({ onOpenChange }: CollectionDialogProps) {

    return (
        <Dialog onOpenChange={onOpenChange}>
            <DialogTrigger asChild>
                            <Button
                                variant="outline"
                                size="icon"
                            >
                                <Share/>
                            </Button>
                </DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>
                        Collection sharing
                    </DialogTitle>
                    <DialogDescription>
                        Configure collection sharing options
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    Footer
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}