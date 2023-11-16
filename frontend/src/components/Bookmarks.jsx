import { Alert, Button, Col, Form, Offcanvas, Row } from "solid-bootstrap";
import { For, Show, createEffect, createSignal, on } from "solid-js";
import BookmarkItem from "./BookmarkItem";

import { isValidURL } from "../lib/utils";
import { useNavigate } from "@solidjs/router";
import api from "../lib/api";

export default ({ archived, fetchItems, addItem }) => {
  const navigate = useNavigate();

  let [bookmarks, setBookmarks] = createSignal([]);
  let [form, setForm] = createSignal({
    title: "",
  });
  let [url, setURL] = createSignal("");
  let [selected, setSelected] = createSignal({});

  let [showAlert, setAlert] = createSignal("");
  let [showSelected, setShowSelected] = createSignal(false);
  let [loading, setLoading] = createSignal(false);

  let updateForm = (e) => {
    setForm({ ...form(), [e.target.name]: e.target.value });
  };

  let reloadItems = async () => {
    let res = await fetchItems();
    if (!res.ok) {
      return;
    }
    setBookmarks(await res.json());
  };

  let onAddItem = async (e) => {
    e.preventDefault();

    let res = await addItem({ ...form(), url: url() });

    if (!res.ok) {
      setAlert("error");
    } else {
      setAlert("success");
      setURL("");
      setForm({
        title: "",
      });
      await reloadItems();
    }
  };

  let pasteUrl = async () => {
    let text = await navigator.clipboard.readText();
    if (isValidURL(text)) {
      setForm({ ...form(), url: text });
      setURL(text);
    }
  };

  let onDeleteItem = async (item) => {
    let res = await api.deleteItem(item);
    if (!res.ok) {
      return;
    }
    setSelected({});
    setShowSelected(false);
    await reloadItems();
  };

  let selectItem = (item) => {
    setSelected(item);
    setShowSelected(true);
  };

  createEffect(
    on(url, (u) => {
      let fetchInfo = async () => {
        if (!u) {
          return;
        }
        setLoading(true);
        setForm({
          title: "",
        });
        try {
          let res = await api.getPageInfo(u);
          if (!res.ok) {
            setLoading(false);
            return;
          }
          let json = await res.json();
          setForm({
            title: json.title,
            desc: "",
          });
        } catch (e) {}
        setLoading(false);
      };
      fetchInfo();
    }),
  );

  createEffect(() => {
    let f = async () => {
      let res = await fetchItems();
      let json = await res.json();

      if (!res.ok) {
        return;
      }

      setBookmarks(json);
    };
    f();
  });

  return (
    <>
      <Show when={!archived()}>
        <Col class="mt-1 no-overflow" md={4}>
          <Form onSubmit={onAddItem}>
            <Row>
              <Col>
                <Form.Control
                  name="url"
                  type="text"
                  size="sm"
                  placeholder="URL"
                  required
                  value={url()}
                  onInput={(e) => setURL(e.target.value)}
                  disabled={loading()}
                />
              </Col>
            </Row>
            <Row class="mt-1">
              <Col>
                <Form.Control
                  name="title"
                  type="text"
                  size="sm"
                  placeholder={loading() ? "Loading..." : "Title"}
                  required
                  value={form().title}
                  onInput={updateForm}
                  disabled={loading()}
                />
              </Col>
            </Row>
            <Row class="mt-1 mb-2">
              <Col class="text-right spaced-left">
                <Button size="sm" onClick={pasteUrl}>
                  <svg
                    fill="currentColor"
                    stroke-width="0"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 16 16"
                    style="overflow: visible; color: currentcolor; margin-top: -4px;"
                    height="1em"
                    width="1em"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M10 1.5a.5.5 0 0 0-.5-.5h-3a.5.5 0 0 0-.5.5v1a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-1Zm-5 0A1.5 1.5 0 0 1 6.5 0h3A1.5 1.5 0 0 1 11 1.5v1A1.5 1.5 0 0 1 9.5 4h-3A1.5 1.5 0 0 1 5 2.5v-1Zm-2 0h1v1A2.5 2.5 0 0 0 6.5 5h3A2.5 2.5 0 0 0 12 2.5v-1h1a2 2 0 0 1 2 2V14a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V3.5a2 2 0 0 1 2-2Z"
                    ></path>
                  </svg>
                </Button>
                <Button size="sm" type="submit">
                  Add Bookmark
                </Button>
              </Col>
            </Row>
            <Row class="mt-1">
              <Col>
                <Alert
                  variant="danger"
                  dismissible
                  show={showAlert() == "error"}
                  onClose={() => setAlert("")}
                >
                  Unable to add bookmark
                </Alert>
                <Alert
                  variant="success"
                  dismissible
                  show={showAlert() == "success"}
                  onClose={() => setAlert("")}
                >
                  Bookmark added
                </Alert>
              </Col>
            </Row>
          </Form>
        </Col>
      </Show>
      <Col class="mt-1 scrollable no-scrollbar">
        <For
          each={bookmarks()}
          fallback={
            <Alert variant="dark" class="text-center">
              No Bookmarks Found
            </Alert>
          }
        >
          {(item, index) => (
            <BookmarkItem
              index={index}
              item={item}
              select={selectItem}
              deleteItem={onDeleteItem}
            />
          )}
        </For>
      </Col>
      <Offcanvas
        placement="bottom"
        show={showSelected()}
        onHide={() => setShowSelected(!showSelected())}
      >
        <Offcanvas.Header closeButton>
          <Col>{selected().title}</Col>
        </Offcanvas.Header>
        <Offcanvas.Body>
          <Row>
            <Button
              class="view-btn"
              size="sm"
              onClick={() => window.open(selected().url, "_blank")}
            >
              Open
            </Button>
          </Row>
          <Row class="mt-2">
            <Button class="view-btn" size="sm" onClick={() => {}}>
              View
            </Button>
          </Row>
          <Row class="mt-4 mb-4">
            <Button
              class="remove-btn"
              size="sm"
              variant="dark"
              onClick={() => onDeleteItem(selected())}
            >
              Delete
            </Button>
          </Row>
        </Offcanvas.Body>
      </Offcanvas>
    </>
  );
};
