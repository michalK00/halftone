import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { PlusCircle } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import {Gallery, getGalleries} from "../../../api/galleries.ts";
import { GalleriesTable } from "./galleries-table.tsx";
import {GalleryDialog} from "@/features/galleries/components/gallery-dialog.tsx";
import {useState} from "react";


export function GalleriesSection({ collectionId }: { collectionId: string }) {
    const [editGallery, setEditGallery] = useState<Gallery | null>(null);
    const [isDialogOpen, setIsDialogOpen] = useState(false);

    const { data, isLoading, error } = useQuery({
        queryKey: ['galleries', collectionId],
        queryFn: () => getGalleries(collectionId),
    });

    const handleEdit = (gallery: Gallery) => {
        setEditGallery(gallery);
        setIsDialogOpen(true);
    };

    return (
        <Card>
            <CardHeader>
                <div className="flex justify-between items-center">
                    <div>
                        <CardTitle>Galleries</CardTitle>
                        <CardDescription>Manage photo galleries in this collection</CardDescription>
                    </div>
                    <Button
                        className="flex items-center gap-2"
                        onClick={() => {
                            setEditGallery(null);
                            setIsDialogOpen(true);
                        }}
                    >
                        <PlusCircle className="w-4 h-4" />
                        New Gallery
                    </Button>
                </div>
            </CardHeader>
            <CardContent>
                <GalleriesTable
                    galleries={data ?? []}
                    isLoading={isLoading}
                    error={error as Error}
                    onEdit={handleEdit}
                />
            </CardContent>

            <GalleryDialog
                mode={editGallery ? 'edit' : 'create'}
                collectionId={collectionId}
                gallery={editGallery ?? undefined}
                open={isDialogOpen}
                onOpenChange={(open) => {
                    setIsDialogOpen(open);
                    if (!open) setEditGallery(null);
                }}
            />
        </Card>
    );
}