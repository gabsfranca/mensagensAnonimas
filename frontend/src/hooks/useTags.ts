import { createSignal } from "solid-js";
import { Tag } from "../types";
import { fetchAvailableTags, addTagsToReport } from "../services/AdminServices";

export const useTags = () => {
    const [availableTags, setAvailableTags] = createSignal<Tag[]>([]);
    const [selectedTags, setSelectedTags] = createSignal<string[]>([]);
    const [isLoading, setIsLoading] = createSignal(false);
    const [error, setError] = createSignal('');

    const loadTags = async () => {
        setIsLoading(true);
        try {
            const tags = await fetchAvailableTags();
            console.log('Tags carregadas: ', tags);
            setAvailableTags(tags);
        } catch (e) {
            console.error('Erro ao carregar tags: ', e);
            setError('Erro ao carregar Tags disponiveis');
        } finally {
            setIsLoading(false);
        }
    };

    const addTagsToMessage = async (messageId: string, onSuccess: (updated: any) => void) => {
        if (!messageId) {
            console.log('nenhuma mensagem selecionada');
            return
        }

        try {
            console.log('Tentando adicionar tags');
            const updated = await addTagsToReport(messageId, selectedTags());
            console.log('Mensagem atualizada: ', updated);
            onSuccess(updated);
            setSelectedTags([]);
            setError('');
        } catch (e) {
            console.error('Erro ao adicionar tags: ', e);
            const errorMessage = e instanceof Error ? e.message : String(e);
            setError(`Erro ao adicionar tags à mensagem: ${errorMessage}`);
        }
    };

    return {
        availableTags, 
        selectedTags, 
        isLoading, 
        error, 
        loadTags, 
        addTagsToMessage, 
        setSelectedTags
    };
};