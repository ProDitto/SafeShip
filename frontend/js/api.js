const BASE_URL = '/v1';

async function fetchJSON(url, options = {}) {
    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ message: response.statusText }));
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }
        return await response.json();
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

