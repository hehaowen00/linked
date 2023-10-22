import { Col, Container } from "solid-bootstrap";
import Header from "../components/Header";

export default function Home() {
  return (
    <>
      <Header />
      <Container>
        <Col md={6} className="mt-4 m-auto text-center">
          <h1>Linked</h1>
          <p>Bookmark Manager</p>
        </Col>
      </Container>
    </>
  );
}
