import { Show, For } from 'solid-js';
import { MessageResponse, MessageStatus } from '../types';
import { MediaViewer } from './MediaViewer';
import '../styles/components/MessageDetail.css'

interface MessageDetailProps {
  message: MessageResponse | null;
  onStatusChange: (status: MessageStatus) => void;
  onAddObservation: (text: string) => void;
}

export const MessageDetail = (props: MessageDetailProps) => {
  const handleObservationSubmit = (e: SubmitEvent) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget as HTMLFormElement);
    const observation = formData.get('observation') as string;
    if (observation?.trim()) {
      props.onAddObservation(observation.trim());
      (e.target as HTMLFormElement).reset();
    }
  };

  return (
    <div class="message-detail">
      <Show 
        when={props.message} 
        fallback={<div class="empty-state">Selecione uma mensagem</div>}
      >
        {message => (
          <>  
            <div class="detail-header">
              <select
                value={message().status}
                onChange={(e) => props.onStatusChange(e.target.value as MessageStatus)}
                class="status-select"
              >
                <For each={['recebido', 'em análise', 'concluído']}>
                  {(status) => <option value={status}>{status}</option>}
                </For>
              </select>
              <span class="message-date">
                Recebida em: {new Date(message().createdAt!).toLocaleDateString()}
              </span>
            </div>

            <div class="content-section">
              <h3>Conteúdo:</h3>
              <p class="message-content">{message().content}</p>
            </div>

            <MediaViewer media={message().media || []} />

            <div class="observations-section">
              <h3>Observações:</h3>
              <form onSubmit={handleObservationSubmit}>
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

              <Show when={message().obs}>
                <div class="obs-history">
                  <For each={message().obs?.split('\n').filter(Boolean)}>
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
        )}
      </Show>
    </div>
  );
};