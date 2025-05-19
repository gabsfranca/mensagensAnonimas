import { createSignal, JSX } from 'solid-js';
import { MessageResponse, MediaAtachment } from '../types';
import { MediaViewer } from './MediaViewer';

interface ReportDetailProps {
    report: MessageResponse | null;
    onAddNote: (note: string) => void; 
    onStatusChange: (status: MessageResponse['status']) => void;
}

export const ReportDetail = (props: ReportDetailProps): JSX.Element => {
    const [newNote, setNewNote] = createSignal('');

    const handleSubmitNote = (e: Event) => {
        e.preventDefault();
        if (newNote().trim()) {
            props.onAddNote(newNote());
            setNewNote('');
        }
    }

    return (
    <div class="report-detail">
      {props.report ? (
        <>
          <div class="status-header">
            <select 
              value={props.report.status}
              onChange={(e) => props.onStatusChange(e.target.value as MessageResponse['status'])}
            >
              <option value="recebido">Recebido</option>
              <option value="em análise">Em Análise</option>
              <option value="concluído">Concluído</option>
            </select>
          </div>

          <div class="message-content">
            <h3>Mensagem:</h3>
            <p>{props.report.content}</p>
          </div>

          <MediaViewer media={props.report.media!} />

          <div class="notes-section">
            <h3>Observações:</h3>
            <form onSubmit={handleSubmitNote}>
              <textarea
                value={newNote()}
                onInput={(e) => setNewNote(e.currentTarget.value)}
                placeholder="Adicionar nova observação..."
              />
              <button type="submit">Adicionar Observação</button>
            </form>

            <div class="notes-list">
                <div class="note-item">
                  <p>{props.report.obs}</p>
                </div>
            </div>
          </div>
        </>
      ) : (
        <div class="empty-state">Selecione uma denúncia para visualizar</div>
      )}
    </div>
  );
};
