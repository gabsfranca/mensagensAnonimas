import { createSignal, createMemo } from 'solid-js';
import { MessageResponse } from '../types';

export const useTagFilters = (messages: () => MessageResponse[]) => {
    const [activeTagFilters, setActiveTagFilters] = createSignal<string[]>([]);
    const [searchTerm, setSearchTerm] = createSignal('');

    const filteredMessages = createMemo(() => {
        const selectedTags = activeTagFilters();
        const term = searchTerm().toLowerCase();
        const allMessages = messages();
        
        return allMessages.filter((msg) => {
            const matchesTags = 
                selectedTags.length === 0 ||
                msg.tags?.some((tag) => selectedTags.includes(tag.id));

            const matchesText = 
                msg.content.toLowerCase().includes(term) ||
                (msg.obs?.toLowerCase().includes(term) ?? false);
        
            return matchesTags && matchesText;
        });
    });

    return {
        activeTagFilters,
        setActiveTagFilters, 
        searchTerm, 
        setSearchTerm,
        filteredMessages,
    };
};