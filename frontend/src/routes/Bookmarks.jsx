import { Container } from "solid-bootstrap";
import Header from "../components/Header";
import Bookmarks from "../components/Bookmarks";

export default () => {
  let fetchItems = async () => {
    let res = await fetch("https://localhost:8000/api/items", {
      credentials: "include",
    });
    if (!res.ok) {
      return;
    }
    return res;
  };

  let addItem = async (item) => {
    let res = await fetch("https://localhost:8000/api/items", {
      method: "POST",
      credentials: "include",
      body: JSON.stringify(item),
    });
    return res;
  };

  return (
    <>
      <Header authenticated={true} />
      <Container class="mt-2 content">
        <Bookmarks fetchItems={fetchItems} addItem={addItem} />
      </Container>
    </>
  );
};
