import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";

export default function Logout() {
  const navigate = useNavigate();
  createEffect(() => {
    let logout = async () => {
      let res = await fetch("https://localhost:8000/auth/logout");
      if (!res.ok) {
        return;
      }
      navigate("/");
    };
    logout();
  });
  return <></>;
}
