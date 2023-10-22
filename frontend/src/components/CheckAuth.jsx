import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";

export default ({ children }) => {
  const navigate = useNavigate();

  createEffect(() => {
    let checkAuth = async () => {
      let res = await fetch("https://localhost:8000/auth/validate", {
        credentials: "include",
      });
      if (res.ok) {
        navigate("/bookmarks");
      }
    };
    checkAuth();
  });

  return <>{children}</>;
};
