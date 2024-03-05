export default (() => {
    const mapOrders = (orders) => {
        var mappedOrders = orders.map((order) => {
            return {
                id: order?.ID,
                code: order?.Code,
                customerName: `${order?.Customer?.FirstName} ${order?.Customer?.LastName}`,
                customerEmail: order?.Customer?.Email,
                date: order?.CreatedAt,
                items: order?.OrderItems,
            };
        });
        return mappedOrders;
    };

    const expandPath = (path) => {
        return `${process.env.BACKEND}/${path}`;
    };

    return {
        searchOrders(searchTerm) {
            return fetch(expandPath(`orders?q=${searchTerm}`))
                .then((response) => response.json())
                .then((orders) => orders || [])
                .then((orders) =>
                    orders.filter(
                        (order) => order !== null && order !== undefined
                    )
                )
                .then((orders) => mapOrders(orders));
        },

        getOrder(orderId) {
            return fetch(expandPath(`orders/${orderId}`)).then((response) =>
                response.json()
            );
        },
    };
})();
