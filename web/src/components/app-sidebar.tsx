import {
    Sidebar,
    SidebarContent, SidebarFooter,

    SidebarGroup, SidebarGroupContent, SidebarGroupLabel,
    SidebarMenu, SidebarMenuButton, SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Link } from "react-router-dom";
import {ChevronUp, Images, ShoppingBag, User2,} from "lucide-react"
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from "@radix-ui/react-dropdown-menu";

const items = [
    {
        title: "Collections",
        url: "/collections",
        icon: Images,
    },
    {
        title: "Orders",
        url: "/orders",
        icon: ShoppingBag,
    },
]

export function AppSidebar() {
  return (
      <Sidebar>
          <SidebarContent>
              <SidebarGroup>
                  <SidebarGroupLabel>Application</SidebarGroupLabel>
                  <SidebarGroupContent>
                      <SidebarMenu>
                          {items.map((item) => (
                              <SidebarMenuItem key={item.title}>
                                  <SidebarMenuButton asChild>
                                      <Link to={item.url}>
                                          <item.icon />
                                          <span>{item.title}</span>
                                      </Link>
                                  </SidebarMenuButton>
                              </SidebarMenuItem>
                          ))}
                      </SidebarMenu>
                  </SidebarGroupContent>
              </SidebarGroup>
          </SidebarContent>
          <SidebarFooter>
              <SidebarMenu>
                  <SidebarMenuItem>

                      <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                              <SidebarMenuButton>
                                  <User2 /> Username
                                  <ChevronUp className="ml-auto" />
                              </SidebarMenuButton>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent
                              side="top"
                              className="w-[--radix-popper-anchor-width]"
                          >
                              <DropdownMenuItem>
                                  <span>Account</span>
                              </DropdownMenuItem>
                              <DropdownMenuItem>
                                  <span>Sign out</span>
                              </DropdownMenuItem>
                          </DropdownMenuContent>
                      </DropdownMenu>
                  </SidebarMenuItem>
              </SidebarMenu>
          </SidebarFooter>
      </Sidebar>
  );
}
