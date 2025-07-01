//eu acho que o problema ta no hook 

import { Show, For, createMemo } from 'solid-js';
import { MessageResponse, Tag } from '../types';
import '../styles/components/TagsSection.css'

interface TagsSectionProps {
  message: MessageResponse | null;
  availableTags: Tag[];
  selectedTags: string[];
  onTagsChange: (tags: string[]) => void;
  onAddTags: () => void;
  onRemoveTag: (tagId: string) => void;
  isLoading?: boolean;
}

export const TagsSection = (props: TagsSectionProps) => {
  const handleTagSelection = (e: Event) => {
    const select = e.currentTarget as HTMLSelectElement;
    const options = select.selectedOptions;
    const selected = Array.from(options).map(opt => opt.value);
    props.onTagsChange(selected);
  };

  return (
    <Show when={props.message}>
      {message => (
        <div class="tags-section">
          <h3>Tags:</h3>
          
          <div class="current-tags">
            <Show 
              when={(props.message?.tags ?? []).length > 0} 
              fallback={<span class="no-tags">Nenhuma tag adicionada</span>}
            >
              <For each={props.message?.tags ?? []}>
                {(tag) => (
                  <span class="tag-badge">
                    <span>{tag.Name || "Tag sem nome"}</span>
                    <button
                      class='tag-remove'
                      onclick={(e) => {
                        e.stopPropagation();
                        props.onRemoveTag(tag.id);
                      }}
                    >x</button>
                  </span>
                )}
              </For>
            </Show>
          </div>
          
          <div class="add-tags">
            <Show 
              when={props.availableTags.length > 0}
              fallback={<div>Carregando tags...</div>}
            >
              <select 
                multiple
                size={5}
                value={props.selectedTags}
                onChange={handleTagSelection}
              >
                <For each={props.availableTags}>
                  {tag => (
                    <option value={tag.id}>
                      {tag.Name || "sem nome"}
                    </option>
                  )}
                </For>
              </select>
            </Show>
            <button 
              onClick={props.onAddTags}
              disabled={props.isLoading}
            >
              {props.isLoading ? 'Adicionando...' : 'Adicionar Tags'}
            </button>
          </div>
        </div>
      )}
    </Show>
  );
};
