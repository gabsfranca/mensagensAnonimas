import { For, createSignal, onMount, Show, createEffect } from 'solid-js';
import { MessageResponse, MessageStatus } from '../types';
import { Spinner } from './spinner';
import { MediaViewer } from './MediaViewer'; // ADICIONE ESTA IMPORTAÇÃO
import { fetchMessages, updateMessageStatus, addMessageObs } from '../services/AdminServices';
import { isAuthenticated } from '../services/AuthServices';
import './AdminPanel.css'

const AdminPanel = () => {
    const [messages, setMessages] = createSignal<MessageResponse[]>([]);
    const [selectedMessageId, setSelectedMessageId] = createSignal<string | null>(null);
    const [loading, setLoading] = createSignal(true);
    const [error, setError] = createSignal('');

    const selectedMessage = () => {
        const id = selectedMessageId();
        if (!id) return null;
        return messages().find(m => m.id === id) || null;
    };

    onMount(async () => {
        console.log("AdminPanel montado, verificando autenticação...");
        
        const loadMessages = async () => {
            try {
                console.log("Carregando mensagens...");
                const data = await fetchMessages();
                setMessages(data);
                setError(''); // Limpar erro anterior
                console.log(`${data.length} mensagens carregadas com sucesso`);
            } catch (e) {
                console.error("Erro ao carregar mensagens:", e);
                const errorMessage = e instanceof Error ? e.message : 'Erro ao carregar mensagens';
                setError(errorMessage);
                
                // Se for erro de autenticação, limpar token e redirecionar
                if (errorMessage.includes('autenticado') || errorMessage.includes('Sessão expirada')) {
                    console.log("Erro de autenticação, limpando token e redirecionando...");
                    localStorage.removeItem("auth_token");
                    setTimeout(() => {
                        window.location.href = '/login';
                    }, 1000);
                }
            } finally {
                setLoading(false);
            }
        };

        // Verificar se há token antes de tentar carregar
        const token = localStorage.getItem("auth_token");
        if (!token || token.trim() === '') {
            console.log("Sem token válido, redirecionando para login");
            setError('Usuário não autenticado');
            setLoading(false);
            setTimeout(() => {
                window.location.href = '/login';
            }, 1000);
            return;
        }

        // Aguardar um pouco para garantir que o token foi salvo corretamente
        setTimeout(loadMessages, 200);
    });

    const handleSelectMessage = async (message: MessageResponse) => {
        setSelectedMessageId(message.id!);
        
        try {
            if (message.status === 'recebido') {
                console.log(`Selecionando mensagem ${message.id} e atualizando status`);
                const updated = await updateMessageStatus(message.id!, 'em análise');
                setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
                setError(''); // Limpar erro anterior
            }
        } catch (e) {
            console.error("Erro ao atualizar status:", e);
            const errorMessage = e instanceof Error ? e.message : String(e);
            setError(`Erro ao atualizar status: ${errorMessage}`);
        }
    };

    const handleAddObservation = async (text: string) => {
        const current = selectedMessage();
        if (!current) return;

        try {
            console.log(`Adicionando observação à mensagem ${current.id}`);
            const updated = await addMessageObs(current.id!, text);
            setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
            setError(''); // Limpar erro anterior
        } catch (e) {
            console.error("Erro ao adicionar observação:", e);
            const errorMessage = e instanceof Error ? e.message : String(e);
            setError(`Erro ao adicionar observação: ${errorMessage}`);
        }
    };

    const handleStatusChange = async (status: MessageStatus) => {
        const current = selectedMessage();
        if (!current) return;

        try {
            console.log(`Alterando status da mensagem ${current.id} para ${status}`);
            const updated = await updateMessageStatus(current.id!, status);
            setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
            setError(''); // Limpar erro anterior
        } catch (e) {
            console.error("Erro ao alterar status:", e);
            const errorMessage = e instanceof Error ? e.message : String(e);
            setError(`Erro ao alterar status: ${errorMessage}`);
        }
    };
    
    return (
        <div class="admin-container">
            <Show when={!loading()} fallback={<Spinner show={true} />}>
                <Show when={error()}>
                    <div class="alert error">
                        {error()}
                        <Show when={error().includes('autenticado')}>
                            <button 
                                onClick={() => window.location.href = '/login'}
                                style="margin-left: 10px; padding: 5px 10px; background: #007bff; color: white; border: none; border-radius: 3px; cursor: pointer;"
                            >
                                Ir para Login
                            </button>
                        </Show>
                    </div>
                </Show>
                
                <div class="messages-list">
                    <h2>Mensagens Recebidas ({messages().length})</h2>
                    <For each={messages()}>
                        {(message) => (
                            <div 
                                class={`message-card ${selectedMessageId() === message.id ? 'selected' : ''}`}
                                onClick={() => handleSelectMessage(message)}
                            >
                                <div class={`status-badge ${message.status.replace(' ', '-')}`}>
                                    {message.status}
                                </div>
                                <div class="message-preview">
                                    <p>{message.content.substring(0, 60)}...</p>
                                    <div class="media-indicator">
                                        <Show when={message.media && message.media.length > 0}>
                                            📎 {message.media!.length} anexos
                                        </Show>
                                    </div>
                                </div>
                            </div>
                        )}
                    </For>
                </div>

                <div class="message-detail">
                    <Show 
                        when={selectedMessage()} 
                        fallback={<div class="empty-state">Selecione uma mensagem</div>}
                    >
                        {() => {
                            const message = selectedMessage()!;
                            return (
                                <>  
                                    <div class="detail-header">
                                        <select
                                            value={message.status}
                                            onChange={(e) => handleStatusChange(e.target.value as MessageStatus)}
                                            class="status-select"
                                        >
                                            <For each={['recebido', 'em análise', 'concluído']}>
                                                {(status) => <option value={status}>{status}</option>}
                                            </For>
                                        </select>
                                        <span class="message-date">
                                            Recebida em: {new Date(message.createdAt!).toLocaleDateString()}
                                        </span>
                                    </div>

                                    <div class="content-section">
                                        <h3>Conteúdo:</h3>
                                        <p class="message-content">{message.content}</p>
                                    </div>

                                    {/* SUBSTITUA ESTA PARTE PELO MEDIAVIEWER */}
                                    <MediaViewer media={message.media || []} />

                                    <div class="observations-section">
                                        <h3>Observações:</h3>
                                        <form onSubmit={(e) => {
                                            e.preventDefault();
                                            const formData = new FormData(e.currentTarget);
                                            const observation = formData.get('observation') as string;
                                            if (observation.trim()) {
                                                handleAddObservation(observation.trim());
                                                (e.target as HTMLFormElement).reset();
                                            }
                                        }}>
                                            <textarea
                                                name="observation"
                                                placeholder="Adicione uma observação..."
                                                class="obs-textarea"
                                                required
                                            />
                                            <button type="submit" class="add-obs-button">
                                                Adicionar Observação
                                            </button>
                                        </form>

                                        <Show when={message.obs}>
                                            <div class="obs-history">
                                                <For each={message.obs?.split('\n').filter(Boolean)}>
                                                    {(obs, index) => (
                                                        <div class="obs-item">
                                                            <div class="obs-index">#{index() + 1}</div>
                                                            <p>{obs}</p>
                                                        </div>
                                                    )}
                                                </For>
                                            </div>
                                        </Show>
                                    </div>
                                </>
                            );
                        }}
                    </Show>
                </div>
            </Show>
        </div>
    );
}

export default AdminPanel;