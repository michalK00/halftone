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
import { getCollection } from "@/features/collections/api/collections.ts";
import { useQuery } from "@tanstack/react-query";
import {Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import {Gallery, getGallery} from "@/features/galleries/api/galleries.ts";
import {DeleteGalleryDialog} from "@/features/galleries/components/delete-gallery-dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import {ImageUp, Lock, Share2} from "lucide-react";
import {EditGalleryDialog} from "@/features/galleries/components/edit-gallery-dialog.tsx";
import {useEffect, useState} from "react";
import {Dialog, DialogContent, DialogTitle} from "@/components/ui/dialog.tsx";
import {
    Carousel,
    CarouselApi,
    CarouselContent,
    CarouselItem,
    CarouselNext,
    CarouselPrevious
} from "@/components/ui/carousel.tsx";

type GalleryContentProps = {
    gallery: Gallery;
    collectionId: string
};
// New Image type
interface ImageInfo {
    url: string;
    name: string;
    size: string;
}

// New ImageCard component
function ImageCard({ image, images, index, onSelect }: {
    image: ImageInfo;
    images: ImageInfo[];
    index: number;
    onSelect: (index: number) => void;
}) {
    return (
        <Card
            className="aspect-square flex flex-col overflow-hidden cursor-pointer"
            onClick={() => onSelect(index)}
        >
            <div className="relative flex-1">
                <img
                    src={image.url}
                    alt={image.name}
                    className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 max-w-full max-h-full object-contain"
                />
            </div>
            <CardFooter className="p-2 border-t bg-muted/50 shrink-0 flex justify-between">
                <span className="text-sm text-muted-foreground truncate">
                    {image.name}
                </span>
                <span className="text-sm text-muted-foreground truncate">
                    {image.size}
                </span>
            </CardFooter>
        </Card>
    );
}

function GalleryContent({gallery, collectionId} : GalleryContentProps) {
    const navigate = useNavigate();
    const [api, setApi] = useState<CarouselApi>()
    const [current, setCurrent] = useState<number | null>(null)
    const [selectedImageIndex, setSelectedImageIndex] = useState<number | null>(null);

    useEffect(() => {
        if (!api) {
            return
        }
        setCurrent(api.selectedScrollSnap() + 1)

        api.on("select", () => {
            setCurrent(api.selectedScrollSnap() + 1)
        })
    }, [api])

    // Example images array - replace with your actual images data
    const images: ImageInfo[] = [
        {
            url: "https://letsenhance.io/static/8f5e523ee6b2479e26ecc91b9c25261e/1015f/MainAfter.jpg",
            name: "image.png",
            size: "3.23 Mb"
        },
        {
            url: "https://plus.unsplash.com/premium_photo-1664474619075-644dd191935f",
            name: "image.png",
            size: "3.23 Mb"
        },
        {
            url: "https://i0.wp.com/picjumbo.com/wp-content/uploads/silhouette-of-a-guy-with-a-cap-at-red-sky-sunset-free-image.jpeg",
            name: "image.png",
            size: "3.23 Mb"
        },
        {
            url: "https://letsenhance.io/static/8f5e523ee6b2479e26ecc91b9c25261e/1015f/MainAfter.jpg",
            name: "image.png",
            size: "3.23 Mb"
        }
    ];

    return (
        <div className="space-y-3">
            <Card className="w-full">
                <CardHeader>
                    <div className="flex justify-between items-start ">
                        <div>
                            <CardTitle>{gallery.name}</CardTitle>
                            <CardDescription className="flex flex-col sm:flex-row sm:gap-2">
                                Created on {new Date(gallery.createdAt).toLocaleDateString()}
                                {gallery.sharingEnabled ? (
                                    <div className="flex gap-1">
                                        <Share2 className="w-4 h-4 text-green-500"/>
                                        <span className="text-sm">
                                            Until {new Date(gallery.sharingExpiryDate).toLocaleDateString()}
                                        </span>
                                    </div>
                                ) : (
                                    <div className="flex gap-1">
                                        <Lock className="w-4 h-4"/>
                                        <span className="text-sm">
                                            Disabled
                                        </span>
                                    </div>
                                )}
                            </CardDescription>

                        </div>


                        <div className="flex gap-2 justify-end">
                            <Button
                                variant="outline"
                                size="icon"
                            >
                                {gallery.sharingEnabled ? (
                                    <Lock className="w-4 h-4" />
                                ) : (
                                    <Share2 className="w-4 h-4" />
                                )}
                            </Button>
                            <EditGalleryDialog gallery={gallery}/>
                            <DeleteGalleryDialog gallery={gallery} onDelete={() => {navigate(`/collections/${collectionId}`)}}/>
                        </div>
                    </div>
                </CardHeader>
            </Card>

            <Card className="w-full">
                <CardHeader>
                    <div className="flex justify-between items-start ">
                        <CardTitle>Photos</CardTitle>
                        <div className="flex gap-2 justify-end">
                            <Button variant="outline" size="icon">
                                <ImageUp className="w-4 h-4" />
                            </Button>
                        </div>
                    </div>
                </CardHeader>
                <CardContent className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 auto-rows-fr">
                    {images.map((image, index) => (
                        <ImageCard
                            key={index}
                            image={image}
                            images={images}
                            index={index}
                            onSelect={setSelectedImageIndex}
                        />
                    ))}
                </CardContent>
            </Card>

            <Dialog open={selectedImageIndex !== null} onOpenChange={() => setSelectedImageIndex(null)}>
                <DialogContent className="max-w-7xl w-[90%] sm:w-10/12 rounded-lg">
                    <Carousel setApi={setApi}  className="w-full" opts={{startIndex: selectedImageIndex!, loop: true}}>
                        <CarouselContent>
                            {images.map((image, index) => (
                                <CarouselItem  key={index}>
                                    <DialogTitle className="text-center pb-2 text-sm sm:text-base">
                                        {image.name}
                                    </DialogTitle>
                                    <div className="relative w-full aspect-video ">
                                        <img
                                            src={image.url}
                                            alt={image.name}
                                            className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 max-w-full max-h-full object-contain"
                                        />
                                    </div>
                                </CarouselItem>
                            ))}
                        </CarouselContent>
                        <CarouselPrevious className="hidden sm:flex"/>
                        <CarouselNext className="hidden sm:flex"/>
                    </Carousel>
                    <div className="text-center text-sm sm:text-base -mt-2 -mb-2">Photo {current} out of {images.length}</div>
                </DialogContent>
            </Dialog>
        </div>
    );
}


export default function GalleryPage() {
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

        return <GalleryContent gallery={galleryQuery.data} collectionId={collectionId!} />;
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