import {
  Alert,
  Container,
  Form,
  Row,
  Col,
  Button,
  Card,
} from "solid-bootstrap";
import Header from "../components/Header";
import { createEffect, createSignal, For } from "solid-js";

export default function Collections() {
  let [showAlert, setAlert] = createSignal("");
  let [collection, setCollection] = createSignal({
    name: "",
  });
  let [collections, setCollections] = createSignal([]);

  let updateCollection = (e) => {
    setCollection({ ...collection(), [e.target.name]: e.target.value });
  };

  let getCollections = async () => {
    let res = await fetch("https://localhost:8000/api/collections", {
      credentials: "include",
    });
    if (!res.ok) {
    }
    let json = await res.json();
    setCollections(json.data);
  };
  createEffect(() => {
    getCollections();
  });

  let addCollection = async (e) => {
    e.preventDefault();

    let res = await fetch("https://localhost:8000/api/collections", {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({
        ...collection(),
      }),
    });

    if (!res.ok) {
      setAlert("error");
      return;
    }

    setAlert("success");
    setCollection({ name: "" });
    await getCollections();
  };

  return (
    <>
      <Header authenticated={true} />
      <Container class="mt-2 content">
        <Row class="reversed flexed">
          <Col class="mt-2" md={4}>
            <Form onSubmit={addCollection}>
              <Row>
                <Col>
                  <Form.Control
                    name="name"
                    type="text"
                    size="sm"
                    placeholder="New Collection"
                    required
                    onInput={updateCollection}
                    value={collection().name}
                  />
                </Col>
              </Row>
              <Row class="mt-1">
                <Col class="text-right">
                  <Button size="sm" type="submit">
                    Add Collection
                  </Button>
                </Col>
              </Row>
              <Row class="mt-2">
                <Col>
                  <Show when={showAlert() === "error"}>
                    <Alert
                      variant="danger"
                      dismissible
                      onClose={() => setAlert("")}
                    >
                      Unable to create collection
                    </Alert>
                  </Show>
                  <Show
                    when={showAlert() === "success"}
                    onClose={() => setAlert("")}
                  >
                    <Alert variant="success" dismissible>
                      Collection added
                    </Alert>
                  </Show>
                </Col>
              </Row>
            </Form>
          </Col>
          <Col>
            <Row class="mt-2 mb-4">
              <For each={collections()}>
                {(collection, index) => (
                  <div class="bookmark-item">
                    <Card>
                      <Card.Body>
                        <a href={`/collections/${collection.id}`}>
                          {collection.name}
                        </a>
                        <br />
                        <span>
                          Created At:{" "}
                          {new Date(collection.created_at).toLocaleString()}{" "}
                        </span>
                        <Show when={collection.archived_at > 0}>
                          <br />
                          <span>
                            Archived At:{" "}
                            {new Date(collection.archived_at).toLocaleString()}
                          </span>
                        </Show>
                      </Card.Body>
                    </Card>
                  </div>
                )}
              </For>
            </Row>
          </Col>
        </Row>
      </Container>
    </>
  );
}
