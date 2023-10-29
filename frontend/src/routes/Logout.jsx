import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";
import api from "../lib/api";

export default function Logout() {
  const navigate = useNavigate();
  createEffect(() => {
    let onLogout = async () => {
      let res = await api.logout();
      if (!res.ok) {
        return;
      }
      navigate("/");
    };
    onLogout();
  });
  return <></>;
}
