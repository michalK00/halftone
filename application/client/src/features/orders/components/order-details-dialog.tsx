import { useState } from "react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Eye, Loader2 } from "lucide-react";
import { Order, OrderStatus, updateOrder } from "@/api/orders";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "@/hooks/use-toast";

interface OrderDetailsDialogProps {
    order: Order;
}

export function OrderDetailsDialog({ order }: OrderDetailsDialogProps) {
    const [open, setOpen] = useState(false);
    const [selectedStatus, setSelectedStatus] = useState<OrderStatus | null>(null);
    const queryClient = useQueryClient();
    const { toast } = useToast();

    // Reset selected status when dialog opens/closes or order changes
    const currentStatus = selectedStatus ?? order.status;

    const updateStatusMutation = useMutation({
        mutationFn: (status: OrderStatus) => updateOrder(order.id, { status }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['orders'] });
            toast({
                title: "Order updated",
                description: "Order status has been updated successfully.",
            });
            setSelectedStatus(null);
            setOpen(false);
        },
        onError: () => {
            toast({
                title: "Error",
                description: "Failed to update order status. Please try again.",
                variant: "destructive",
            });
        },
    });

    const handleStatusUpdate = () => {
        if (currentStatus !== order.status) {
            updateStatusMutation.mutate(currentStatus);
        }
    };

    const handleOpenChange = (newOpen: boolean) => {
        setOpen(newOpen);
        if (!newOpen) {
            setSelectedStatus(null);
        }
    };

    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <DialogTrigger asChild>
                <Button variant="ghost" size="sm">
                    <Eye className="h-4 w-4" />
                </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl">
                <DialogHeader>
                    <DialogTitle>Order Details</DialogTitle>
                    <DialogDescription>
                        View and manage order information
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <Label className="text-sm text-muted-foreground">Order ID</Label>
                            <p className="text-sm font-mono">{order.id}</p>
                        </div>
                        <div>
                            <Label className="text-sm text-muted-foreground">Gallery ID</Label>
                            <p className="text-sm font-mono">{order.galleryId}</p>
                        </div>
                    </div>

                    <Separator />

                    <div className="space-y-2">
                        <Label className="text-sm text-muted-foreground">Client Email</Label>
                        <p className="text-sm">{order.clientEmail}</p>
                    </div>

                    <div className="space-y-2">
                        <Label className="text-sm text-muted-foreground">Comment</Label>
                        <p className="text-sm">{order.comment || "No comment provided"}</p>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <Label className="text-sm text-muted-foreground">Created</Label>
                            <p className="text-sm">{new Date(order.createdAt).toLocaleString()}</p>
                        </div>
                        <div>
                            <Label className="text-sm text-muted-foreground">Last Updated</Label>
                            <p className="text-sm">{new Date(order.updatedAt).toLocaleString()}</p>
                        </div>
                    </div>

                    <div className="space-y-2">
                        <Label className="text-sm text-muted-foreground">Photos ({order.photos.length})</Label>
                        <div className="max-h-32 overflow-y-auto">
                            <div className="text-sm font-mono space-y-1">
                                {order.photos.map((photo, index) => (
                                    <div key={photo.photoId} className="text-muted-foreground">
                                        {index + 1}. {photo.photoId}
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>

                    <Separator />

                    <div className="space-y-2">
                        <Label htmlFor="status-select">Status</Label>
                        <div className="flex items-center gap-2">
                            <select
                                id="status-select"
                                value={currentStatus}
                                onChange={(e) => setSelectedStatus(e.target.value as OrderStatus)}
                                className="flex h-10 w-[180px] items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                            >
                                <option value="pending">Pending</option>
                                <option value="completed">Completed</option>
                            </select>
                            <Button
                                onClick={handleStatusUpdate}
                                disabled={currentStatus === order.status || updateStatusMutation.isPending}
                            >
                                {updateStatusMutation.isPending && (
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                )}
                                Update Status
                            </Button>
                        </div>
                    </div>

                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        Current status:
                        <Badge variant={order.status === 'completed' ? 'default' : 'secondary'}>
                            {order.status}
                        </Badge>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}