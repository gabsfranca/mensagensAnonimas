import { MessageResponse } from "../types";

const URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

export const checkReportStatus = async (uuid: string): Promise<MessageResponse> => {
    try {

        console.log('verificando status da denuncia: ', uuid);

        const response = await fetch(`${URL}/reports/${uuid}/status`, {
            method: 'GET', 
            headers: {
                'Content-Type': 'application/json',
            },
        });

        console.log('resposta recebida com sucesso: ', response.status);

        const contentType = response.headers.get("content-type");
        let result;

        if (contentType && contentType.includes("application/json")) {
            try {
                result = await response.json();
                console.log('dados do json recebidos: ', result);
            } catch (e) {
                console.error('erro ao parsear json: ', e);
                const text = await response.text();
                throw new Error(`Resposta invalida do servidor: ${text}`);
            }
        } else {
            const text = await response.text();
            console.log('resposta text recebida: ', text);
            throw new Error(`resposta naojson do sv: ${text}`);
        }

        if (!response.ok) {
            return {
                success: false,
                error: result?.error || `Erro no sv: ${response.status}`,
                content: result?.content,
                status: result?.status
            };
        }

        return {
            id: result.id,
            success: true,
            content: result.content, 
            status: result.status,
        };
    } catch (e) {
        console.error('erro: ', e);
        return {
            success: false, 
            error: e instanceof Error ? e.message : 'Erro ao conectar com o servidor',
            content: '', 
            status: 'recebido', 
        }
    }
}