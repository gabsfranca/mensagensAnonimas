import { MessageResponse, MessageStatus } from "../types";

interface ApiError {
    message: string; 
    statusCode: number;
}

const handleResponse = async <T>(response: Response): Promise<T> => {
    if (!response.ok) {
        const error: ApiError = await response.json();
        throw new Error(error.message || `HTTP error!: ${response.status}`);
    }
    return response.json();
}

export const fetchMessages = async (): Promise<MessageResponse[]> => {
    try {
        const response = await fetch('http://localhost:8080/messages', {
            credentials: 'include'
        });
        return handleResponse<MessageResponse[]>(response);
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
        const response = await fetch(`http://localhost:8080/messages/${id}/status`, {
            method: 'PATCH', 
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ status }),
            credentials: 'include'
        });
        return handleResponse<MessageResponse>(response);
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
        const response = await fetch(`http://localhost:8080/messages/${id}/obs`, {
            method: 'POST', 
            headers: {
                'Content-Type':'application/json',
            },
            body: JSON.stringify({ obs }),
            credentials: 'include'
        });
        return handleResponse<MessageResponse>(response);
    } catch (e) {
        console.error('falha ao add obs: ', e);
        throw e;
    }
};

export const getMessageObs = async(id: string): Promise<MessageResponse> => {
    try {
        const response = await fetch(`http://localhost:8080/messages/${id}/obs`, {
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
        const response = await fetch(`http://localhost:8080/messages/${id}`, {
            credentials: 'include'
        });
        return handleResponse<MessageResponse>(response);
    } catch (e) {
        console.error('falha no fetch dos detalhes da msg: ', e);
        throw e;
    }
};