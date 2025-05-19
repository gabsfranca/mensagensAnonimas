import { MessageResponse } from '../types';

export const sendAnonymousMessage = async (formData: FormData): Promise<MessageResponse> => {
    try {
        const response = await fetch('http://localhost:8080/send-anonymous-message', {
            method: 'POST', 
            body: formData
        });

        const result = await response.json();

        if (!response.ok) {
            return {
                success: false,
                error: result.error || 'Erro desconhecido',
                content: result.content,
                status: result.status,
            };
        }

        return {
            id: result.id,
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