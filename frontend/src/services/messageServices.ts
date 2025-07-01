import { MessageResponse } from '../types';

const URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

export const sendAnonymousMessage = async (formData: FormData): Promise<MessageResponse> => {
    try {
        const response = await fetch(`${URL}/send-anonymous-message`, {
            method: 'POST', 
            body: formData
        });

        const result = await response.json();

        if (!response.ok) {
            return {
                success: false,
                shortId: result.shortId,
                error: result.error || 'Erro desconhecido',
                content: result.content,
                status: result.status,
            };
        }

        return {
            id: result.id,
            shortId: result.shortId,
            success: true,
            content: result.content,
            status: result.status,
        };
    } catch (e) {
        console.error('Erro: ', e);
        return {
            success: false,
            error: 'Erro ao conectar com o sv',
            content: '',
            status: 'recebido',
        };
    }   
}