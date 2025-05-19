import { For, Show, JSX } from 'solid-js';
import { MediaAtachment } from '../types';

interface MediaViewerProps {
    media: MediaAtachment[];
}

const getMediaComponent = (media: MediaAtachment): JSX.Element => {
    switch(media.type) {
        case 'image': 
            return (
                <img
                src={media.thumbnail || media.url}
                alt="Anexo de imagem"
                class="media-content"
                loading="lazy"
                />
            );
        case 'video': 
            return(
                <video controls class="media-content">
                <source src={media.url} type={`video/${media.url.split('.').pop()}`} />
                Seu navegador não suporta o elemento de vídeo.
                </video>
            );
        case 'audio':
            return(
                <audio controls class="media-content">
                <source src={media.url} type={`audio/${media.url.split('.').pop()}`} />
                Seu navegador não suporta o elemento de áudio.
                </audio>
            );
        default:
            return <div>Tipo de midia nao suportado</div>
    }
};

export const MediaViewer = (props: MediaViewerProps) => {
  return (
    <Show when={props.media.length > 0}>
      <div class="media-section">
        <h3>Anexos:</h3>
        <div class="media-grid">
          <For each={props.media}>
            {(media) => (
              <div class="media-item">
                <div class="media-container">
                  {getMediaComponent(media)}
                  <div class="media-type-badge">
                    {media.type[0].toUpperCase()}
                  </div>
                </div>
              </div>
            )}
          </For>
        </div>
      </div>
    </Show>
  );
};