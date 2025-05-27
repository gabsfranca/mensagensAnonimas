import { createSignal, JSX } from 'solid-js'; 
import { sendAnonymousMessage } from '../services/messageServices';
import { Spinner } from './spinner';
import { generateId } from '../utils/idGenerator';
import { send } from 'vite';

export const MessageForm = (): JSX.Element => {
    //createSignal é tipo o useState do React 
    const [message, setMessage] = createSignal('');
    const [isLoading, setIsLoading] = createSignal(false);
    const [alertMessage, setAlertMessage] = createSignal('');
    const [files, setFiles] = createSignal<File[]>([]);    

    const handleMessageChange = (e: Event) => {
        const target = e.target as HTMLTextAreaElement;
        setMessage(target.value);
    };

    const handleFileChange = (e: Event) => {
      const target = e.target as HTMLInputElement;
      if (target.files) {
        setFiles(Array.from(target.files));
      }
    };

    const handleSubmit = async () => {
        const trimmedMessage = message().trim();
        setAlertMessage('');

        if (!trimmedMessage) {
            setAlertMessage('Por favor, digite uma mensagem');
            return; 
        }

        setIsLoading(true);

        try {

            const formData = new FormData();
            formData.append('content', trimmedMessage);

            files().forEach((file, index) => {
              formData.append('files', file);
            });

            const result = await sendAnonymousMessage(formData);

            if (result.success) {
              alert(
                `Mensagem enviada com sucesso!\nID copiado para sua área de transferência! Guarde ele, é importante.`
              );
              if (result.id) {
                if (navigator.clipboard && navigator.clipboard.writeText) {
                    navigator.clipboard.writeText(result.id);
                } else {
                    setAlertMessage('erro ao copiar texto para a área de trabalho')
                }
              }
              setMessage(`copie seu código: ${result.id}`);
              setFiles([]);
              setAlertMessage('');
            } else {
              setAlertMessage(`Erro ao enviar mensagem: ${result.error}`);
            }
        } catch (e) {
            console.error('Err: ', e);
            setAlertMessage('erro ao conectar com o sv');
        } finally{
            setIsLoading(false);
        }
    };

     return (
        <div class="container">
            <h1>Mensagem Anônima</h1>
            <textarea
                id="messageInput"
                value={message()}
                onInput={handleMessageChange}
                placeholder="Digite sua mensagem aqui..."
            />
            
            {/* Input para upload de arquivos */}
            <div class="file-upload">
                <label for="fileInput"></label>
                <input
                    type="file"
                    id="fileInput"
                    multiple
                    onChange={handleFileChange}
                    accept="image/*, video/*, audio/*"
                />

                <div class="file-list">
                    {files().map(file => (
                        <div class="file-item">
                            {file.name} ({Math.round(file.size / 1024)} KB)
                        </div>
                    ))}
                </div>
            </div>

            <button
                onClick={handleSubmit}
                disabled={isLoading()}
            >
                <span class="button-content">
                    <span>{isLoading() ? 'Enviando...' : 'Enviar Mensagem'}</span>
                    <Spinner show={isLoading()} />
                </span>
            </button>
            <div class="alert">{alertMessage()}</div>
        </div>
    );
};