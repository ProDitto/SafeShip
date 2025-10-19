const TOKEN_KEY = 'authToken';

export function login() {
    // In a real OIDC flow, this would redirect to an identity provider.
    // For this mock, we just set a hardcoded token.
    localStorage.setItem(TOKEN_KEY, 'mock-jwt-token');
    console.log('Mock login successful.');
}

export function logout() {
    localStorage.removeItem(TOKEN_KEY);
    console.log('Logged out.');
}

export function getToken() {
    return localStorage.getItem(TOKEN_KEY);
}

export function isAuthenticated() {
    return !!getToken();
}

