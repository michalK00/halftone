import {Button} from "@/components/ui/button";
import {LoadingSpinner} from "@/components/ui/loading-spinner";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {Pencil} from "lucide-react";
import {Link} from "react-router-dom";
import {Collection} from "../api/collections";
import {DeleteCollectionDialog} from "./delete-collection-dialog";
import {ShareCollectionDialog} from "@/features/collections/components/share-collections-dialog.tsx";

type CollectionsTableProps = {
    collections: Collection[];
    galleryCounts: Array<{
        status: 'pending' | 'error' | 'success';
        data?: number;
    }>;
    onEdit: (collection: Collection) => void;
};

export function CollectionsTable({ collections, galleryCounts, onEdit }: CollectionsTableProps) {
    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Last Updated</TableHead>
                    <TableHead>Galleries No.</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {collections.map((collection, index) => (
                    <TableRow key={collection.id}>
                        <TableCell className="font-medium">
                            <Link
                                to={`/collections/${collection.id}`}
                                className="text-blue-500 hover:underline"
                            >
                                {collection.name}
                            </Link>
                        </TableCell>
                        <TableCell>
                            {new Date(collection.createdAt).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                            {new Date(collection.updatedAt).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                            {galleryCounts[index].status === 'pending' ? (
                                <LoadingSpinner className="w-4 h-4"/>
                            ) : galleryCounts[index].status === 'error' ? (
                                'Error'
                            ) : (
                                galleryCounts[index].data
                            )}
                        </TableCell>
                        <TableCell className="text-right">
                            <ShareCollectionDialog
                                collection={collection}
                                onOpenChange={() => {}}
                            />
                            <Button
                                className="ml-2"
                                variant="outline"
                                size="icon"
                                onClick={() => onEdit(collection)}
                            >
                                <Pencil/>
                            </Button>
                            <DeleteCollectionDialog
                                collection={collection}
                                onDelete={() => {}}
                            />

                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    );
}