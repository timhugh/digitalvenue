import React from "react";
import { Table } from "react-bootstrap";
import OrderListItem from "./OrderListItem";

const OrderList = ({ orders }) => {
    return (
        <Table>
            <thead>
                <tr>
                    <th>Order Code</th>
                    <th>Customer Name</th>
                    <th>Customer Email</th>
                    <th>Date</th>
                </tr>
            </thead>
            <tbody>
                {orders.map((order) => (
                    <OrderListItem order={order}></OrderListItem>
                ))}
                {orders.length === 0 && (
                    <tr>
                        <td colSpan="4">No orders found</td>
                    </tr>
                )}
            </tbody>
        </Table>
    );
};

export default OrderList;
