import { Button, Card } from "solid-bootstrap";
import { Show } from "solid-js";

export default function BookmarkItem({ index, item, deleteItem, select }) {
  if (mobileAndTabletCheck()) {
    return <MobileComponent index={index} item={item} select={select} />;
  } else {
    return (
      <div class="bookmark-item">
        <Card>
          <Card.Body>
            <a href={item.url} target="_blank">
              {item.title}
            </a>
            <br />
            <span class="text-sm">
              Created At: {new Date(item.created_at).toLocaleString()}{" "}
            </span>
            <br />
            <div class="flex-row fifty text-right spaced-left mt-1">
              <Button class="xs view-btn" size="sm">
                View
              </Button>
              <Button
                class="xs remove-btn"
                size="sm"
                variant="light"
                onClick={() => deleteItem(item)}
              >
                Delete
              </Button>
            </div>
          </Card.Body>
        </Card>
      </div>
    );
  }
}

function MobileComponent({ index, item, select }) {
  return (
    <>
      <div class="bookmark-item cursor" onClick={() => select(item)}>
        <Card>
          <Card.Body>
            <a
              href={item.url}
              target="_blank"
              onClick={(e) => e.preventDefault()}
            >
              {item.title}
            </a>
            <p />
            <span class="mt-3 mb-2 text-sm">
              Created At: {new Date(item.created_at).toLocaleString()}{" "}
            </span>
          </Card.Body>
        </Card>
      </div>
    </>
  );
}
