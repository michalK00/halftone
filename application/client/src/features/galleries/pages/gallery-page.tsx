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
import { getCollection } from "@/api/collections.ts";
import { useQuery } from "@tanstack/react-query";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import {getGallery} from "@/api/galleries.ts";
import PhotoGallery from "@/features/galleries/components/photo-gallery.tsx";
import {Lock, Share2} from "lucide-react";
import {EditGalleryDialog} from "@/features/galleries/components/edit-gallery-dialog.tsx";
import {DeleteGalleryDialog} from "@/features/galleries/components/delete-gallery-dialog.tsx";
import ImageShareModal from "@/features/galleries/components/image-share-modal.tsx";


export default function GalleryPage() {
    const navigate = useNavigate();
    const params = useParams();
    const galleryQuery = useQuery({
        queryKey: ['gallery', params.galleryId],
        queryFn: () => getGallery(params.galleryId!),
        enabled: !!params.galleryId,
    });

    const collectionId = galleryQuery.data?.collectionId;

    const collectionQuery = useQuery({
        queryKey: ['collections', collectionId],
        queryFn: () => getCollection(collectionId!),
        enabled: !!collectionId,
    });

    function createSharingUrl(galleryPath: string) {
        return `https://prod-alb-1747799725.eu-north-1.elb.amazonaws.com/client${galleryPath}`;
    }

    const renderContent = () => {
        if (galleryQuery.status === 'pending') {
            return (
                <Card className="w-full">
                    <CardContent className="flex p-6 h-64 justify-center items-center">
                        <LoadingSpinner />
                    </CardContent>
                </Card>
            );
        }

        if (galleryQuery.status === 'error') {
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

        return <div className="space-y-3">
                <Card className="w-full">
                    <CardHeader>
                        <div className="flex justify-between items-start">
                            <div>
                                <CardTitle>{galleryQuery.data.name}</CardTitle>
                                <CardDescription className="flex flex-col sm:flex-row sm:gap-2">
                                    Created on {new Date(galleryQuery.data.createdAt).toLocaleDateString()}
                                    {galleryQuery.data.sharing.sharingEnabled ? (
                                        <div className="flex gap-1">
                                            <Share2 className="w-4 h-4 text-green-500"/>
                                            <span className="text-sm">
                                                    Until {new Date(galleryQuery.data.sharing.sharingExpiryDate).toLocaleDateString()}
                                                </span>
                                        </div>
                                    ) : (
                                        <div className="flex gap-1">
                                            <Lock className="w-4 h-4"/>
                                            <span className="text-sm">Disabled</span>
                                        </div>
                                    )}
                                </CardDescription>
                            </div>

                            <div className="flex gap-2 justify-end">
                                <ImageShareModal
                                    galleryId={galleryQuery.data.id}
                                    sharingEnabled={galleryQuery.data.sharing.sharingEnabled}
                                    sharingExpiryDate={new Date(galleryQuery.data.sharing.sharingExpiryDate)}
                                    sharingUrl={createSharingUrl(galleryQuery.data.sharing.sharingUrl)}
                                />
                                <EditGalleryDialog gallery={galleryQuery.data}/>
                                <DeleteGalleryDialog gallery={galleryQuery.data} onDelete={() => {navigate(`/collections/${collectionId}`)}}/>
                            </div>
                        </div>
                    </CardHeader>
                </Card>
                <PhotoGallery gallery={galleryQuery.data} />
            </div>;
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
                                    <Link to="/collections">Collections</Link>
                                </BreadcrumbLink>
                            </BreadcrumbItem>
                            <BreadcrumbSeparator/>
                            <BreadcrumbItem>
                                <BreadcrumbLink asChild>
                                    <Link to={`/collections/${collectionId}`}>
                                        {collectionQuery.data?.name ?? 'Loading...'}
                                    </Link>
                                </BreadcrumbLink>
                            </BreadcrumbItem>
                            <BreadcrumbSeparator/>
                            <BreadcrumbItem>
                                <BreadcrumbLink asChild>
                                    <Link to="." relative="path">
                                        Gallery: {galleryQuery.data?.name ?? 'Loading...'}
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

