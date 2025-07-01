const URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

export const getObservations = async (shordId: string) => {
    try {
        const res = await fetch(`${URL}/reports/${shordId}/observations`);
        return await res.json();
    } catch (e) {
        console.error('erro ao pegar obs da mensagem!');
    }
}

export const postObservation = async (shortId: string, content: string) => {
    const res = await fetch(`${URL}/reports/${shortId}/observations`, {
        method: 'POST', 
        headers: { "Content-Type": "application/json" }, 
        body: JSON.stringify({ content }),
    });
    return await res.json();
};