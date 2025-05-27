import {SidebarTrigger} from "@/components/ui/sidebar.tsx";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList
} from "@/components/ui/breadcrumb.tsx";
import {Link} from "react-router-dom";
import {ModeToggle} from "@/components/mode-toggle.tsx";

const Orders = () => {
    return (
        <main className="w-full">
            <div className="w-full flex p-2 justify-between items-center">
                <div className="flex gap-4 items-center">
                    <SidebarTrigger/>
                    <Breadcrumb>
                        <BreadcrumbList>
                            <BreadcrumbItem>
                                <BreadcrumbLink asChild>
                                    <Link to="/orders">Orders</Link>
                                </BreadcrumbLink>
                            </BreadcrumbItem>

                        </BreadcrumbList>
                    </Breadcrumb>
                </div>
                <ModeToggle></ModeToggle>
            </div>
            <div className="flex p-2">


            </div>
        </main>

    );
};

export default Orders;