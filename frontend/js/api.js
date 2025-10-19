import { getToken } from './auth.js';

const BASE_URL = 'https://ominous-succotash-x5rw57r6gxw72j7v-8080.app.github.dev/v1';

async function fetchJSON(url, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    const token = getToken();
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    try {
        const response = await fetch(url, { ...options, headers });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ message: response.statusText }));
            throw new Error(`API Error: ${response.status} ${errorData.message || response.statusText}`);
        }

        if (response.status === 204) { // No Content
            return null;
        }

        return response.json();
    } catch (error) {
        console.error('API call failed:', error);
        throw error;
    }
}

export const apiClient = {
    getImages: () => fetchJSON(`${BASE_URL}/images`),
    getImage: (id) => fetchJSON(`${BASE_URL}/images/${id}`),
    getImageSBOMs: (id) => fetchJSON(`${BASE_URL}/images/${id}/sbom`),
    getImageCVEs: (id) => fetchJSON(`${BASE_URL}/images/${id}/cves`),
    getImageVerification: (id) => fetchJSON(`${BASE_URL}/images/${id}/verify`),
};

