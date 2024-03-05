import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import APIService from "./APIService";

const PrintOrder = () => {
    const { id } = useParams();
    const [order, setOrder] = useState({});

    useEffect(async () => {
        console.log("fetching order");
        const order = await APIService.getOrder(id);
        console.log("order found", order);
        setOrder(order);
    }, []);

    return (
        <>
            {order ? (
                order.OrderItems?.map((item) => {
                    return (
                        <div
                            style={{
                                width: "50%",
                                height: "50%",
                                margin: "0px",
                                border: "solid 1px black",
                                textAlign: "center",
                                display: "inline-block",
                            }}
                        >
                            <h2 style={{ fontSize: "24px" }}>{item.Name}</h2>
                            <img
                                src={`${process.env.BACKEND}/tickets/${item.Ticket.Code}/qrcode`}
                            ></img>
                            <p>{item.Ticket.Code}</p>
                        </div>
                    );
                })
            ) : (
                <>No order found</>
            )}
        </>
    );
};

export default PrintOrder;
