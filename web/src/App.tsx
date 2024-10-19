import "./App.css";
import { ThemeProvider } from "./components/theme-provider";
import LoginButton from "./components/login-button";
import { ModeToggle } from "./components/mode-toggle";
import { SidebarProvider, SidebarTrigger } from "./components/ui/sidebar";
import { AppSidebar } from "./components/app-sidebar";
import { useAuth0 } from "@auth0/auth0-react";
import LogoutButton from "./components/logout-button";
import Profile from "./components/profile";
import { Button } from "./components/ui/button";

function App() {
  const { isAuthenticated, getAccessTokenSilently } = useAuth0();

  const showApiCall = () => {
    getAccessTokenSilently().then((token) => {
      fetch("http://localhost:8080/collections", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
        .then((response) => response.json())
        .then((data) => alert(JSON.stringify(data)));
    });
  };

  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <SidebarProvider>
        <AppSidebar />
        <main className="w-full">
          <div className="w-full flex p-2 justify-between items-center">
            <SidebarTrigger />
            <div className="flex gap-3">
              <ModeToggle></ModeToggle>
              {isAuthenticated ? (
                <>
                  <LogoutButton />
                  <Profile />
                  <Button onClick={() => showApiCall()}>Call api</Button>
                </>
              ) : (
                <LoginButton />
              )}
            </div>
          </div>
        </main>
      </SidebarProvider>
    </ThemeProvider>
  );
}

export default App;
