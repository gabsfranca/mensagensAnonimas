import { JSX } from 'solid-js';
import { MessageResponse } from '../types';

interface ReportListProps {
    reports: MessageResponse[];
    selectedId?: string;
    onSelect: (report: MessageResponse) => void;
}

export const ReportList = (props: ReportListProps): JSX.Element => {
    return (
    <div class="report-list">
      <h2>Denúncias Recebidas ({props.reports.length})</h2>
      <div class="list-container">
        {props.reports.map(report => (
          <div 
            class={`report-item ${report.id === props.selectedId ? 'selected' : ''}`}
            onClick={() => props.onSelect(report)}
          >
            <div class="report-status-indicator" data-status={report.status} />
            <div class="report-preview">
              <h3>{report.content.substring(0, 50)}...</h3>
              <small>{new Date(report.createdAt!).toLocaleDateString()}</small>
              {report.media && report.media.length > 0 && (
                <div class="media-preview">
                  <span>📷 {report.media.length}</span>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};