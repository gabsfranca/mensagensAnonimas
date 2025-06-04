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

// 1. Defina a URL base da API como constante
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const getMediaComponent = (media: MediaAtachment): JSX.Element => {
    console.log('🎬 Renderizando mídia:', media);

    // 2. Construa as URLs completas corretamente
    const fullUrl = `${API_BASE_URL}${media.url}`;
    const fullThumbnail = media.thumbnail ? `${API_BASE_URL}${media.thumbnail}` : null;
    
    console.log(`fullUrl: ${fullUrl}`);
    console.log(`fullThumbnail: ${fullThumbnail}`);
    
    switch(media.type) {
        case 'image': 
            return (
                <img
                    src={fullThumbnail || fullUrl}  // Use a URL completa
                    alt="Anexo de imagem"
                    class="media-content"
                    loading="lazy"
                    onLoad={() => console.log('✅ Imagem carregada:', fullUrl)}
                    onError={(e) => console.error('❌ Erro ao carregar imagem:', fullUrl, e)}
                />
            );
        case 'video': 
            return(
                <video 
                    controls 
                    class="media-content"
                    onLoadedData={() => console.log('✅ Vídeo carregado:', fullUrl)}
                    onError={(e) => console.error('❌ Erro ao carregar vídeo:', fullUrl, e)}
                >
                    {/* 3. Corrija para usar a URL completa no source */}
                    <source src={fullUrl} type={`video/${media.url.split('.').pop()}`} />
                    Seu navegador não suporta o elemento de vídeo.
                </video>
            );
        case 'audio':
            return(
                <audio 
                    controls 
                    class="media-content"
                    onLoadedData={() => console.log('✅ Áudio carregado:', fullUrl)}
                    onError={(e) => console.error('❌ Erro ao carregar áudio:', fullUrl, e)}
                >
                    {/* 4. Corrija para usar a URL completa no source */}
                    <source src={fullUrl} type={`audio/${media.url.split('.').pop()}`} />
                    Seu navegador não suporta o elemento de áudio.
                </audio>
            );
        default:
            console.warn('⚠️ Tipo de mídia não suportado:', media.type);
            return <div>Tipo de mídia não suportado: {media.type}</div>
    }
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