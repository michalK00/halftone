import "./App.css";
import { ThemeProvider } from "./components/theme-provider";
import { SidebarProvider } from "./components/ui/sidebar";
import { AppSidebar } from "./components/app-sidebar";
import {Navigate, Route, Routes} from "react-router-dom";
import Orders from "./features/orders/pages/orders-page.tsx";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import CollectionsPage from "@/features/collections/pages/collections-page.tsx";
import {Toaster} from "@/components/ui/toaster.tsx";
import CollectionPage from "@/features/collections/pages/collection-page.tsx";
import GalleryPage from "@/features/galleries/pages/gallery-page.tsx";
import { AuthProvider, useAuth } from "@/context/auth-context.tsx";
import SignInPage from "@/features/auth/pages/SignInPage.tsx";
import SignUpPage from "@/features/auth/pages/SignUpPage.tsx";
import VerifyAccountPage from "@/features/auth/pages/VerifyAccountPage.tsx";
import ClientPage from "@/features/client/client-page.tsx";
const queryClient = new QueryClient()

function AppContent() {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-screen">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
            </div>
        );
    }

    if (!isAuthenticated) {
        return (
            <Routes>
                <Route path="/client/*" element={<ClientPage />} />
                <Route path="/signin" element={<SignInPage />} />
                <Route path="/signup" element={<SignUpPage />} />
                <Route path="/verify-account" element={<VerifyAccountPage />} />
                {/*<Route path="/forgot-password" element={<ForgotPasswordPage />} />*/}
                {/*<Route path="/reset-password" element={<ResetPasswordPage />} />*/}
                <Route path="*" element={<Navigate to="/signin" replace />} />
            </Routes>
        );
    }

    return (

            <SidebarProvider>
                <AppSidebar />
                <Routes>
                    <Route path="/collections" element={<CollectionsPage />} />
                    <Route path="/collections/:collectionId" element={<CollectionPage />} />
                    <Route path="/galleries/:galleryId" element={<GalleryPage />} />
                    <Route path="/orders" element={<Orders />} />
                    <Route path="*" element={<Navigate to="/collections" replace />} />
                </Routes>
            </SidebarProvider>

    );
}

function App() {
    return (
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
            <QueryClientProvider client={queryClient}>
            <AuthProvider>
                <Toaster />
                <AppContent />
            </AuthProvider>
            </QueryClientProvider>
        </ThemeProvider>
    );
}

export default App;