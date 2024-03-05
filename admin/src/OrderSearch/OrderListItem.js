import React, { useState } from "react";
import { Button, Col, Collapse, Row } from "react-bootstrap";
import { format, parseISO } from "date-fns";

const OrderListItem = ({ order }) => {
    const [collapsed, setCollapsed] = useState(false);

    const toggleCollapsed = () => {
        setCollapsed(!collapsed);
    };

    const ticketEmailLink = `${process.env.BACKEND}/orders/${order.code}/email`;

    const openTicketEmail = () => {
        window.open(ticketEmailLink, "_blank");
    };

    const openPrintOrder = () => {
        window.open(`/orders/${order.id}/print`);
    };

    return (
        <>
            <tr key={order.id} onClick={toggleCollapsed}>
                <td>{order.code}</td>
                <td>{order.customerName}</td>
                <td>{order.customerEmail}</td>
                <td>{format(parseISO(order.date), "eee M/d/yyyy h:mm b")}</td>
            </tr>
            <Collapse in={collapsed} style={{ backgroundColor: "lightgray" }}>
                <tr>
                    <td colSpan="100%">
                        <Row>
                            <Col sm="8">
                                Tickets:
                                <ul>
                                    {order.items.map((item) => (
                                        <li key="item.ID">{item.Name}</li>
                                    ))}
                                </ul>
                            </Col>
                            <Col
                                sm="4"
                                style={{
                                    display: "flex",
                                    flexDirection: "column",
                                    justifyContent: "space-between",
                                    alignItems: "flex-end",
                                }}
                            >
                                <Button onClick={openPrintOrder} size="sm">
                                    Print Tickets
                                </Button>
                                <Button onClick={openTicketEmail} size="sm">
                                    View Ticket Email
                                </Button>
                            </Col>
                        </Row>
                    </td>
                </tr>
            </Collapse>
        </>
    );
};

export default OrderListItem;
