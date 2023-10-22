import { Nav } from "solid-bootstrap";
import { createSignal } from "solid-js";
import { useLocation } from "@solidjs/router";

export default function LinkContainer({ href, children }) {
  const location = useLocation();
  const [isActive] = createSignal(location.pathname === href);

  return (
    <Nav.Link href={href} active={isActive()}>
      {children}
    </Nav.Link>
  );
}
