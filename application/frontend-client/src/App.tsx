import { useState, useEffect } from 'react';
import { ChevronLeft, ShoppingCart, Check, X, Image } from 'lucide-react';
import axios from 'axios';

// Configure axios base URL
const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
});

// Types
interface Gallery {
    id: string;
    name: string;
    photoOptions: {
        downsize: boolean;
        watermark: boolean;
    };
}

interface Photo {
    id: string;
    originalFilename: string;
    url: string;
}

interface SelectedPhoto {
    photoId: string;
    photo: Photo;
}

const App = () => {
    const [gallery, setGallery] = useState<Gallery | null>(null);
    const [photos, setPhotos] = useState<Photo[]>([]);
    const [selectedPhotos, setSelectedPhotos] = useState<SelectedPhoto[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [showOrderForm, setShowOrderForm] = useState(false);
    const [orderSuccess, setOrderSuccess] = useState(false);
    const [selectedPhotoView, setSelectedPhotoView] = useState<Photo | null>(null);

    // Order form state
    const [email, setEmail] = useState('');
    const [comment, setComment] = useState('');
    const [submitting, setSubmitting] = useState(false);

    // Extract token from URL
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get('token');
    const galleryId = window.location.pathname.split('/').pop();

    useEffect(() => {
        if (token && galleryId) {
            fetchGalleryData();
        } else {
            setError('Invalid gallery link');
            setLoading(false);
        }
    }, [token, galleryId]);

    const fetchGalleryData = async () => {
        try {
            // Set authorization header for all requests
            const config = {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            };

            // Fetch gallery info
            const galleryRes = await api.get(`/api/v1/client/galleries/${galleryId}`, config);
            setGallery(galleryRes.data);

            // Fetch photos
            const photosRes = await api.get(`/api/v1/client/galleries/${galleryId}/photos`, config);
            setPhotos(photosRes.data);
        } catch (err) {
            if (axios.isAxiosError(err)) {
                if (err.response?.status === 401) {
                    setError('Access denied - invalid or expired token');
                } else if (err.response?.status === 404) {
                    setError('Gallery not found');
                } else {
                    setError('Failed to load gallery');
                }
            } else {
                setError('An unexpected error occurred');
            }
        } finally {
            setLoading(false);
        }
    };

    const togglePhotoSelection = (photo: Photo) => {
        setSelectedPhotos(prev => {
            const exists = prev.find(p => p.photoId === photo.id);
            if (exists) {
                return prev.filter(p => p.photoId !== photo.id);
            }
            return [...prev, { photoId: photo.id, photo }];
        });
    };

    const handleSubmitOrder = async () => {
        if (!email) {
            alert('Please enter your email address');
            return;
        }

        setSubmitting(true);

        try {
            const config = {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            };

            await api.post(`/api/v1/client/galleries/${galleryId}`, {
                clientEmail: email,
                comment: comment,
                photoIds: selectedPhotos.map(p => p.photoId)
            }, config);

            setOrderSuccess(true);
            setShowOrderForm(false);
            setSelectedPhotos([]);
            setEmail('');
            setComment('');
        } catch (err) {
            if (axios.isAxiosError(err)) {
                alert(err.response?.data?.error || 'Failed to submit order. Please try again.');
            } else {
                alert('Failed to submit order. Please try again.');
            }
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <div className="min-h-screen bg-zinc-900 flex items-center justify-center">
                <div className="text-white">Loading...</div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen bg-zinc-900 flex items-center justify-center">
                <div className="text-center">
                    <h1 className="text-2xl font-bold text-white mb-4">Access Denied</h1>
                    <p className="text-gray-400">{error}</p>
                </div>
            </div>
        );
    }

    if (orderSuccess) {
        return (
            <div className="min-h-screen bg-zinc-900 flex items-center justify-center">
                <div className="text-center bg-zinc-800 p-8 rounded-lg max-w-md">
                    <div className="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mx-auto mb-4">
                        <Check className="w-8 h-8 text-white" />
                    </div>
                    <h1 className="text-2xl font-bold text-white mb-4">Order Submitted Successfully!</h1>
                    <p className="text-gray-400 mb-6">
                        We've received your order and will process it shortly.
                        You'll receive a confirmation email.
                    </p>
                    <button
                        onClick={() => {
                            setOrderSuccess(false);
                            setShowOrderForm(false);
                        }}
                        className="bg-zinc-700 text-white px-6 py-2 rounded-md hover:bg-zinc-600 transition-colors"
                    >
                        Back to Gallery
                    </button>
                </div>
            </div>
        );
    }

    if (showOrderForm) {
        return (
            <div className="min-h-screen bg-zinc-900 p-4">
                <div className="max-w-4xl mx-auto">
                    <button
                        onClick={() => setShowOrderForm(false)}
                        className="flex items-center text-gray-400 hover:text-white mb-6 transition-colors"
                    >
                        <ChevronLeft className="w-5 h-5 mr-1" />
                        Back to Gallery
                    </button>

                    <div className="bg-zinc-800 rounded-lg p-6">
                        <h2 className="text-2xl font-bold text-white mb-6">Create Order</h2>

                        <div className="mb-6">
                            <h3 className="text-lg font-medium text-white mb-3">Selected Photos ({selectedPhotos.length})</h3>
                            <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                                {selectedPhotos.map(({ photo }) => (
                                    <div key={photo.id} className="relative group">
                                        <img
                                            src={photo.url}
                                            alt={photo.originalFilename}
                                            className="w-full aspect-square object-cover rounded"
                                        />
                                        <button
                                            onClick={() => togglePhotoSelection(photo)}
                                            className="absolute top-2 right-2 bg-red-600 text-white p-1 rounded opacity-0 group-hover:opacity-100 transition-opacity"
                                        >
                                            <X className="w-4 h-4" />
                                        </button>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="space-y-4">
                            <div>
                                <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-1">
                                    Email Address *
                                </label>
                                <input
                                    type="email"
                                    id="email"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    required
                                    className="w-full px-3 py-2 bg-zinc-700 border border-zinc-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    placeholder="your@email.com"
                                />
                            </div>

                            <div>
                                <label htmlFor="comment" className="block text-sm font-medium text-gray-300 mb-1">
                                    Special Instructions
                                </label>
                                <textarea
                                    id="comment"
                                    value={comment}
                                    onChange={(e) => setComment(e.target.value)}
                                    rows={4}
                                    className="w-full px-3 py-2 bg-zinc-700 border border-zinc-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    placeholder="Any special requests for your order..."
                                />
                            </div>

                            <div className="flex gap-3 pt-4">
                                <button
                                    type="submit"
                                    disabled={submitting}
                                    onClick={handleSubmitOrder}
                                    className="flex-1 bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    {submitting ? 'Submitting...' : 'Submit Order'}
                                </button>
                                <button
                                    type="button"
                                    onClick={() => setShowOrderForm(false)}
                                    className="px-6 py-2 bg-zinc-700 text-white rounded-md hover:bg-zinc-600 transition-colors"
                                >
                                    Cancel
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    if (selectedPhotoView) {
        return (
            <div className="fixed inset-0 bg-black flex items-center justify-center z-50">
                <button
                    onClick={() => setSelectedPhotoView(null)}
                    className="absolute top-4 right-4 text-white p-2 hover:bg-white/10 rounded-full transition-colors"
                >
                    <X className="w-6 h-6" />
                </button>
                <img
                    src={selectedPhotoView.url}
                    alt={selectedPhotoView.originalFilename}
                    className="max-w-full max-h-full object-contain"
                />
                <div className="absolute bottom-4 left-4 text-white">
                    <p className="text-sm opacity-75">{selectedPhotoView.originalFilename}</p>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-zinc-900">
            {/* Header */}
            <div className="bg-zinc-800 border-b border-zinc-700 p-4">
                <div className="max-w-7xl mx-auto flex justify-between items-center">
                    <h1 className="text-xl font-bold text-white">{gallery?.name}</h1>
                    {selectedPhotos.length > 0 && (
                        <button
                            onClick={() => setShowOrderForm(true)}
                            className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors flex items-center gap-2"
                        >
                            <ShoppingCart className="w-4 h-4" />
                            Create Order ({selectedPhotos.length})
                        </button>
                    )}
                </div>
            </div>

            {/* Photo Grid */}
            <div className="max-w-7xl mx-auto p-4">
                {photos.length === 0 ? (
                    <div className="text-center py-16">
                        <Image className="w-16 h-16 text-gray-600 mx-auto mb-4" />
                        <p className="text-gray-400">No photos in this gallery yet</p>
                    </div>
                ) : (
                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
                        {photos.map((photo) => {
                            const isSelected = selectedPhotos.some(p => p.photoId === photo.id);
                            return (
                                <div
                                    key={photo.id}
                                    className={`relative group cursor-pointer rounded-lg overflow-hidden ${
                                        isSelected ? 'ring-2 ring-blue-500' : ''
                                    }`}
                                >
                                    <img
                                        src={photo.url}
                                        alt={photo.originalFilename}
                                        className="w-full aspect-square object-cover"
                                        onClick={() => setSelectedPhotoView(photo)}
                                    />
                                    <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-40 transition-all duration-200">
                                        <div className="absolute top-2 right-2 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                            <button
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    togglePhotoSelection(photo);
                                                }}
                                                className={`p-2 rounded-full transition-colors ${
                                                    isSelected
                                                        ? 'bg-blue-600 text-white'
                                                        : 'bg-white/90 text-gray-800 hover:bg-white'
                                                }`}
                                            >
                                                {isSelected ? <Check className="w-4 h-4" /> : <ShoppingCart className="w-4 h-4" />}
                                            </button>
                                        </div>
                                    </div>
                                    <div className="absolute bottom-0 left-0 right-0 p-2 bg-gradient-to-t from-black/70 to-transparent opacity-0 group-hover:opacity-100 transition-opacity">
                                        <p className="text-white text-sm truncate">{photo.originalFilename}</p>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>

            {/* Selection Bar */}
            {selectedPhotos.length > 0 && !showOrderForm && (
                <div className="fixed bottom-0 left-0 right-0 bg-zinc-800 border-t border-zinc-700 p-4">
                    <div className="max-w-7xl mx-auto flex justify-between items-center">
                        <p className="text-white">
                            {selectedPhotos.length} photo{selectedPhotos.length !== 1 ? 's' : ''} selected
                        </p>
                        <div className="flex gap-3">
                            <button
                                onClick={() => setSelectedPhotos([])}
                                className="px-4 py-2 bg-zinc-700 text-white rounded-md hover:bg-zinc-600 transition-colors"
                            >
                                Clear Selection
                            </button>
                            <button
                                onClick={() => setShowOrderForm(true)}
                                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors flex items-center gap-2"
                            >
                                <ShoppingCart className="w-4 h-4" />
                                Create Order
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default App;