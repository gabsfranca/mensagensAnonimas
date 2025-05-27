import { useNavigate } from "@solidjs/router";
import logoutAdmin from "../services/AuthServices";

export default function AdminNavbar() {
  const navigate = useNavigate();

  const handleLogout = async () => {
    const success = await logoutAdmin();
    if (success) {
      navigate('/login', { replace: true });
    }
  };

  return (
    <nav class="bg-indigo-600 text-white p-4">
      <div class="container mx-auto flex justify-between items-center">
        <div class="flex items-center space-x-4">
          <h1 class="text-xl font-bold">Painel Administrativo</h1>
          <div class="hidden md:flex space-x-4">
            <button 
              onClick={() => navigate('/admin')}
              class="hover:text-indigo-200 transition-colors"
            >
              Dashboard
            </button>
            <button 
              onClick={() => navigate('/admin/messages')}
              class="hover:text-indigo-200 transition-colors"
            >
              Mensagens
            </button>
          </div>
        </div>
        <button 
          onClick={handleLogout}
          class="bg-indigo-700 hover:bg-indigo-800 px-4 py-2 rounded transition-colors"
        >
          Sair
        </button>
      </div>
    </nav>
  );
}