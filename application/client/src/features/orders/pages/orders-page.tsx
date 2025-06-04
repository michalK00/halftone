import {SidebarTrigger} from "@/components/ui/sidebar.tsx";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList
} from "@/components/ui/breadcrumb.tsx";
import {Link} from "react-router-dom";
import {ModeToggle} from "@/components/mode-toggle.tsx";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card";
import {LoadingSpinner} from "@/components/ui/loading-spinner";
import {useQuery} from "@tanstack/react-query";
import {getOrders} from "@/api/orders.ts";
import {Badge} from "@/components/ui/badge";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import {OrderDetailsDialog} from "@/features/orders/components/order-details-dialog";

const OrdersPage = () => {
    const ordersQuery = useQuery({
        queryKey: ['orders'],
        queryFn: getOrders,
    });

    const renderContent = () => {
        if (ordersQuery.status === 'pending') {
            return (
                <Card className="w-full">
                    <CardContent className="flex p-6 h-64 justify-center items-center">
                        <LoadingSpinner />
                    </CardContent>
                </Card>
            );
        }

        if (ordersQuery.status === 'error') {
            return (
                <Card className="w-full">
                    <CardContent className="p-6">
                        <div className="text-center text-red-500">
                            Error loading orders. Please try again later.
                        </div>
                    </CardContent>
                </Card>
            );
        }

        const orders = ordersQuery.data || [];

        if (orders.length === 0) {
            return (
                <Card className="w-full">
                    <CardContent className="p-6">
                        <div className="text-center text-muted-foreground">
                            No orders found.
                        </div>
                    </CardContent>
                </Card>
            );
        }

        return (
            <Card className="w-full">
                <CardHeader>
                    <CardTitle>Orders</CardTitle>
                    <CardDescription>
                        Manage all orders from your galleries
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Client Email</TableHead>
                                <TableHead>Photos</TableHead>
                                <TableHead>Status</TableHead>
                                <TableHead>Created</TableHead>
                                <TableHead>Comment</TableHead>
                                <TableHead>Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {orders.map((order) => (
                                <TableRow key={order.id}>
                                    <TableCell>{order.clientEmail}</TableCell>
                                    <TableCell>{order.photos.length} photos</TableCell>
                                    <TableCell>
                                        <Badge variant={order.status === 'completed' ? 'default' : 'secondary'}>
                                            {order.status}
                                        </Badge>
                                    </TableCell>
                                    <TableCell>{new Date(order.createdAt).toLocaleDateString()}</TableCell>
                                    <TableCell className="max-w-xs truncate">
                                        {order.comment || '-'}
                                    </TableCell>
                                    <TableCell>
                                        <OrderDetailsDialog order={order} />
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>
        );
    };

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
                <ModeToggle />
            </div>
            <div className="p-4">
                {renderContent()}
            </div>
        </main>
    );
};

export default OrdersPage;