import crypto from 'crypto';

export const generateId = (): string => {
    return crypto.randomUUID();
}