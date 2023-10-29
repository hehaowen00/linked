import { Show, createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";

import {
  Alert,
  Button,
  Col,
  Container,
  Form,
  InputGroup,
} from "solid-bootstrap";
import Header from "../components/Header";
import api from "../lib/api";

export default function Register() {
  const navigate = useNavigate();

  let [registered, setRegistered] = createSignal(false);
  let [form, setForm] = createSignal({
    email: "",
    first: "",
    last: "",
  });
  let [passCode, setPassCode] = createSignal("");
  let [qrCode, setQrCode] = createSignal("");
  let [showError, setShowError] = createSignal(false);

  let updateForm = (e) => {
    setForm({ ...form(), [e.target.name]: e.target.value });
  };

  let submit = async (e) => {
    e.preventDefault();

    try {
      let res = await api.register(form());

      let json = await res.json();

      if (res.ok) {
        console.log(json);
        setRegistered(true);
        setQrCode(json.qr_code);
      } else {
        setShowError(true);
      }
    } catch (e) {
      setShowError(true);
    }
  };

  let onLogin = async (e) => {
    e.preventDefault();
    // let res = await fetch("https://localhost:8000/auth/login", {
    //   method: "POST",
    //   body: JSON.stringify({
    //     email: form().email,
    //     passCode: passCode(),
    //   }),
    // });

    let res = await api.login({
      email: form().email,
      passCode: passCode(),
    });

    if (res.ok) {
      navigate("/");
    }
  };

  return (
    <>
      <Header />
      <Container>
        <Col md={4} class="m-auto">
          <h1 class="mt-4 text-center">Register</h1>
          <Show when={registered()}>
            <Alert variant="success">
              <span>Registration Successful</span>
            </Alert>
            <Col>
              <h4>To Login</h4>
              <ol>
                <li>Download an Authenticator App</li>
                <li>Scan the QR Code</li>
                <li>
                  Type in the one time password (TOTP) within the time limit
                </li>
              </ol>
            </Col>
            <div class="w-full text-center">
              <img class="ml-auto" src={qrCode()} />
            </div>
            <div class="d-flex justify-content-center mt-4">
              <Form onSubmit={onLogin}>
                <Form.Group>
                  <InputGroup>
                    <Form.Control
                      type="text"
                      name="passCode"
                      placeholder="Pass Code"
                      value={passCode()}
                      onInput={(e) => setPassCode(e.target.value)}
                    />
                    <Button type="submit">Submit</Button>
                  </InputGroup>
                </Form.Group>
              </Form>
            </div>
          </Show>
          <Show when={!registered()}>
            <Show when={showError()}>
              <Alert variant="danger">Unable to create account</Alert>
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
                <Form.Label>First Name</Form.Label>
                <Form.Control
                  name="first"
                  type="text"
                  placeholder="First Name"
                  required
                  value={form().first}
                  onInput={updateForm}
                />
              </Form.Group>
              <Form.Group class="mb-3">
                <Form.Label>Last Name</Form.Label>
                <Form.Control
                  name="last"
                  type="text"
                  placeholder="Last Name"
                  required
                  value={form().last}
                  onInput={updateForm}
                />
              </Form.Group>
              <Form.Check
                class="mb-3"
                label="I accept the terms and conditions"
                required
              />
              <p class="mb-3">
                <a href="/login">Already have an account?</a>
              </p>
              <div class="w-full text-right">
                <Button type="submit" variant="primary">
                  Create Account
                </Button>
              </div>
            </Form>
          </Show>
        </Col>
      </Container>
    </>
  );
}
