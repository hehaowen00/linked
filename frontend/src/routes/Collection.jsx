import { createEffect, createSignal } from "solid-js";
import { useParams } from "@solidjs/router";

import { Button, Col, Container, Form, Row } from "solid-bootstrap";
import Header from "../components/Header";
import Bookmarks from "../components/Bookmarks";

export default () => {
  let [collection, setCollection] = createSignal({ name: "" });
  let [edit, setEdit] = createSignal(false);
  let [changed, setChanged] = createSignal({
    name: "",
  });
  let params = useParams();

  createEffect(() => {
    let getCollection = async () => {
      let res = await fetch(
        `https://localhost:8000/api/collections/${params.collection}`,
        {
          credentials: "include",
        },
      );
      let json = await res.json();
      if (!res.ok) {
        return;
      }
      setCollection(json.data);
    };
    getCollection();
  });

  let fetchItems = async () => {
    let res = await fetch(
      `https://localhost:8000/api/items/${params.collection}`,
      {
        credentials: "include",
      },
    );
    return res;
  };

  let addItem = async (item) => {
    let res = await fetch(`https://localhost:8000/api/items`, {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({
        ...item,
        collection_id: params.collection,
      }),
    });
    return res;
  };

  let editCollection = () => {
    setChanged({ name: collection().name });
    setEdit(true);
  };

  let updateCollection = async () => {
    let res = await fetch(
      `https://localhost:8000/api/collections/${params.collection}`,
      {
        method: "PUT",
        credentials: "include",
        body: JSON.stringify({
          ...collection(),
          name: changed().name,
        }),
      },
    );
    if (!res.ok) {
      return;
    }
    setEdit(false);
  };

  return (
    <>
      <Header authenticated={true} />
      <Container class="mt-2 content">
        <Show when={!edit()}>
          <Row>
            <Col md={8}>
              <Button variant="light" onClick={editCollection}>
                {collection().name}
              </Button>
            </Col>
          </Row>
          <Bookmarks fetchItems={fetchItems} addItem={addItem} />
        </Show>
        <Show when={edit()}>
          <Row>
            <Col>
              <Form onSubmit={updateCollection}>
                <Row>
                  <Col>
                    <Form.Control
                      name="name"
                      type="text"
                      placeholder="Name"
                      required
                      value={changed().name}
                    />
                  </Col>
                </Row>
                <Row class="mt-2">
                  <Col class="text-right spaced-left">
                    <Button type="submit" size="sm">
                      Save
                    </Button>
                    <Button size="sm" onClick={() => setEdit(false)}>
                      Cancel
                    </Button>
                  </Col>
                </Row>
              </Form>
            </Col>
          </Row>
        </Show>
      </Container>
    </>
  );
};
