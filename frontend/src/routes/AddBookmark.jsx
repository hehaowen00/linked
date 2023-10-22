import { Col, Container, Button, Form, Row } from "solid-bootstrap";
import Header from "../components/Header";

export default function AddBookmark() {
  let pasteUrl = () => {};
  let addBookmark = () => {};
  return (
    <>
      <Header authenticated={true} />
      <Container class="mt-2">
        <Row>
          <Col>
            <h5>Add Bookmark</h5>
          </Col>
        </Row>
        <Row>
          <Col>
            <Row>
              <Col class="mt-1">
                <Form.Control type="text" size="sm" placeholder="URL" />
              </Col>
            </Row>
            <Row class="mt-1">
              <Col>
                <Form.Control type="text" size="sm" placeholder="Title" />
              </Col>
            </Row>
            <Row class="mt-1">
              <Col>
                <Form.Control type="text" size="sm" placeholder="Description" />
              </Col>
            </Row>
            <Row class="mt-2">
              <Col class="text-right spaced-left">
                <Button size="sm" onClick={pasteUrl}>
                  Paste URL
                </Button>
                <Button size="sm" onClick={addBookmark}>
                  Add Bookmark
                </Button>
              </Col>
            </Row>
          </Col>
        </Row>
        <Row></Row>
      </Container>
    </>
  );
}
