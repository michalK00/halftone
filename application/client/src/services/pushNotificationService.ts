import { initializeApp } from 'firebase/app';
import { getMessaging, getToken, onMessage } from 'firebase/messaging';
import {api} from "@/api";

interface FirebaseConfig {
    apiKey: string;
    authDomain: string;
    projectId: string;
    storageBucket: string;
    messagingSenderId: string;
    appId: string;
    measurementId: string;
}

interface SubscriptionRequest {
    token: string;
}

const firebaseConfig: FirebaseConfig = {
    apiKey: "AIzaSyCAgCjC_lepDLSDjbpsIyYCzjM21OtudXY",
    authDomain: "halftone-f9f6b.firebaseapp.com",
    projectId: "halftone-f9f6b",
    storageBucket: "halftone-f9f6b.firebasestorage.app",
    messagingSenderId: "943968764365",
    appId: "1:943968764365:web:3f73c823d36301a6a70484",
    measurementId: "G-Y0CCNQP7HK"
};

const VAPID_KEY = "BAmp03J42z4bOTSjkfhngEyLJDnPl1GjX-8SVCnJEfSIRaGWAoYjImqmJ1hGE335T8PaGURVqKYxE-A-kdk8c7A";

class PushNotificationService {
    private readonly app;
    private readonly messaging;
    private currentToken: string | null = null;

    constructor() {
        this.app = initializeApp(firebaseConfig);
        this.messaging = getMessaging(this.app);
    }

    isSupported(): boolean {
        return 'serviceWorker' in navigator && 'PushManager' in window;
    }

    async requestPermission(): Promise<string | null> {
        try {
            const permission = await Notification.requestPermission();

            if (permission === 'granted') {
                if ('serviceWorker' in navigator) {
                    await navigator.serviceWorker.register('/firebase-messaging-sw.js');
                }

                const token = await getToken(this.messaging, {
                    vapidKey: VAPID_KEY,
                });

                if (token) {
                    this.currentToken = token;
                    await this.registerWithBackend(token);
                    return token;
                }
            }
            return null;
        } catch (error) {
            console.error('Error getting permission or token:', error);
            throw error;
        }
    }

    private async registerWithBackend(token: string): Promise<void> {
        const request: SubscriptionRequest = {
            token
        };

        try {
            const response = await api.post('/api/v1/push/subscribe', JSON.stringify(request));
            return response.data;
        } catch (error) {
            throw new Error('Failed to register token with backend');
        }
    }

    onMessageListener(callback: (payload: any) => void): void {
        onMessage(this.messaging, callback);
    }

    getCurrentToken(): string | null {
        return this.currentToken;
    }
}

export default new PushNotificationService();