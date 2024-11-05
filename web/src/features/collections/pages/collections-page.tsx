import {SidebarTrigger} from "@/components/ui/sidebar";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
} from "@/components/ui/breadcrumb";
import {Link} from "react-router-dom";
import {ModeToggle} from "@/components/mode-toggle";
import {useQueries, useQuery} from "@tanstack/react-query";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {PlusCircle} from "lucide-react";
import {LoadingSpinner} from "@/components/ui/loading-spinner";
import {useState} from "react";
import {Collection, getCollections, getGalleryCount} from "../api/collections";
import {CollectionDialog} from "../components/collection-dialog";
import {CollectionsTable} from "../components/collections-table";

function CollectionsQuery() {
    const [editCollection, setEditCollection] = useState<Collection | null>(null);
    const [isDialogOpen, setIsDialogOpen] = useState(false);

    const collectionsQuery = useQuery({
        queryKey: ['collections'],
        queryFn: getCollections
    });

    const galleryCountQueries = useQueries({
        queries: (collectionsQuery.data ?? []).map((collection) => ({
            queryKey: ['galleryCount', collection.id],
            queryFn: () => getGalleryCount(collection.id),
            enabled: !!collectionsQuery.data,
        })),
    });

    if (collectionsQuery.status === 'pending') {
        return (
            <Card className="w-full">
                <CardContent className="flex p-6 h-64 justify-center items-center">
                    <LoadingSpinner />
                </CardContent>
            </Card>
        );
    }

    if (collectionsQuery.status === 'error') {
        return (
            <Card className="w-full">
                <CardContent className="p-6">
                    <div className="text-center text-red-500">
                        Error loading collections. Please try again later.
                    </div>
                </CardContent>
            </Card>
        );
    }

    const handleEdit = (collection: Collection) => {
        setEditCollection(collection);
        setIsDialogOpen(true);
    };

    const galleryCounts = galleryCountQueries.map(query => ({
        status: query.status,
        data: query.data,
    }));

    return (
        <Card className="w-full">
            <CardHeader>
                <div className="flex justify-between items-center">
                    <div>
                        <CardTitle>Collections</CardTitle>
                        <CardDescription>
                            Manage your collections and their contents
                        </CardDescription>
                    </div>
                    <CollectionDialog
                        mode="create"
                        open={isDialogOpen && !editCollection}
                        onOpenChange={(open) => setIsDialogOpen(open)}
                        trigger={
                            <Button className="flex items-center gap-2">
                                <PlusCircle className="w-4 h-4" />
                                New Collection
                            </Button>
                        }
                    />
                </div>
            </CardHeader>
            <CardContent>
                <CollectionsTable
                    collections={collectionsQuery.data}
                    galleryCounts={galleryCounts}
                    onEdit={handleEdit}
                />
            </CardContent>

            {editCollection && (
                <CollectionDialog
                    mode="edit"
                    collection={editCollection}
                    open={isDialogOpen}
                    onOpenChange={(open) => {
                        setIsDialogOpen(open);
                        if (!open) setEditCollection(null);
                    }}
                />
            )}
        </Card>
    );
}

export default function CollectionsPage() {
    return (
        <main className="w-full">
            <div className="w-full flex p-2 justify-between items-center">
                <div className="flex gap-4 items-center">
                    <SidebarTrigger/>
                    <Breadcrumb>
                        <BreadcrumbList>
                            <BreadcrumbItem>
                                <BreadcrumbLink asChild>
                                    <Link to="/collections">Collections</Link>
                                </BreadcrumbLink>
                            </BreadcrumbItem>
                        </BreadcrumbList>
                    </Breadcrumb>
                </div>
                <ModeToggle/>
            </div>
            <div className="flex p-2">
                <CollectionsQuery />
            </div>
        </main>
    );
}