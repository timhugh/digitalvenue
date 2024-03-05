import React from "react";
import { useRouteError } from "react-router-dom";

export default () => {
    const error = useRouteError();
    console.error(error);

    return (
        <div
            style={{
                display: "flex",
                height: "100%",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
            }}
        >
            <h2>Oops!</h2>
            <p>Sorry, an unexpected error has occurred:</p>
            <p>
                <i>{error.statusText || error.message}</i>
            </p>
        </div>
    );
};
