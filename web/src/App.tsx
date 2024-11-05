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

const queryClient = new QueryClient()

function App() {
  // const { isAuthenticated, getAccessTokenSilently } = useAuth0();
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
        <QueryClientProvider client={queryClient}>
          <SidebarProvider>
            <AppSidebar />
            <Routes>
              <Route path="/collections" element={<CollectionsPage/>}/>
              <Route path="/collections/:collectionId"  element={<CollectionPage/>}/>
              <Route path="/orders" element={<Orders/>}/>
              <Route path="*" element={<Navigate to="/collections" replace/>} />
            </Routes>
          </SidebarProvider>
        </QueryClientProvider>
    </ThemeProvider>
  );
}

export default App;
