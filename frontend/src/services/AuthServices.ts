const URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

interface LoginReponse {
    message: string;
    token: string;
}

interface ErrorReponse {
    error: string;
    details?: string
}

export async function loginAdmin(email: string, password: string):Promise<Boolean> {

    try {
        const response = await fetch(`${URL}/login`, {
            method: "POST", 
            credentials: "include", 
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({ email, password }),
        });

        if (!response.ok) {
            const errorData: ErrorReponse = await response.json();
            console.error("erro de login: ", errorData.error);
            return false;
        }

        const data: LoginReponse = await response.json();
        localStorage.setItem("auth_token", data.token);

        return true;
    } catch (e) {
        console.error("erro na req de login: ", e);
        return false;
    }
}

export async function registerAdmin(
    email: string,
    password: string
): Promise<{
     success: boolean,
     message: string
}> {
    try {
        const response = await fetch(`${URL}/register`, {
            method: "POST", 
            headers: {"Content-Type":"application/json"},
            body: JSON.stringify({ email, password }),
        });

        const data = await response.json(); 

        if (!response.ok) {
            return {
                success: false, 
                message: data.error || "erro ao registrar adm"
            };
        }

        return {
            success: true, 
            message: "adm registrado com sucesso"
        };
    } catch (e) {
        console.error("erro na req de registro: ", e);
        return {
            success: false, 
            message: "erro de rede ao tentar registrar"
        };
    }
}

export default async function logoutAdmin(): Promise<boolean> {
    try {
        const response = await fetch(`${URL}/logout`, {
            method: "POST", 
            credentials: "include",
        });

        localStorage.removeItem("auth_token");

        return response.ok;
    } catch (e) {
        console.error("erro na requisição de logout: ", e);
        return false;
    }
}

export function isAuthenticated(): boolean {
    return !!localStorage.getItem("auth_token");
}

export function getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem("auth_tokem");
    return {
        "Content-Type":"application/json",
        "Authorization": token ? `Bearer ${token}` : ""
    };
}