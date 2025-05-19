import { For, createSignal, onMount, Show, Accessor } from 'solid-js';
import { MessageResponse, MessageStatus } from '../types';
import { Spinner } from './spinner';
import { fetchMessages, updateMessageStatus, addMessageObs } from '../services/AdminServices';

export const AdminPanel = () => {
    const [messages, setMessages] = createSignal<MessageResponse[]>([]);
    const [selectedMessage, setSelectedMessage] = createSignal<MessageResponse | null>(null);
    const [loading, setLoading] = createSignal(true);
    const [error, setError] = createSignal('');

    onMount(async () => {
        try {
            const data = await fetchMessages();
            setMessages(data);
        } catch (e) {
            setError('Erro ao carregar msgs');
        } finally {
            setLoading(false);
        }
    });

    const handleSelectMessage = async (message: MessageResponse) => {
        try {
            if (message.status === 'recebido') {
                const updated = await updateMessageStatus(message.id!, 'em análise');
                setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
                setSelectedMessage(updated);
            } else {
                setSelectedMessage(message);
            } 
        } catch (e) {
            setError(`erro ao atualizar status: ${e instanceof Error ? e.message : String(e)}`);
        }
    };

    const handleAddObservation = async (text: string) => {
        const current = selectedMessage();
        if (!current) return;

        try {
            const updated = await addMessageObs(current.id!, text);
            setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
            setSelectedMessage(updated);
        } catch (e) {
            setError(`erro ao add obs: ${e}`);
        }
    };

    const handleStatusChange = async (status: MessageStatus) => {
        const current = selectedMessage();
        if (!current) return;

        try {
            const updated = await updateMessageStatus(current.id!, status);
            setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
            setSelectedMessage(updated);
        } catch (e) {
            setError(`erro ao atualizar status: ${e}`);
        }
    };
    
    return (
  <div class="admin-container">
    <Show when={!loading()} fallback={<Spinner show={true} />}>
      <Show when={error()}>
        <div class="alert error">{error()}</div>
      </Show>
      
      <div class="messages-list">
        <h2>Mensagens Recebidas ({messages().length})</h2>
        <For each={messages()}>
          {(message) => (
            <div 
              class={`message-card ${selectedMessage()?.id === message.id ? 'selected' : ''}`}
              onClick={() => handleSelectMessage(message)}
            >
              <div class={`status-badge ${message.status.replace(' ', '-')}`}>
                {message.status}
              </div>
              <div class="message-preview">
                <p>{message.content.substring(0, 60)}...</p>
                <div class="media-indicator">
                  <Show when={message.media!.length > 0}>
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
          {(msg: Accessor<MessageResponse>) => {
            const message = msg();
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

                <Show when={message.media!.length > 0}>
                  <div class="media-section">
                    <h3>Anexos:</h3>
                    <div class="media-grid">
                      <For each={message.media}>
                        {(media) => (
                          <div class="media-item">
                            <Show
                              when={media.type === 'image'}
                              fallback={
                                media.type === 'video' ? (
                                  <video controls src={media.url} />
                                ) : (
                                  <audio controls src={media.url} />
                                )
                              }
                            >
                              <img 
                                src={media.thumbnail || media.url} 
                                alt="Anexo" 
                                class="media-content"
                              />
                            </Show>
                          </div>
                        )}
                      </For>
                    </div>
                  </div>
                </Show>

                <div class="observations-section">
                  <h3>Observações:</h3>
                  <form onSubmit={(e) => {
                    e.preventDefault();
                    const formData = new FormData(e.currentTarget);
                    handleAddObservation(formData.get('observation') as string);
                  }}>
                    <textarea
                      name="observation"
                      placeholder="Adicione uma observação..."
                      class="obs-textarea"
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