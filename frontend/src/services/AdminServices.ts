import { MessageResponse, MessageStatus } from "../types";
import { getAuthHeaders, getTokenSafely } from "./AuthServices";
import { MediaType } from "../types";

interface ApiError {
    message: string; 
    statusCode: number;
}

export interface MediaAtachment {
    id: string;
    url: string;
    type: 'image' | 'video' | 'audio';
}

const URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

const mapReportToMessage = (report: any): MessageResponse => ({
    id: report.id,
    success: report.success,
    content: report.message,
    status: report.status,
    createdAt: report.createdAt,
    obs: report.obs,
    media: report.media?.map((m: any) => ({
        id: m.id,
        url: m.url,
        type: m.type.toLowerCase() as MediaType // Conversão para tipo compatível
    })) || []
});

const handleResponse = async <T>(response: Response): Promise<T> => {
    if (!response.ok) {
        if (response.status === 401) {
            localStorage.removeItem("auth_token");
            throw new Error("sessão expirada, faça login novamente");
        }

        try {
            const error: ApiError = await response.json();
            throw new Error(error.message || `erro HTTP!: ${response.status}`);
        } catch {
            throw new Error(`erro HTTP!: ${response.status}`);
        }
    }
    return response.json();
}

const makeAuthenticatedRequest = async (url: string, options: RequestInit = {}) => {
    const token = getTokenSafely();
    if (!token || token.trim() == '') {
        throw new Error("usuario nao autenticado");
    }

    console.log('fazendo req autenticada para: ', url);

    const headers = {
        "Content-Type": "application/json", 
        "Authorization": `Bearer ${token}`,
        ...(options.headers || {})
    };

    return fetch(url, {
        ...options, 
        headers: {
            ...headers, 
            ...(options.headers || {})
        }
    });
};

export const fetchMessages = async (): Promise<MessageResponse[]> => {
    try {
        console.log('buscando msgs...');
        const response = await makeAuthenticatedRequest(`${URL}/messages`, {
            method: 'GET',
        });

        const result = await handleResponse<MessageResponse[]>(response);
        console.log(`${result.length} msgs carregadas`);

        return result;
    } catch (e) {
        console.error(`erro no fetch das mensagens: ${e}`);
        throw e;
    }
};

export const updateMessageStatus = async (
    id: string, 
    status: MessageStatus
): Promise<MessageResponse> => {
    try {
        console.log(`atualizando status da msg ${id} para ${status}`)
        const response = await makeAuthenticatedRequest(`${URL}/messages/${id}/status`, {
            method: 'PATCH', 
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ status }),
        });

        const result = await handleResponse<MessageResponse>(response);

        console.log('status atualizado com sucesso!');

        return result;
        
    } catch (e) {
        console.error('falha ao atualizar status: ', e);    
        throw e;
    }
};

export const addMessageObs = async (
    id: string, 
    obs: string
): Promise<MessageResponse> => {
    try {
        console.log(`adicionando obs a msg ${id}`);
        const response = await makeAuthenticatedRequest(`${URL}/messages/${id}/obs`, {
            method: 'POST', 
            headers: {
                'Content-Type':'application/json',
            },
            body: JSON.stringify({ obs }),
        });
        const result = await handleResponse<MessageResponse>(response);
        console.log('obs add com successooooo');
        return result
    } catch (e) {
        console.error('falha ao add obs: ', e);
        throw e;
    }
};

export const getMessageObs = async(id: string): Promise<MessageResponse> => {
    try {
        console.log('buscando obs da msg: ', id);
        const response = await fetch(`${URL}/messages/${id}/obs`, {
            method: 'GET',
        });
        return handleResponse<MessageResponse>(response);
    } catch (e) {
        console.error('falha ao dar get nessa porraaa de observacao: ', e);
        throw e;
    }
}

export const getMessageDetails = async (
    id: string
): Promise<MessageResponse> => {
    try{
        console.log(`buscando detalhes da msg: ${id}`);
        const response = await makeAuthenticatedRequest(`${URL}/messages/${id}`, {
            method: 'GET',
        });
        return handleResponse<MessageResponse>(response);
    } catch (e) {
        console.error('falha no fetch dos detalhes da msg: ', e);
        throw e;
    }
};