import { Button, Col, Container, Form, Row } from "solid-bootstrap";
import Header from "../components/Header";
import Bookmarks from "../components/Bookmarks";

import { createEffect, createSignal } from "solid-js";
import { useNavigate, useParams } from "@solidjs/router";
import api from "../lib/api";

export default () => {
  const navigate = useNavigate();
  const params = useParams();

  let [collection, setCollection] = createSignal({
    name: "",
    created_at: 0,
    deleted_at: 0,
  });
  let [loaded, setLoaded] = createSignal(false);

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

  return (
    <>
      <Header authenticated={true} />
      <Show when={loaded()}>
        <Container>
          <Row class="mt-2">
            <Col>
              <Button size="sm" variant="light">
                {collection().name}
              </Button>
            </Col>
          </Row>
        </Container>
        <Container class="mt-1 content overflow">
          <Bookmarks
            archived={() => collection().deleted_at > 0}
            fetchItems={fetchItems}
            addItem={addItem}
          />
        </Container>
      </Show>
    </>
  );
};
