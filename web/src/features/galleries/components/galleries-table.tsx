import { LoadingSpinner } from "@/components/ui/loading-spinner.tsx";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "react-router-dom";
import { Pencil, Lock, Share2 } from "lucide-react";
import {Gallery} from "@/features/galleries/api/galleries.ts";
import {DeleteGalleryDialog} from "@/features/galleries/components/delete-gallery-dialog.tsx";

type GalleriesTableProps = {
    galleries: Gallery[];
    isLoading: boolean;
    error: Error | null;
    onEdit: (gallery: Gallery) => void;
};

export function GalleriesTable({ galleries, isLoading, error, onEdit }: GalleriesTableProps) {
    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-32">
                <LoadingSpinner />
            </div>
        );
    }

    if (error) {
        return (
            <div className="text-center text-red-500">
                Error loading galleries
            </div>
        );
    }

    if (galleries.length === 0) {
        return (
            <div className="text-center text-gray-500">
                No galleries found
            </div>
        );
    }

    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Last Updated</TableHead>
                    <TableHead>Sharing Status</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {galleries.map((gallery) => (
                    <TableRow key={gallery.id}>
                        <TableCell className="font-medium">
                            <Link
                                to={`/galleries/${gallery.id}`}
                                className="text-blue-500 hover:underline flex items-center gap-2"
                            >
                                {gallery.name}
                            </Link>
                        </TableCell>
                        <TableCell>
                            {new Date(gallery.createdAt).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                            {new Date(gallery.updatedAt).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                            <div className="flex items-center gap-2">
                                {gallery.sharing.sharingEnabled ? (
                                    <>
                                        <Share2 className="w-4 h-4 text-green-500" />
                                        <span className="text-sm">
                                            Until {new Date(gallery.sharing.sharingExpiryDate).toLocaleDateString()}
                                        </span>
                                    </>
                                ) : (
                                    <>
                                    <Lock className="w-4 h-4" />
                                    <span className="text-sm">
                                    Disabled
                            </span>
                                    </>
                                )}
                            </div>
                        </TableCell>
                        <TableCell className="text-right flex flex-col sm:flex-row gap-2 justify-end">
                            <Button
                                variant="outline"
                                size="icon"
                            >
                                {gallery.sharing.sharingEnabled ? (
                                    <Lock className="w-4 h-4" />
                                ) : (
                                    <Share2 className="w-4 h-4" />
                                )}
                            </Button>
                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => onEdit(gallery)}
                            >
                                <Pencil className="w-4 h-4" />
                            </Button>
                            <DeleteGalleryDialog gallery={gallery} />
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    );
}