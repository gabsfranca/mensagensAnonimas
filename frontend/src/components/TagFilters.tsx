import { createSignal, For, Show } from 'solid-js';
import { useTags } from '../hooks/useTags';
import { onMount } from 'solid-js';

interface TagFiltersProps {
    onFilterChange: (selectedTags: string[]) => void
}

export const TagFilters = (props: TagFiltersProps) => {
    const {
        availableTags, 
        selectedTags, 
        setSelectedTags, 
        loadTags, 
        isLoading, 
    } = useTags();

    const [localSelection, setLocalSelection] = createSignal<string[]>([]);

    onMount(() => {
      loadTags(); // carrega a tag da primeira vez
      setLocalSelection(selectedTags());
    });

    

    const toggleTags = (tagId: string) => {
        const current = localSelection();
        if (current.includes(tagId)) {
            setLocalSelection(current.filter((id) => id !== tagId));
        } else {
            setLocalSelection([...current, tagId]);
        }

        // props.onFilterChange(selectedTags().includes(tagId) COMENTADO POIS VOU COLOCAR O BOTAO DE FILTRAR
        //     ? selectedTags().filter((id) => id !== tagId)
        //     : [...selectedTags(), tagId]);
    };

    const clearFilters = () => {
        setLocalSelection([]);
        setSelectedTags([]);
        props.onFilterChange([]);
    };

    const applyFilters = () => {
      setSelectedTags(localSelection());
      props.onFilterChange(localSelection());
    };

     return (
    <div class="tag-filters">
      <h3>Filtrar por Tag</h3>
      <Show when={!isLoading()}>
        <For each={availableTags()}>
          {(tag) => (
            <label>
              <input
                type="checkbox"
                value={tag.id}
                checked={localSelection().includes(tag.id)}
                onChange={() => toggleTags(tag.id)}
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
  );
};



