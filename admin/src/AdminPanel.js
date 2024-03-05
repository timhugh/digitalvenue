import React from "react";
import { Container, Nav, Navbar } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";
import { Outlet } from "react-router-dom";

const AdminPanel = () => {
    return (
        <>
            <Navbar expand="md">
                <Container>
                    <LinkContainer to="/">
                        <Navbar.Brand>Box Office</Navbar.Brand>
                    </LinkContainer>
                    <Navbar.Toggle />
                    <Navbar.Collapse>
                        <Nav>
                            <LinkContainer to="/orders">
                                <Nav.Link>Order Search</Nav.Link>
                            </LinkContainer>
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>
            <Container>
                <Outlet />
            </Container>
        </>
    );
};

export default AdminPanel;
