import { apiClient } from './api.js';
import { renderImageTable } from './components/imageTable.js';
import { renderImageDetail } from './components/imageDetail.js';

const app = document.getElementById('app');

const showLoading = () => {
    app.innerHTML = '<p>Loading...</p>';
};

const showError = (message) => {
    app.innerHTML = `<p style="color: red;">Error: ${message}</p>`;
};

const showImageListView = async () => {
    showLoading();
    try {
        const images = await apiClient.getImages();
        app.innerHTML = renderImageTable(images);
    } catch (error) {
        showError(error.message);
    }
};

const showImageDetailView = async (id) => {
    showLoading();
    try {
        const [image, sboms, cves, verification] = await Promise.all([
            apiClient.getImage(id),
            apiClient.getImageSBOMs(id),
            apiClient.getImageCVEs(id),
            apiClient.getImageVerification(id),
        ]);
        app.innerHTML = renderImageDetail({ image, sboms, cves, verification });
    } catch (error) {
        showError(error.message);
    }
};

const router = () => {
    const hash = window.location.hash;
    const match = hash.match(/^#\/images\/(\d+)$/);

    if (match) {
        const imageId = match[1];
        showImageDetailView(imageId);
    } else {
        showImageListView();
    }
};

window.addEventListener('hashchange', router);
window.addEventListener('DOMContentLoaded', router);
