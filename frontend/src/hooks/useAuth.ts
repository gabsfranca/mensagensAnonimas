import { createSignal, onMount } from 'solid-js';

export const useAuth = () => {
    const [isAuthenticated, setIsAuthenticated] = createSignal(false);
    const [isLoading, setIsLoading] = createSignal(false);
    const [error, setError] = createSignal('');

    onMount(() => {
        const token = localStorage.getItem("auth_token");
        if (!token || token.trim() === "") {
            setError('Usuário não encontrado');
            setIsLoading(false);
            setTimeout(() => {
                window.location.href = '/login';
            }, 1000);
            return;
        }

        setIsAuthenticated(true);
        setIsLoading(false);
    });

    const logout = () => {
        localStorage.removeItem("auth_token");
        window.location.href = '/login';
    };

    return {
        isAuthenticated,
        isLoading, 
        error, 
        logout
    };
};