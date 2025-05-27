import { createSignal } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { loginAdmin } from "../services/AuthServices";

export default function Login() {
    const [email, setEmail] = createSignal('');
    const [password, setPassword] = createSignal('');
    const [error, setError] = createSignal('');
    const [isLoading, setIsLoading] = createSignal(false);

    const navigate = useNavigate();

    const handleSubmit = async (e: Event) => {
        e.preventDefault();
        setError('');

        if (!email() || !password()) {
            setError('Por favor, preencha todos os campos')
            return;
        }

        if (password().length < 8) {
            setError('A senha deve ter no minimo 8 carcteres');
            return;
        }

        setIsLoading(true);

        try {
            const success = await loginAdmin(email(), password());
            if (success) navigate('/admin');
            else setError('credenciais inválidas. Verifique seu email e senha');
        } catch (e) {
            console.error('erro ao fazer login: ', e);
            setError('ocorreu um erro ao tentar fazer login, tente novamente mais tarde, ou fale com o gabriel')
        } finally {
            setIsLoading(false);            
        }
    };

    return (
        <div class="w-full max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow-md">
            <h1 class="text-2xl font-bold mb-6 text-center">Login de Administrador</h1>
            
            {error() && (
                <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4" role="alert">
                    <span class="block sm:inline">{error()}</span>
                </div>
            )}
            
            <form onSubmit={handleSubmit} class="space-y-4">
                <div>
                    <label for="email" class="block text-sm font-medium text-gray-700">
                        Email
                    </label>
                    <input
                        id="email"
                        type="email"
                        value={email()}
                        onInput={(e) => setEmail(e.currentTarget.value)}
                        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        placeholder="Digite seu email"
                        required
                    />
                </div>
                
                <div>
                    <label for="password" class="block text-sm font-medium text-gray-700">
                        Senha
                    </label>
                    <input
                        id="password"
                        type="password"
                        value={password()}
                        onInput={(e) => setPassword(e.currentTarget.value)}
                        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                        placeholder="Digite sua senha"
                        minLength={8}
                        required
                    />
                </div>
                
                <div>
                    <button
                        type="submit"
                        disabled={isLoading()}
                        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:bg-indigo-300"
                    >
                        {isLoading() ? 'Entrando...' : 'Entrar'}
                    </button>
                </div>
            </form>
        </div>
    );
}