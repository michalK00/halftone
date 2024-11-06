import "./App.css";
import { ThemeProvider } from "./components/theme-provider";
import { SidebarProvider } from "./components/ui/sidebar";
import { AppSidebar } from "./components/app-sidebar";
import {Navigate, Route, Routes} from "react-router-dom";
import Orders from "./pages/orders.tsx";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import CollectionsPage from "@/features/collections/pages/collections-page.tsx";
import {Toaster} from "@/components/ui/toaster.tsx";
import CollectionPage from "@/features/collections/pages/collection-page.tsx";
import {useAuth0} from "@auth0/auth0-react";
import LoginButton from "@/components/login-button.tsx";
import GalleryPage from "@/features/galleries/pages/gallery-page.tsx";

const queryClient = new QueryClient()

function App() {
  const { isAuthenticated } = useAuth0();
  //
  // const showApiCall = () => {
  //   getAccessTokenSilently().then((token) => {
  //     fetch("http://localhost:8080/collections", {
  //       headers: {
  //         Authorization: `Bearer ${token}`,
  //       },
  //     })
  //       .then((response) => response.json())
  //       .then((data) => alert(JSON.stringify(data)));
  //   });
  // };

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Toaster />
      {isAuthenticated ? <> <QueryClientProvider client={queryClient}>
        <SidebarProvider>
          <AppSidebar />
          <Routes>
            <Route path="/collections" element={<CollectionsPage/>}/>
            <Route path="/collections/:collectionId"  element={<CollectionPage/>}/>
            <Route path="/galleries/:galleryId"  element={<GalleryPage/>}/>
            <Route path="/orders" element={<Orders/>}/>
            <Route path="*" element={<Navigate to="/collections" replace/>} />
          </Routes>
        </SidebarProvider>
      </QueryClientProvider></> : <div className="flex justify-center items-center h-screen	"><LoginButton/></div>}

    </ThemeProvider>
  );
}

export default App;
