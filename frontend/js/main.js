import { apiClient } from './api.js';
import { renderImageTable } from './components/imageTable.js';
import { renderImageDetail } from './components/imageDetail.js';
import { isAuthenticated, login, logout } from './auth.js';

const app = document.getElementById('app');
const authControls = document.getElementById('auth-controls');

const showLoading = () => {
    app.innerHTML = '<div class="card">Loading...</div>';
};

const showError = (message) => {
    app.innerHTML = `<div class="card" style="color: red;">Error: ${message}</div>`;
};

const showLoginView = () => {
    app.innerHTML = `
        <div class="card">
            <h2>Please Log In</h2>
            <p>You must be logged in to view the image catalog.</p>
        </div>
    `;
};

const showImageListView = async () => {
    showLoading();
    try {
        const images = await apiClient.getImages();
        app.innerHTML = renderImageTable(images);
    } catch (error) {
        showError(error.message);
        if (error.message.includes('401')) {
            showLoginView();
        }
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
    if (!isAuthenticated()) {
        showLoginView();
        return;
    }

    const hash = window.location.hash;
    if (hash.startsWith('#/images/')) {
        const id = hash.split('/')[2];
        showImageDetailView(id);
    } else {
        showImageListView();
    }
};

const renderAuthUI = () => {
    if (isAuthenticated()) {
        authControls.innerHTML = '<button id="logout-btn">Logout</button>';
        document.getElementById('logout-btn').addEventListener('click', () => {
            logout();
            renderAuthUI();
            router();
        });
    } else {
        authControls.innerHTML = '<button id="login-btn">Login</button>';
        document.getElementById('login-btn').addEventListener('click', () => {
            login();
            renderAuthUI();
            router();
        });
    }
};

window.addEventListener('hashchange', router);
window.addEventListener('DOMContentLoaded', () => {
    renderAuthUI();
    router();
});
