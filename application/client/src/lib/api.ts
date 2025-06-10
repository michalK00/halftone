import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';

const api = axios.create({
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

api.interceptors.response.use(
    (response: AxiosResponse) => response,
    async (error: AxiosError) => {
        const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };

        if (
            error.response?.status === 401 &&
            !originalRequest._retry &&
            localStorage.getItem('refreshToken')
        ) {
            originalRequest._retry = true;

            try {
                const refreshToken = localStorage.getItem('refreshToken');
                const response = await api.post(
                    '/auth/refresh-token',
                    { refresh_token: refreshToken }
                );

                const { id_token, refresh_token } = response.data;
                localStorage.setItem('token', id_token);
                localStorage.setItem('refreshToken', refresh_token);

                if (originalRequest.headers) {
                    originalRequest.headers.Authorization = `Bearer ${id_token}`;
                } else {
                    originalRequest.headers = { Authorization: `Bearer ${id_token}` };
                }

                return axios(originalRequest);
            } catch (refreshError) {
                localStorage.removeItem('token');
                localStorage.removeItem('refreshToken');

                window.location.href = '/signin';
                return Promise.reject(refreshError);
            }
        }

        if (error.response?.data) {
            return Promise.reject({
                status: error.response.status,
                message: 'An error occurred',
                data: error.response.data
            });
        }

        return Promise.reject(error);
    }
);

export default api;