import React, { useState } from "react";
import { Col, Form, FormControl, Row, Button, Table } from "react-bootstrap";

import APIService from "../APIService.js";
import OrderList from "./OrderList.js";

const OrderSearch = () => {
    const [orders, setOrders] = useState([]);

    const handleSearch = (e) => {
        e.preventDefault();

        const formData = new FormData(e.target);
        const search = Object.fromEntries(formData.entries());
        const searchTerm = search.search;

        APIService.searchOrders(searchTerm).then((orders) => {
            setOrders(orders || []);
        });
    };

    return (
        <>
            <Row>
                <Col>
                    <Form onSubmit={handleSearch} className="d-flex">
                        <FormControl
                            name="search"
                            type="text"
                            placeholder="Search by customer name or order code"
                        />
                        <Button variant="primary" type="submit">
                            Search
                        </Button>
                    </Form>
                </Col>
            </Row>
            <Row>
                <Col>
                    <h2>Results</h2>
                    <OrderList orders={orders}></OrderList>
                </Col>
            </Row>
        </>
    );
};

export default OrderSearch;
