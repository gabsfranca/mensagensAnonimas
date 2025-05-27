import { createEffect } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { isAuthenticated } from "../services/AuthServices";

interface AuthGuardProps {
    children: any;
}

export default function AuthGuard(props: AuthGuardProps) {
    const navigate = useNavigate();

    createEffect(() => {
        if (!isAuthenticated()) {
            navigate('/login', { replace: true });
        }
    });

    return (
    <>
      {isAuthenticated() && props.children}
    </>
  );
}