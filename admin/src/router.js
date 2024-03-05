import React from "react";
import { createBrowserRouter } from "react-router-dom";

import AdminPanel from "./AdminPanel";
import RouteError from "./route-error";
import OrderSearch from "./OrderSearch/OrderSearch";
import PrintOrder from "./PrintOrder";

const router = createBrowserRouter([
    {
        path: "/",
        element: <AdminPanel />,
        errorElement: <RouteError />,
        children: [
            {
                path: "orders",
                element: <OrderSearch />,
            },
        ],
    },
    {
        path: "/orders/:id/print",
        element: <PrintOrder />,
    },
]);

export default router;
