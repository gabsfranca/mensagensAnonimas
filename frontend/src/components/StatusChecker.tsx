import { createSignal, JSX, For } from "solid-js";
import { checkReportStatus } from "../services/statusServices";
import { getObservations, postObservation } from "../services/ObservationServices";
import { Spinner } from "./spinner";

export const StatusChecker = (): JSX.Element => {
  const [uuid, setUuid] = createSignal('');
  const [status, setStatus] = createSignal('');
  const [isLoading, setIsLoading] = createSignal(false);
  const [alertMessage, setAlertMessage] = createSignal('');
  const [observations, setObservations] = createSignal<{ content: string, author: string, createdAt: string }[]>([]);
  const [newObs, setNewObs] = createSignal('');

  const statusSteps = [
    { id: 1, label: 'Recebido', key: 'recebido' },
    { id: 2, label: 'Em análise', key: 'em análise' },
    { id: 3, label: 'Concluído', key: 'concluído' }
  ];

  const handleUuidChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    setUuid(target.value.trim());
  };

  const isValidUuid = (id: string) => /^[a-zA-Z0-9]{5,10}$/.test(id);

  const handleCheckStatus = async () => {
    setAlertMessage('');
    setStatus('');
    setObservations([]);

    if (!uuid()) {
      setAlertMessage('Por favor, digite o ID da denúncia');
      return;
    }

    if (!isValidUuid(uuid())) {
      setAlertMessage('ID inválido');
      return;
    }

    setIsLoading(true);
    try {
      const result = await checkReportStatus(uuid());
      if (result.success) {
        setStatus(result.status);
        setAlertMessage('');
        // Buscar observações
        const obsRes = await getObservations(uuid());
        if (obsRes.success) {
          setObservations(obsRes.observations);
        } else {
          setAlertMessage('Erro ao buscar observações');
        }
      } else {
        setAlertMessage(result.error || 'Erro ao buscar status');
      }
    } catch (e) {
      console.error('erro:', e);
      setAlertMessage('Erro ao conectar com o servidor');
    } finally {
      setIsLoading(false);
    }
  };

  const handleNewObservationSubmit = async () => {
    if (!newObs().trim()) return;

    try {
      const response = await postObservation(uuid(), newObs().trim());
      if (response.success) {
        // Adiciona à lista atual sem resetar
        const obsRes = await getObservations(uuid())
        if (obsRes.success) {
          setObservations(obsRes.observations);
        }
      } else {
        setAlertMessage('Erro ao enviar observação');
      }
    } catch (e) {
      console.error(e);
      setAlertMessage('Erro ao enviar observação');
    }
  };

  const getStatusColor = () => {
    switch (status().toLowerCase()) {
      case 'recebido':
        return 'status-received';
      case 'em análise':
        return 'status-analyzing';
      case 'concluído':
        return 'status-completed';
      default:
        return '';
    }
  };

  return (
    <div class="status-container">
      <h2>Deseja consultar o andamento?</h2>
      <div class="status-form">
        <input
          type="text"
          value={uuid()}
          onInput={handleUuidChange}
          placeholder="Digite o ID da denúncia"
          class="uuid-input"
        />
        <button
          onClick={handleCheckStatus}
          disabled={isLoading()}
          class="check-button"
        >
          <span class="button-content">
            <span>{isLoading() ? 'Consultando...' : 'Consultar Status'}</span>
            <Spinner show={isLoading()} />
          </span>
        </button>
      </div>

      {status() && (
        <div class="status-tracker">
          <div class="status-progress">
            {statusSteps.map((step) => (
              <div class={`status-step ${
                statusSteps.findIndex(s => s.key === status().toLowerCase()) >= step.id - 1 ? 'active' : ''
              }`}>
                <div class="step-icon">
                  {status().toLowerCase() === step.key ? (
                    <div class="current-step" />
                  ) : (
                    <div class="completed-icon">✓</div>
                  )}
                </div>
                <div class="step-label">{step.label}</div>
              </div>
            ))}
          </div>

          <div class="current-status-message">
            Status atual: <strong>{status()}</strong>
          </div>

          {/* HISTÓRICO DE OBSERVAÇÕES */}
          <div class="obs-box">
            <h3>Observações:</h3>
            <For each={observations()}>
              {(obs) => (
                <div class="observation-message">
                  <p><strong>{obs.author}:</strong> {obs.content}</p>
                  <small>{new Date(obs.createdAt).toLocaleString()}</small>
                </div>
              )}
            </For>
          </div>

          {/* CAMPO PARA NOVA OBSERVAÇÃO */}
          <div class="new-observation">
            <textarea
              placeholder="Escreva uma nova mensagem ou atualização..."
              value={newObs()}
              onInput={(e) => setNewObs((e.target as HTMLTextAreaElement).value)}
            />
            <button onClick={handleNewObservationSubmit}>
              Enviar
            </button>
          </div>
        </div>
      )}

      <div class="alert">{alertMessage()}</div>
    </div>
  );
};
