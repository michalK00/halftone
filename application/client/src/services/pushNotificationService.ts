import { initializeApp } from 'firebase/app';
import { getMessaging, getToken, onMessage } from 'firebase/messaging';

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
    userId: string;
    userAttributes?: Record<string, string>;
}

// Firebase configuration
const firebaseConfig: FirebaseConfig = {

    apiKey: "AIzaSyDvgVZNmJpawxPchhGJjmcvhqeodfXpaTo",
    authDomain: "halftone-63e5a.firebaseapp.com",
    projectId: "halftone-63e5a",
    storageBucket: "halftone-63e5a.firebasestorage.app",
    messagingSenderId: "430276000092",
    appId: "1:430276000092:web:c69cc5976610b79ab7e999",
    measurementId: "G-F8NED80L3D"
};

const VAPID_KEY = "BAF1rvan9p5Z0_zMhR-TV1indHpYhWKc8Z_twPonYRkRZPiaKMctrp1YFRxsvHiG0FUhVLpjcgFGw_cmpGkK1pA";

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

    async requestPermission(userId: string): Promise<string | null> {
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
                    await this.registerWithBackend(token, userId);
                    return token;
                }
            }
            return null;
        } catch (error) {
            console.error('Error getting permission or token:', error);
            throw error;
        }
    }

    // Register token with backend
    private async registerWithBackend(token: string, userId: string): Promise<void> {
        const request: SubscriptionRequest = {
            token,
            userId,
            userAttributes: {
                registeredAt: new Date().toISOString(),
            }
        };

        const response = await fetch('/api/push/subscribe', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(request),
        });

        if (!response.ok) {
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