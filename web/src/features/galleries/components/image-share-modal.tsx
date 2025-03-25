import { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import {CalendarIcon, Share2, Copy, Check, CalendarCog} from "lucide-react";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import {useToast} from "@/hooks/use-toast.ts";
import {shareGallery, stopGallerySharing} from "@/api/share.ts";
import {useQueryClient} from "@tanstack/react-query";

type ImageShareModalProps = {
    galleryId: string;
    sharingEnabled: boolean;
    sharingExpiryDate?: Date;
    sharingUrl?: string;
};

function ImageShareModal({
                             galleryId,
                             sharingEnabled,
                             sharingExpiryDate,
                             sharingUrl,
                         }: ImageShareModalProps) {
    const [isOpen, setIsOpen] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const displayDate = sharingEnabled ? sharingExpiryDate : undefined
    const [endDate, setEndDate] = useState<Date | undefined>(displayDate);
    const [copied, setCopied] = useState(false);

    const { toast } = useToast();
    const queryClient = useQueryClient();

    const closeModal = () => {
        setIsOpen(false);
    };

    const handleShare = async () => {
        if (!endDate) return;


        try {
            setIsLoading(true);
            if (sharingEnabled) {
                const response = await stopGallerySharing(galleryId)
                sharingEnabled = response.sharing.sharingEnabled

                toast({
                    title: "Success",
                    description: `Gallery sharing stopped successfully`,
                });
            } else {
                const response = await shareGallery(galleryId, { sharingExpiry: endDate })
                sharingEnabled = true
                sharingUrl = response.shareUrl
                sharingExpiryDate = response.sharingExpiry
                toast({
                    title: "Success",
                    description: `Gallery shared successfully`,
                });
            }
            queryClient.invalidateQueries({ queryKey: ['gallery', galleryId] });

        } catch (error) {
            toast({
                title: "Error",
                description: `Failed to ${sharingEnabled ? "stop sharing" : "share"} gallery`,
                variant: "destructive",
            });
            console.error(`Failed to ${sharingEnabled ? "stop sharing" : "share"} gallery:`, error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleCopy = async () => {
        try {
            await navigator.clipboard.writeText(sharingUrl!);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        } catch (err) {
            console.error('Failed to copy text: ', err);
        }
    };

    // Calculate date constraints
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const maxDate = new Date();
    maxDate.setFullYear(maxDate.getFullYear() + 1);

    return (
        <div>
            <Button
                variant="outline"
                size="icon"
                onClick={() => setIsOpen(true)}
                className="hover:bg-accent hover:text-accent-foreground transition-colors"
                title="Modify sharing settings"
            >
                <Share2 className="w-4 h-4" />
            </Button>

            <Dialog open={isOpen} onOpenChange={closeModal}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle className="text-xl font-semibold">Share Gallery</DialogTitle>
                    </DialogHeader>

                    <div className="space-y-6">
                        <div className="space-y-4">
                            <div className="text-sm text-muted-foreground">
                                {sharingEnabled ? (
                                    <div className="flex items-center gap-2">
                                        <div className="w-2 h-2 bg-green-500 rounded-full" />
                                        <span>Sharing until {sharingExpiryDate?.toLocaleDateString()}</span>
                                    </div>
                                ) : (
                                    <div className="flex items-center gap-2">
                                        <div className="w-2 h-2 bg-gray-300 rounded-full" />
                                        <span>Currently not sharing</span>
                                    </div>
                                )}
                            </div>

                            <div className="flex flex-col space-y-2">
                                <label className="text-sm font-medium">Expiry Date</label>
                                <Popover>
                                    <PopoverTrigger asChild>
                                        <Button
                                            variant="outline"
                                            className={cn(
                                                "w-full pl-3 text-left font-normal",
                                                !endDate && "text-muted-foreground"
                                            )}
                                        >
                                            {endDate ? (
                                                endDate.toLocaleDateString()
                                            ) : (
                                                <span>Select expiry date</span>
                                            )}
                                            <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                        </Button>
                                    </PopoverTrigger>
                                    <PopoverContent className="w-auto p-0" align="start">
                                        <Calendar
                                            mode="single"
                                            selected={endDate}
                                            onSelect={setEndDate}
                                            disabled={(date) => date <= today || date > maxDate}
                                            initialFocus
                                        />
                                    </PopoverContent>
                                </Popover>
                                {sharingEnabled && (endDate?.toDateString() != sharingExpiryDate?.toDateString()) && <Button variant="outline">Reschedule expiry date <CalendarCog/></Button>}
                            </div>

                        </div>

                        {sharingEnabled && (
                            <Card className="bg-accent/50">
                                <CardContent className="p-4 space-y-4">
                                    <div className="flex justify-center bg-white rounded-lg p-2">
                                        <img
                                            src={`${import.meta.env.VITE_BACKEND_URL}/api/v1/qr?url=${sharingUrl}`}
                                            alt="QR Code"
                                            className="w-48 h-48 object-contain"
                                        />
                                    </div>
                                    <div className="flex gap-2">
                                        <Input
                                            value={sharingUrl}
                                            readOnly
                                            className="text-sm bg-background"
                                        />
                                        <Button
                                            variant="outline"
                                            size="icon"
                                            onClick={handleCopy}
                                            className="shrink-0 bg-background hover:bg-accent"
                                        >
                                            {copied ? (
                                                <Check className="w-4 h-4 text-green-500" />
                                            ) : (
                                                <Copy className="w-4 h-4" />
                                            )}
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>
                        )}
                    </div>

                    <DialogFooter>
                        <div className="flex gap-2 justify-end">
                            <Button
                                variant="outline"
                                onClick={closeModal}
                                disabled={isLoading}
                            >
                                Close
                            </Button>
                            <Button
                                onClick={handleShare}
                                disabled={!endDate || isLoading}
                                className="bg-primary hover:bg-primary/90"
                            >
                                {sharingEnabled ? "Stop sharing" : "Share"}
                            </Button>
                        </div>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}

export default ImageShareModal;