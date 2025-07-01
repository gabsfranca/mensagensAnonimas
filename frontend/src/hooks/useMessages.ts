import { createSignal, createMemo } from 'solid-js';
import { MessageResponse, MessageStatus } from '../types';
import {
    fetchMessages,
    updateMessageStatus,
    addMessageObs
} from '../services/AdminServices';

export const useMessages = () => {
    const [messages, setMessages] = createSignal<MessageResponse[]>([]);
    const [selectedMessageId, setSelectedMessageId] = createSignal<string | null>(null);
    const [isLoading, setIsLoading] = createSignal(false);
    const [error, setError] = createSignal('');

    const selectedMessage = createMemo(() => {
        const id = selectedMessageId();
        if (!id) return null;
        return messages().find(m => m.id === id) || null;
    });

    const loadMessages = async () => {
        setIsLoading(true);
        setError('');

        try {
            // console.log('carregando mensagens...')
            const data = await fetchMessages();
            setMessages(data);
            // console.log(`${data.length} mensagens carregadas com sucesso`);
        } catch (e) {
            // console.error('Erro ao carregar mensagens: ', e);
            const errorMessage = e instanceof Error ? e.message : 'Erro ao carregar mensagens'
            setError(errorMessage);

            if (errorMessage.includes('autenticado') || errorMessage.includes('Sessão expirada')) {
                localStorage.removeItem("auth_token");
                setTimeout(() => {
                    window.location.href = '/login';
                }, 1000);
            }
        } finally {
            setIsLoading(false);
        }
    };

    const selectMessage =async (message: MessageResponse) => {
        setSelectedMessageId(message.id!);

        try {
            if(message.status === 'recebido') {
                // console.log(`Selecionando mensagem ${message.id} e atualizando o status`);
                await updateMessageStatusLocal(message.id!, 'em análise');
            }
        } catch (e) {
            console.error('erro ao atualizar status: ', e);
            const errorMessage = e instanceof Error ? e.message : String(e);
            setError(errorMessage);
        }
    };

    const updateMessageStatusLocal = async (messageId: string, status: MessageStatus) => {
        const updated = await updateMessageStatus(messageId, status);
        setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
        setError('');
    }

    const addObservation = async (text: string) => {
        const current = selectedMessage();
        if (!current) return;

        try {
            console.log(`Adicionando observação à mensagem ${current.id}`);
            const updated = await addMessageObs(current.id!, text);
            setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
            setError('');
        } catch (e) {
            console.error('Erro ao adicionar observação: ', e);
            const errorMessage = e instanceof Error ? e.message : String(e);
            setError(`Erro ao adicionar observação: ${errorMessage}`);
        }
    };

    const updateMessageLocal = (updated: MessageResponse) => {
        setMessages(prev => prev.map(m => m.id === updated.id ? updated: m));
    }

    return {
        messages, 
        selectedMessage, 
        selectedMessageId,
        isLoading, 
        error, 
        loadMessages, 
        selectMessage, 
        updateMessageStatusLocal, 
        addObservation, 
        setError,
        updateMessageLocal
    };
};
