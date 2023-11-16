import { Container, Nav, Navbar } from "solid-bootstrap";
import LinkContainer from "./LinkContainer";
import { Show } from "solid-js";

export default function Header({ authenticated }) {
  return (
    <Navbar expand="lg" collapseOnSelect sticky="top">
      <Container>
        <Navbar.Brand href="/">Linked</Navbar.Brand>
        <Navbar.Toggle />
        <Navbar.Collapse>
          <Nav class="me-auto">
            <Show when={authenticated}>
              <LinkContainer href="/bookmarks">Bookmarks</LinkContainer>
              <LinkContainer href="/collections">Collections</LinkContainer>
            </Show>
          </Nav>
          <Show when={authenticated}>
            <Nav class="ml-auto">
              <LinkContainer href="/logout">Logout</LinkContainer>
            </Nav>
          </Show>
          <Show when={!authenticated}>
            <Nav class="ml-auto">
              <LinkContainer href="/login">Login</LinkContainer>
              <LinkContainer href="/register">Register</LinkContainer>
            </Nav>
          </Show>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}
