import { For, createMemo } from 'solid-js';
import { MessageResponse } from '../types';
import '../styles/components/Messageslist.css'

interface MessageLIstProps {
    messages: MessageResponse[];
    selectedMessageId: string | null;
    onSelectMessage: (message: MessageResponse) => void;
}

export const MessagesList = (props: MessageLIstProps) => {
    const messageCount = createMemo(() => props.messages.length);

    return (
    <div class="messages-list">
      <h2>Mensagens Recebidas ({messageCount()})</h2>
      <For each={props.messages}>
        {(message) => (
          <div 
            class={`message-card ${props.selectedMessageId === message.id ? 'selected' : ''}`}
            onClick={() => props.onSelectMessage(message)}
          >
            <div class={`status-badge ${message.status.replace(' ', '-')}`}>
              {message.status}
            </div>
            <div class="message-preview">
              <p>{message.content.substring(0, 60)}...</p>
              <div class="media-indicator">
                {message.media && message.media.length > 0 && (
                  <>📎 {message.media.length} anexos</>
                )}
              </div>
            </div>
          </div>
        )}
      </For>
    </div>
  );
};