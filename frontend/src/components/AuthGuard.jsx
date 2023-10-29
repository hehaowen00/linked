import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";
import api from "../lib/api";

export default ({ children }) => {
  const navigate = useNavigate();

  createEffect(() => {
    let checkAuth = async () => {
      let res = await api.validate();
      if (!res.ok) {
        navigate("/login");
      }
    };

    checkAuth();
  });

  return <>{children}</>;
};
