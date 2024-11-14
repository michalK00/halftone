import { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { ImageUp, Trash2 } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import {toast} from "@/hooks/use-toast.ts";
import {confirmUpload, getUploadUrls, UploadRequestBody, uploadToAWS} from "@/features/galleries/api/photos.ts";
import {useQueryClient} from "@tanstack/react-query";

type FileWithPreview = {
    file: File;
    preview: string;
};

type ImageUploadModalProps = {
    galleryId: string;
};

function ImageUploadModal({galleryId}: ImageUploadModalProps){
    const [isOpen, setIsOpen] = useState(false);
    const [selectedFiles, setSelectedFiles] = useState<FileWithPreview[]>([]);
    const [isUploading, setIsUploading] = useState(false);
    const queryClient = useQueryClient()

    const handleFiles = (files: FileList) => {
        const imageFiles = Array.from(files).filter(file => file.type.startsWith('image/'));

        const newFiles = imageFiles.map(file => ({
            file,
            preview: URL.createObjectURL(file)
        }));

        setSelectedFiles(prev => [...prev, ...newFiles]);
    };

    const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        handleFiles(e.dataTransfer.files);
    };

    const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files?.length) {
            handleFiles(e.target.files);
        }
    };

    const removeFile = (index: number) => {
        setSelectedFiles(prev => {
            const newFiles = [...prev];
            URL.revokeObjectURL(newFiles[index].preview);
            newFiles.splice(index, 1);
            return newFiles;
        });
    };

    async function uploadFiles(selectedFiles: FileWithPreview[]) {
        setIsUploading(true)
        const body: UploadRequestBody[] = selectedFiles.map(fileWithPreview => {return {originalFilename: fileWithPreview.file.name}})

        const photoUploads = await getUploadUrls(galleryId, body)

        for (let i = 0; i < selectedFiles.length; i++) {
            await uploadToAWS(selectedFiles[i].file, photoUploads[i])
            await confirmUpload(photoUploads[i].id);
        }
    }

    const handleUpload = async () => {
        setIsUploading(true);

        try {
            selectedFiles.forEach(file => URL.revokeObjectURL(file.preview));

            await uploadFiles(selectedFiles)

            toast({
                title: "Success!",
                description: `Successfully uploaded ${selectedFiles.length} ${selectedFiles.length === 1 ? 'image' : 'images'}`,
                variant: "default",
            });

            setSelectedFiles([]);
            setIsUploading(false);
            queryClient.invalidateQueries({queryKey: ['photos']})
            setIsOpen(false);

        } catch (error) {
            setIsUploading(false);
            toast({
                title: "Upload failed",
                description: error instanceof Error ? error.message : "An unexpected error occurred",
                variant: "destructive",
            });
            console.error('Upload failed:', error);
        }
    };

    const closeModal = () => {
        selectedFiles.forEach(file => URL.revokeObjectURL(file.preview));
        setSelectedFiles([]);
        setIsOpen(false);
    };

    return (
        <>
            <Button variant="outline" size="icon" onClick={() => setIsOpen(true)}>
                <ImageUp className="w-4 h-4" />
            </Button>

            <Dialog open={isOpen} onOpenChange={closeModal}>
                <DialogContent className="sm:max-w-xl">
                    <DialogHeader>
                        <DialogTitle>Upload Images</DialogTitle>
                    </DialogHeader>

                    <div className="space-y-4">
                        <div
                            className="border-2 border-dashed border-gray-200 rounded-lg p-6 cursor-pointer hover:border-gray-300 transition-colors"
                            onDragOver={(e) => e.preventDefault()}
                            onDrop={handleDrop}
                            onClick={() => document.getElementById('file-upload')?.click()}
                        >
                            <div className="flex flex-col items-center gap-2 text-center">
                                <div className="w-10 h-10 rounded-full bg-blue-50 flex items-center justify-center">
                                    <ImageUp className="w-5 h-5 text-blue-500" />
                                </div>
                                <div>
                                    <p className="text-sm font-medium">Select image files to upload</p>
                                    <p className="text-xs text-gray-500">or drag and drop them here</p>
                                </div>
                            </div>
                            <input
                                id="file-upload"
                                type="file"
                                className="hidden"
                                accept="image/*"
                                multiple
                                onChange={handleFileInput}
                            />
                        </div>

                        {selectedFiles.length > 0 && (
                            <ScrollArea className="h-[200px] w-full rounded-md border p-4">
                                <div className="grid grid-cols-4 gap-4">
                                    {selectedFiles.map((file, index) => (
                                        <div key={index} className="relative group">
                                            <img
                                                src={file.preview}
                                                alt={file.file.name}
                                                className="w-full aspect-square rounded-lg object-cover"
                                            />
                                            <Button
                                                variant="destructive"
                                                size="icon"
                                                className="absolute top-1 right-1 w-6 h-6 opacity-0 group-hover:opacity-100 transition-opacity"
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    removeFile(index);
                                                }}
                                            >
                                                <Trash2 className="w-3 h-3" />
                                            </Button>
                                            <p className="text-xs mt-1 truncate">{file.file.name}</p>
                                        </div>
                                    ))}
                                </div>
                            </ScrollArea>
                        )}
                    </div>

                    <DialogFooter>
                        <div className="flex gap-2 justify-end">
                            <Button
                                variant="outline"
                                onClick={closeModal}
                                disabled={isUploading}
                            >
                                Cancel
                            </Button>
                            <Button
                                onClick={handleUpload}
                                disabled={selectedFiles.length === 0 || isUploading}
                            >
                                {isUploading ? (
                                    <>
                                        Uploading...
                                    </>
                                ) : (
                                    'Upload'
                                )}
                            </Button>
                        </div>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
};

export default ImageUploadModal;