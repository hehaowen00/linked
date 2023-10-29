import { Container } from "solid-bootstrap";
import Header from "../components/Header";
import Bookmarks from "../components/Bookmarks";
import api from "../lib/api";

export default () => {
  return (
    <>
      <Header authenticated={true} />
      <Container class="mt-2 content">
        <Bookmarks
          archived={() => false}
          fetchItems={api.getItems}
          addItem={api.addItem}
        />
      </Container>
    </>
  );
};
