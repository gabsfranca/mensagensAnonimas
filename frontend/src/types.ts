export type MessageStatus = 'recebido' | 'em análise' | 'concluído';
export type MediaType = 'image' | 'video' | 'audio'

export interface MediaAtachment {
    id: string; 
    url: string; 
    type: MediaType;
    thumbnail?: string;
}

export interface MessageResponse {
    id?: string;
    success: boolean; 
    error?: string; 
    content: string;
    status: MessageStatus;
    obs?: string;
    media?: MediaAtachment[], 
    createdAt?: string;
    updatedAt?: string;
}

// export interface MessageData {
//     id: string
//     content: string; 
//     status: MessageStatus;
//     obs?: string;
// }