import { createSignal, createMemo } from 'solid-js';
import { MessageResponse, MessageStatus } from '../types';
import {
    fetchMessages,
    updateMessageStatus,
    addMessageObs
} from '../services/AdminServices';
import { getObservations } from '../services/ObservationServices';
import type { Observation } from '../types'; // ajuste o caminho se necessário
import crypto from 'crypto';

export const useMessages = () => {
    const [messages, setMessages] = createSignal<MessageResponse[]>([]);
    const [selectedMessageId, setSelectedMessageId] = createSignal<string | null>(null);
    const [isLoading, setIsLoading] = createSignal(false);
    const [error, setError] = createSignal('');
    const [observations, setObservations] = createSignal<Observation[]>([]); // ✅

    const selectedMessage = createMemo(() => {
        const id = selectedMessageId();
        if (!id) return null;
        return messages().find(m => m.id === id) || null;
    });

    const loadMessages = async () => {
        setIsLoading(true);
        setError('');
        try {
            const data = await fetchMessages();
            setMessages(data);
        } catch (e) {
            const errorMessage = e instanceof Error ? e.message : 'Erro ao carregar mensagens';
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

    const selectMessage = async (message: MessageResponse) => {
        setSelectedMessageId(message.id!);

        // ✅ Carrega as observações quando uma mensagem for selecionada
        try {
            const result = await getObservations(message.id!);
            if (result.success) {
                setObservations(result.observations);
            } else {
                console.error('Erro ao carregar observações:', result.error);
                setObservations([]);
            }

            if (message.status === 'recebido') {
                await updateMessageStatusLocal(message.id!, 'em análise');
            }
        } catch (e) {
            console.error('Erro ao selecionar mensagem:', e);
            setError(e instanceof Error ? e.message : String(e));
        }
    };

    const updateMessageStatusLocal = async (messageId: string, status: MessageStatus) => {
        const updated = await updateMessageStatus(messageId, status);
        setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
        setError('');
    };

    const addObservation = async (text: string) => {
        const current = selectedMessage();
        if (!current) return;

        try {
            const updated = await addMessageObs(current.id!, text);
            setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));

            
            const result = await getObservations(current.id!);
            if (result.success) {
            setObservations(result.observations);
            } else {
            console.error('Erro ao recarregar observações:', result.error);
            }
            setError('');
        } catch (e) {
            const errorMessage = e instanceof Error ? e.message : String(e);
            setError(`Erro ao adicionar observação: ${errorMessage}`);
        }
    };

    const updateMessageLocal = (updated: MessageResponse) => {
        setMessages(prev => prev.map(m => m.id === updated.id ? updated : m));
    };

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
        updateMessageLocal,
        setError,
        observations,       // ✅ exporta observações para o componente
        setObservations     // opcional, caso precise resetar
    };
};
