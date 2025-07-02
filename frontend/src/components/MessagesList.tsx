import { For, createMemo, createSignal, Show, onMount } from 'solid-js';
import { MessageResponse, Tag } from '../types';
import '../styles/components/Messageslist.css';

interface MessageLIstProps {
  messages: MessageResponse[];
  selectedMessageId: string | null;
  onSelectMessage: (message: MessageResponse) => void;
  searchTerm: string;
  onSearchChange: (value: string) => void;
  onFilterChange: (selectedTags: string[]) => void;
  availableTags: Tag[];
  selectedTags: string[];
  loadTags: () => Promise<void>;
  isLoading: boolean;
}

export const MessagesList = (props: MessageLIstProps) => {
  const messageCount = createMemo(() => props.messages.length);

  const [showFilters, setShowFilters] = createSignal(false);
  const [localSelection, setLocalSelection] = createSignal<string[]>([]);

  onMount(() => {
    props.loadTags();
    setLocalSelection(props.selectedTags);
  });

  const toggleTag = (tagId: string) => {
    const current = localSelection();
    if (current.includes(tagId)) {
      setLocalSelection(current.filter(id => id !== tagId));
    } else {
      setLocalSelection([...current, tagId]);
    }
  };

  const applyFilters = () => {
    props.onFilterChange(localSelection());
  };

  const clearFilters = () => {
    setLocalSelection([]);
    props.onFilterChange([]);
  };

  return (
    <div class="messages-list">
      <input
        type="text"
        class="search-input"
        placeholder="Buscar por palavra-chave..."
        value={props.searchTerm}
        onInput={(e) => props.onSearchChange(e.currentTarget.value)}
      />

      <button onClick={() => setShowFilters(!showFilters())}>
        {showFilters() ? 'Ocultar filtros' : 'Exibir filtros'}
      </button>

      <Show when={showFilters()}>
        <div class="tag-filters">
          <h3>Filtrar por Tag</h3>
          <Show when={!props.isLoading}>
            <For each={props.availableTags}>
              {(tag) => (
                <label>
                  <input
                    type="checkbox"
                    value={tag.id}
                    checked={localSelection().includes(tag.id)}
                    onChange={() => toggleTag(tag.id)}
                  />
                  {tag.Name}
                </label>
              )}
            </For>

            <div style="margin-top: 8px;">
              <button onClick={applyFilters}>Aplicar filtros</button>
              <button onClick={clearFilters}>Limpar Filtros</button>  
            </div>
          </Show>
        </div>
      </Show>

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
