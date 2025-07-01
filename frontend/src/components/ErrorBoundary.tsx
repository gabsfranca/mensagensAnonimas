import { Show } from 'solid-js';
import '../styles/components/ErrorBoundary.css'

interface ErrorBoundaryProps {
  error: string;
  onRetry?: () => void;
  showLoginButton?: boolean;
}

export const ErrorBoundary = (props: ErrorBoundaryProps) => {
  return (
    <Show when={props.error}>
      <div class="alert error">
        {props.error}
        <Show when={props.showLoginButton && props.error.includes('autenticado')}>
          <button 
            onClick={() => window.location.href = '/login'}
            style="margin-left: 10px; padding: 5px 10px; background: #007bff; color: white; border: none; border-radius: 3px; cursor: pointer;"
          >
            Ir para Login
          </button>
        </Show>
        <Show when={props.onRetry}>
          <button 
            onClick={props.onRetry}
            style="margin-left: 10px; padding: 5px 10px; background: #28a745; color: white; border: none; border-radius: 3px; cursor: pointer;"
          >
            Tentar Novamente
          </button>
        </Show>
      </div>
    </Show>
  );
};