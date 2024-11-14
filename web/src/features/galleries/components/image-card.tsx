import {Card, CardFooter} from "@/components/ui/card.tsx";
import {Photo} from "@/features/galleries/api/photos.ts";
import {AlertCircle} from "lucide-react";
import {useState} from "react";
import {LoadingSpinner} from "@/components/ui/loading-spinner.tsx";



function ImageCard({ image, index, onSelect }: {
    image: Photo;
    index: number;
    onSelect: (index: number) => void;
}) {
    const [hasError, setHasError] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    return (
        <Card
            className="aspect-square flex flex-col overflow-hidden cursor-pointer"
            onClick={() => onSelect(index)}
        >
            <div className="relative flex-1">
                {isLoading && (
                    <div className="absolute inset-0 flex items-center justify-center dark:bg-zinc-900 light:bg-zinc-600">
                        <LoadingSpinner/>
                    </div>
                )}
                {hasError ? (
                    <div className="absolute inset-0 flex flex-col items-center justify-center dark:text-zinc-100 light:text-zinc-700">
                        <AlertCircle className="w-8 h-8 mb-2" />
                        <span className="text-sm">Failed to load image</span>
                    </div>
                ) : (
                    <img
                        src={image.url}
                        alt={image.originalFilename}
                        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 max-w-full max-h-full object-contain"
                        onError={() => {
                            setHasError(true);
                            setIsLoading(false);
                        }}
                        onLoad={() => setIsLoading(false)}
                    />
                )}
            </div>
            <CardFooter className="p-2 border-t bg-muted/50 shrink-0 flex justify-between">
                <span className="text-sm text-muted-foreground truncate">
                    {image.originalFilename}
                </span>

            </CardFooter>
        </Card>
    );
}

export default ImageCard