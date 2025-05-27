import {Gallery} from "@/api/galleries.ts";
import {useEffect, useState} from "react";
import {
    Carousel,
    CarouselApi,
    CarouselContent,
    CarouselItem,
    CarouselNext,
    CarouselPrevious
} from "@/components/ui/carousel.tsx";
import {deletePhoto, getPhotos, Photo} from "@/api/photos.ts";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import ImageCard from "@/features/galleries/components/image-card.tsx";
import {Dialog, DialogContent, DialogTitle} from "@/components/ui/dialog.tsx";
import {useQuery, useQueryClient} from "@tanstack/react-query";
import ImageUploadModal from "@/features/galleries/components/image-upload-modal.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Check, Trash2, X} from "lucide-react";
import {useToast} from "@/hooks/use-toast.ts";

type GalleryContentProps = {
    gallery: Gallery;
};

function PhotoGallery({gallery}: GalleryContentProps) {
    const [api, setApi] = useState<CarouselApi>();
    const [current, setCurrent] = useState<number | null>(null);
    const [selectedImageIndex, setSelectedImageIndex] = useState<number | null>(null);
    const [deleting, setDeletingMode] = useState<boolean>(false);
    const [indexesToDelete, setIndexesToDelete] = useState<number[]>([]);
    const {toast} = useToast();
    const queryClient = useQueryClient();

    const {data: images, isLoading, isError} = useQuery({
        queryKey: ['photos', gallery.id],
        queryFn: () => getPhotos(gallery.id),
    });

    function switchDeleteMode(e: React.FormEvent) {
        e.preventDefault();
        setDeletingMode(!deleting);
        setIndexesToDelete([]);
    }

    const handleSelect = (index: number) => {
        if (deleting) {
            setIndexesToDelete(prev =>
                prev.includes(index)
                    ? prev.filter(i => i !== index)
                    : [...prev, index]
            );
        } else {
            setSelectedImageIndex(index);
        }
    };

    const handleDelete = async () => {
        if (images && indexesToDelete.length > 0) {
            try {
                const photos = indexesToDelete.map(index => images[index])
                for (let i = 0; i < photos.length; i++) {
                    await deletePhoto(photos[i].id)
                }
                const queryKey = ['photos', gallery.id];
                queryClient.setQueryData(queryKey, (oldData: Photo[]) => {
                    return oldData.filter((_, index) => !indexesToDelete.includes(index))
                })

                toast({
                    title: "Success",
                    description: `Deleted ${indexesToDelete.length} photo(s)`,
                });


                setDeletingMode(false);
                setIndexesToDelete([]);
            } catch (error) {
                console.error('Delete failed:', error);
                toast({
                    title: "Error",
                    description: "Failed to delete photos",
                    variant: "destructive",
                });
            }
        }
    };

    useEffect(() => {
        if (!api) return;
        setCurrent(api.selectedScrollSnap() + 1);
        api.on("select", () => {
            setCurrent(api.selectedScrollSnap() + 1);
        });
    }, [api]);

    if (isLoading) {
        return (
            <div className="space-y-3">
                <Card className="w-full">
                    <CardHeader>
                        <div className="animate-pulse h-8 rounded w-1/4"></div>
                    </CardHeader>
                </Card>
                <Card className="w-full p-6">
                    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                        {[...Array(4)].map((_, i) => (
                            <div key={i} className="animate-pulse dark:bg-zinc-800 light:bg-zinc-600 rounded aspect-square"></div>
                        ))}
                    </div>
                </Card>
            </div>
        );
    }

    if (isError) {
        return (
            <Card className="w-full p-6">
                <div className="text-center text-red-500">
                    Error loading photos. Please try again later.
                </div>
            </Card>
        );
    }

    return (
        <>
            <Card className="w-full">
                <CardHeader>
                    <div className="flex justify-between items-start">
                        <CardTitle>Photos <span className="text-sm text-muted-foreground">{deleting && `(${indexesToDelete.length} selected)`}</span></CardTitle>
                        <div className="flex gap-2 justify-end">
                            <ImageUploadModal galleryId={gallery.id}/>
                            {deleting ? (
                                <>
                                    <Button
                                        variant="destructive"
                                        size="icon"
                                        onClick={handleDelete}
                                        disabled={indexesToDelete.length === 0}
                                    >
                                        <Check className="w-4 h-4"/>
                                    </Button>
                                    <Button size="icon" onClick={switchDeleteMode}>
                                        <X className="w-4 h-4"/>
                                    </Button>
                                </>
                            ) : (
                                <Button variant="destructive" size="icon" onClick={switchDeleteMode}>
                                    <Trash2 className="w-4 h-4"/>
                                </Button>
                            )}
                        </div>
                    </div>
                </CardHeader>
                <CardContent className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 auto-rows-fr">
                    {images?.map((image, index) => (
                        <ImageCard
                            key={index}
                            image={image}
                            index={index}
                            onSelect={handleSelect}
                            isSelected={indexesToDelete.includes(index)}
                            selectionMode={deleting}
                        />
                    ))}
                </CardContent>
            </Card>

            <Dialog open={selectedImageIndex !== null} onOpenChange={() => setSelectedImageIndex(null)}>
                <DialogContent className="max-w-7xl w-[90%] sm:w-10/12 rounded-lg">
                    <Carousel setApi={setApi} className="w-full" opts={{startIndex: selectedImageIndex!, loop: true}}>
                        <CarouselContent>
                            {images?.map((image, index) => (
                                <CarouselItem key={index}>
                                    <DialogTitle className="text-center pb-2 text-sm sm:text-base">
                                        {image.originalFilename}
                                    </DialogTitle>
                                    <div className="relative w-full aspect-video">
                                        <img
                                            src={image.url}
                                            alt={image.originalFilename}
                                            className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 max-w-full max-h-full object-contain"
                                        />
                                    </div>
                                </CarouselItem>
                            ))}
                        </CarouselContent>
                        <CarouselPrevious className="hidden sm:flex"/>
                        <CarouselNext className="hidden sm:flex"/>
                    </Carousel>
                    <div className="text-center text-sm sm:text-base -mt-2 -mb-2">
                        Photo {current} out of {images?.length}
                    </div>
                </DialogContent>
            </Dialog>
        </>
    );
}

export default PhotoGallery;