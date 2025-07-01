import { onMount, Show } from 'solid-js';
import { Spinner } from './spinner';
import { useAuth } from '../hooks/useAuth';
import { useMessages } from '../hooks/useMessages';
import { useTags } from '../hooks/useTags';
import { MessagesList } from './MessagesList';
import { MessageDetail } from './MessageDetails';
import { TagsSection } from './TagsSection';
import { ErrorBoundary } from './ErrorBoundary';
import { TagFilters } from './TagFilters';
import { useTagFilters } from '../hooks/useTagFilters'

import '../styles/AdminPanel.css';

const URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

const AdminPanel = () => {
  const auth = useAuth();
  const messages = useMessages();
  const observations = messages.observations;
  const tags = useTags();

  const {
    activeTagFilters,
    setActiveTagFilters,
    searchTerm,
    setSearchTerm,
    filteredMessages,
  } = useTagFilters(messages.messages);



  onMount(async () => {
    console.log("AdminPanel montado, verificando autenticação...");
    
    if (!auth.isAuthenticated()) return;

    // Load data in parallel for better performance
    await Promise.all([
      messages.loadMessages(),
      tags.loadTags()
    ]);
  });

  const handleAddTags = () => {
    const messageId = messages.selectedMessageId();
    if (!messageId) return;

    tags.addTagsToMessage(messageId, (updated) => {
      messages.updateMessageLocal(updated);
      const messagesData = messages.messages();
      const updatedMessages = messagesData.map(m => 
        m.id === updated.id 
          ? { 
              ...m, 
              tags: updated.tags?.map((t: any) => ({ id: t.id, Name: t.Name })) 
            } 
          : m
      );
    });
  };

  const handleRemoveTag = async (tagId: string) => {
    const messageId = messages.selectedMessageId();
    if (!messageId || !messages.selectedMessage()) return;

    try {

      const token = localStorage.getItem("auth_token");
      if (!token) {
        throw new Error("Token de autenticação não encontrado!");
      }


      const response = await fetch(`${URL}/messages/${messageId}/tags/${tagId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) throw new Error('falha ao remover tag');

      const currentMessage = messages.selectedMessage();
      if (currentMessage) {
        const updatedTags = currentMessage.tags?.filter(t => t.id !== tagId) || [];

        const updatedMessage = {
          ...currentMessage, 
          tags: updatedTags
        };

        messages.updateMessageLocal(updatedMessage);
      }
    }catch (e) {
      console.error('erro ao remover tag: ', e);
    }
  }

  // Show loading state during auth check
  if (auth.isLoading()) {
    return <Spinner show={true} />;
  }

  // Show auth error
  if (auth.error()) {
    return <ErrorBoundary error={auth.error()} showLoginButton={true} />;
  }

  return (
    <div class="admin-container">
      <Show when={!messages.isLoading()} fallback={<Spinner show={true} />}>
        <ErrorBoundary 
          error={messages.error() || tags.error()} 
          showLoginButton={true}
          onRetry={() => {
            messages.loadMessages();
            tags.loadTags();
          }}
        />
        
        <TagFilters
          onFilterChange={setActiveTagFilters}
        />
        
        <input
            type="text"
            class="search-input"
            placeholder="Buscar por palavra-chave..."
            value={searchTerm()}
            onInput={(e) => setSearchTerm(e.currentTarget.value)}
          />
        <MessagesList 
          messages={filteredMessages()}
          selectedMessageId={messages.selectedMessageId()}
          onSelectMessage={messages.selectMessage}
        />

        <MessageDetail 
          message={messages.selectedMessage()}
          onStatusChange={(status) => {
            const selected = messages.selectedMessage();
            if (selected) {
              messages.updateMessageStatusLocal(String(selected.id), status)
            }
          }}
          onAddObservation={messages.addObservation}
          observations={observations()}
        />

        <TagsSection 
          message={messages.selectedMessage()}
          availableTags={tags.availableTags()}
          selectedTags={tags.selectedTags()}
          onTagsChange={tags.setSelectedTags}
          onAddTags={handleAddTags}
          onRemoveTag={handleRemoveTag}
          isLoading={tags.isLoading()}
        />
      </Show>
    </div>
  );
};

export default AdminPanel;