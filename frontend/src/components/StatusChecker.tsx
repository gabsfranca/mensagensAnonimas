import { createSignal, JSX, } from "solid-js";
import { checkReportStatus } from "../services/statusServices";
import { Spinner } from "./spinner";

export const StatusChecker = (): JSX.Element => {
    const [uuid, setUuid] = createSignal('');
    const [status, setStatus] = createSignal('');
    const [isLoading, setIsLoading] = createSignal(false);
    const [alertMessage, setAlertMessage] = createSignal('');

    const statusSteps = [
        { id: 1, label: 'Recebido', key: 'recebido' },
        { id: 2, label: 'Em análise', key: 'em análise' },
        { id: 3, label: 'Concluído', key: 'concluído' }
    ];

    const handleUuidChange = (e: Event) => {
        const target = e.target as HTMLInputElement;
        setUuid(target.value.trim());
    };

    const handleCheckStatus = async () => {
        setAlertMessage('');
        setStatus('');

        if (!uuid()) {
            setAlertMessage('Port favor digite o ID da denuncia');
            return;
        }

        if (!isValidUuid(uuid())) {
            setAlertMessage('ID invalido');
            return;
        } 

        setIsLoading(true);

        try {
            const result = await checkReportStatus(uuid());

            if (result.success) {
                setStatus(result.status!);
                setAlertMessage('');
            } else {
                setAlertMessage(result.error || 'Erro ao buscar status');
            }
        } catch (e) {
            console.error('erro: ', e);
            setAlertMessage('erro ao conectar com o sv');
        } finally {
            setIsLoading(false);
        }
    };

    const isValidUuid = (uuid: string) => {
        const regex = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/;
        return regex.test(uuid);
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
                    placeholder="Digite o UUID da denúncia"
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
      </div>
    )}

    <div class="alert">{alertMessage()}</div>
  </div>
);
}