import { For, Show, JSX, createEffect } from 'solid-js';

export type MediaType = 'image' | 'video' | 'audio'

export interface MediaAtachment {
    id: string; 
    url: string; 
    type: MediaType;
    thumbnail?: string;
}

interface MediaViewerProps {
    media: MediaAtachment[];
}

//const API_BASE_URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

//const API_BASE_URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

const getMediaComponent = (media: MediaAtachment): JSX.Element => {

  
    const API_BASE_URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080"

    const fullUrl = `${API_BASE_URL}${media.url}`;
    return (
        <a
            href={fullUrl}
            target='_blank'
            rel='noopener noreferrer'
            class='media.link'
        >
            Abrir {media.type} ({media.url.split('/').pop()})
        </a>
    )
};

export const MediaViewer = (props: MediaViewerProps) => {
    console.log('🚀 MediaViewer renderizado!');
    console.log('📊 Props recebidos:', props);
    console.log('📁 Media array:', props.media);
    console.log('🔢 Quantidade de itens:', props.media?.length || 0);

    // Verificar se props.media existe e tem conteúdo
    createEffect(() => {
        console.log('🔄 Effect executado - props.media:', props.media);
        if (!props.media) {
            console.warn('⚠️ props.media é undefined ou null');
        } else if (props.media.length === 0) {
            console.warn('⚠️ props.media está vazio');
        } else {
            console.log('✅ props.media tem', props.media.length, 'itens');
            props.media.forEach((item, index) => {
                console.log(`📄 Item ${index}:`, item);
            });
        }
    });

    // Sempre renderizar algo para teste
    return (
    <div style="padding: 10px; margin: 10px;">            
        <div class="media-section">
            <h3>✅ Anexos: ({props.media.length})</h3>
            <div class="media-grid">
                <For each={props.media}>
                    {(media, index) => {
                        console.log(`🎯 Renderizando For item ${index()}:`, media);
                        return (
                            <div class="media-item" style="border: 1px solid blue; padding: 5px; margin: 5px;">
                                <div class="media-container">
                                    {getMediaComponent(media)}
                                    <div class="media-type-badge">
                                        {media.type?.[0]?.toUpperCase() || '?'}
                                    </div>
                                </div>
                                <div style="font-size: 10px; color: #666; margin-top: 4px;">
                                    ID: {media.id}<br/>
                                    URL: {media.url}<br/>
                                    Tipo: {media.type}
                                </div>
                            </div>
                        );
                    }}
                </For>
            </div>
        </div>
    </div>
);
};
