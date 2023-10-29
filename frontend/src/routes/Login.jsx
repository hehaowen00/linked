import { Alert, Button, Col, Container, Form, Row } from "solid-bootstrap";
import Header from "../components/Header";
import { createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";
import api from "../lib/api";

export default function Home() {
  const navigate = useNavigate();

  let [showError, setShowError] = createSignal(false);
  let [form, setForm] = createSignal({
    email: "",
    passCode: "",
  });

  let updateForm = (e) => {
    setForm({ ...form(), [e.target.name]: e.target.value });
  };

  let submit = async (e) => {
    e.preventDefault();

    try {
      let res = await api.login(form());

      if (res.ok) {
        navigate("/");
      } else {
        setShowError(true);
      }
    } catch (e) {
      setShowError(true);
    }
  };

  return (
    <>
      <Header />
      <Container>
        <Row>
          <Col md={4} class="mx-auto">
            <div class="text-center">
              <h1 class="mt-4 mx-auto">Login</h1>
            </div>
            <Show when={showError()}>
              <Alert variant="danger">Unable to authenticate</Alert>
            </Show>
            <Form onSubmit={submit}>
              <Form.Group class="mb-3">
                <Form.Label>Email Address</Form.Label>
                <Form.Control
                  name="email"
                  type="email"
                  placeholder="Email Address"
                  required
                  value={form().email}
                  onInput={updateForm}
                />
              </Form.Group>
              <Form.Group class="mb-3">
                <Form.Label>Pass Code</Form.Label>
                <Form.Control
                  name="passCode"
                  type="text"
                  placeholder="Pass Code"
                  maxLength={6}
                  required
                  value={form().passCode}
                  onInput={updateForm}
                />
              </Form.Group>
              <a class="mt-1 mb-1" href="/register">
                <p>Don't have an account?</p>
              </a>
              <div class="w-full text-right">
                <Button variant="primary" type="submit" value="submit">
                  Login
                </Button>
              </div>
            </Form>
          </Col>
        </Row>
      </Container>
    </>
  );
}
