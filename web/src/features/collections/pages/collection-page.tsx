import { SidebarTrigger } from "@/components/ui/sidebar";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbSeparator
} from "@/components/ui/breadcrumb";
import {Link, useNavigate, useParams} from "react-router-dom";
import { ModeToggle } from "@/components/mode-toggle";
import {Collection, getCollection} from "@/api/collections.ts";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { Button } from "@/components/ui/button";
import { Pencil } from "lucide-react";
import { GalleriesSection } from "@/features/galleries/components/galleries-section.tsx";
import { useState } from "react";
import { CollectionDialog } from "@/features/collections/components/collection-dialog";
import { DeleteCollectionDialog } from "@/features/collections/components/delete-collection-dialog";

type CollectionContentProps = {
    collection: Collection;
};


function CollectionContent({ collection } : CollectionContentProps) {
    const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
    const navigate = useNavigate();

    return (
        <div className="space-y-3">
            <Card className="w-full">
                <CardHeader>
                    <div className="flex justify-between items-start">
                        <div>
                            <CardTitle>{collection.name}</CardTitle>
                            <CardDescription>
                                Created on {new Date(collection.createdAt).toLocaleDateString()}
                            </CardDescription>
                        </div>
                        <div className="flex gap-2">
                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => setIsEditDialogOpen(true)}
                            >
                                <Pencil className="w-4 h-4" />
                            </Button>
                            <DeleteCollectionDialog collection={collection} onDelete={() => navigate('/collections')} />
                        </div>
                    </div>
                </CardHeader>
            </Card>

            <GalleriesSection collectionId={collection.id} />

            <CollectionDialog
                mode="edit"
                collection={collection}
                open={isEditDialogOpen}
                onOpenChange={setIsEditDialogOpen}
            />
        </div>
    );
}

export default function CollectionPage() {
    const params = useParams();
    const collectionQuery = useQuery({
        queryKey: ['collections', params.collectionId],
        queryFn: () => getCollection(params.collectionId!),
        enabled: !!params.collectionId,
    });

    const renderContent = () => {
        if (collectionQuery.status === 'pending') {
            return (
                <Card className="w-full">
                    <CardContent className="flex p-6 h-64 justify-center items-center">
                        <LoadingSpinner />
                    </CardContent>
                </Card>
            );
        }

        if (collectionQuery.status === 'error') {
            return (
                <Card className="w-full">
                    <CardContent className="p-6">
                        <div className="text-center text-red-500">
                            Error loading collection. Please try again later.
                        </div>
                    </CardContent>
                </Card>
            );
        }

        return <CollectionContent collection={collectionQuery.data} />;
    };

    return (
        <main className="w-full">
            <div className="w-full flex p-2 justify-between items-center">
                <div className="flex gap-4 items-center">
                    <SidebarTrigger/>
                    <Breadcrumb>
                        <BreadcrumbList>
                            <BreadcrumbItem>
                                <BreadcrumbLink asChild>
                                    <Link to=".." relative="path">Collections</Link>
                                </BreadcrumbLink>
                            </BreadcrumbItem>
                            <BreadcrumbSeparator/>
                            <BreadcrumbItem>
                                <BreadcrumbLink asChild>
                                    <Link to="." relative="path">
                                        {collectionQuery.data?.name ?? 'Loading...'}
                                    </Link>
                                </BreadcrumbLink>
                            </BreadcrumbItem>
                        </BreadcrumbList>
                    </Breadcrumb>
                </div>
                <ModeToggle/>
            </div>
            <div className="p-4">
                {renderContent()}
            </div>
        </main>
    );
}