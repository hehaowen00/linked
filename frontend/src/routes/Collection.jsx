import { createEffect, createSignal } from "solid-js";
import { useNavigate, useParams } from "@solidjs/router";
import api from "../lib/api";

import { Button, Col, Container, Form, Row } from "solid-bootstrap";
import Header from "../components/Header";
import Bookmarks from "../components/Bookmarks";

export default () => {
  const navigate = useNavigate();

  let [collection, setCollection] = createSignal({
    name: "",
    created_at: 0,
    deleted_at: 0,
  });
  let [loaded, setLoaded] = createSignal(false);
  let [edit, setEdit] = createSignal(false);
  let [changed, setChanged] = createSignal({
    name: "",
  });
  let params = useParams();

  createEffect(() => {
    let getCollection = async () => {
      let res = await api.getCollection(params.collection);
      let json = await res.json();
      if (!res.ok) {
        return;
      }
      setCollection(json.data);
      setLoaded(true);
    };
    getCollection();
  });

  let fetchItems = async () => {
    let res = await api.getCollectionItems(params.collection);
    return res;
  };

  let addItem = async (item) => {
    let res = await api.addItem({
      ...item,
      collection_id: params.collection,
    });
    return res;
  };

  let editCollection = () => {
    setChanged({ name: collection().name });
    setEdit(true);
  };

  let updateCollection = async () => {
    let res = await api.updateCollection({
      ...collection(),
      name: changed().name,
    });
    if (!res.ok) {
      return;
    }
    setEdit(false);
  };

  let deleteCollection = async () => {
    let res = await api.deleteCollection(collection());
    if (!res.ok) {
      return;
    }
    if (collection().deleted_at === 0) {
      setEdit(false);
      let json = await res.json();
      setCollection({ ...collection(), deleted_at: 0 });
      // setCollection(json.data, { equals: false });
    } else {
      navigate("/collections");
    }
  };

  let unarchiveCollection = async () => {
    setCollection({ ...collection(), deleted_at: 0 }, { equals: false });
    await updateCollection();
  };

  return (
    <>
      <Header authenticated={true} />
      <Show when={loaded()}>
        <Container class="mt-2 content">
          <Show when={!edit()}>
            <Row>
              <Col md={8}>
                <Button variant="light" onClick={editCollection}>
                  {collection().name}
                </Button>
              </Col>
            </Row>
            <Bookmarks
              archived={() => collection().deleted_at > 0}
              fetchItems={fetchItems}
              addItem={addItem}
            />
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
                    <Col>
                      <Button
                        type="submit"
                        size="sm"
                        class="w-full"
                        onClick={updateCollection}
                      >
                        Save
                      </Button>
                    </Col>
                    <Col>
                      <Button
                        size="sm"
                        variant="dark"
                        class="w-full"
                        onClick={() => setEdit(false)}
                      >
                        Cancel
                      </Button>
                    </Col>
                  </Row>
                  <Row>
                    <Col>
                      <Button
                        size="sm"
                        variant="danger"
                        onClick={deleteCollection}
                      >
                        Archive Collection
                      </Button>
                    </Col>
                    <Col>
                      <Button
                        size="sm"
                        variant="danger"
                        onClick={unarchiveCollection}
                      >
                        Unarchive
                      </Button>
                    </Col>
                  </Row>
                </Form>
              </Col>
            </Row>
          </Show>
        </Container>
      </Show>
    </>
  );
};
