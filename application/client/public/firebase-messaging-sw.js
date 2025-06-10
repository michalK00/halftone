importScripts('https://www.gstatic.com/firebasejs/9.0.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/9.0.0/firebase-messaging-compat.js');

const firebaseConfig = {
    apiKey: "AIzaSyCAgCjC_lepDLSDjbpsIyYCzjM21OtudXY",
    authDomain: "halftone-f9f6b.firebaseapp.com",
    projectId: "halftone-f9f6b",
    storageBucket: "halftone-f9f6b.firebasestorage.app",
    messagingSenderId: "943968764365",
    appId: "1:943968764365:web:3f73c823d36301a6a70484",
    measurementId: "G-Y0CCNQP7HK"
};

firebase.initializeApp(firebaseConfig);
const messaging = firebase.messaging();

messaging.onBackgroundMessage((payload) => {
    console.log('Received background message:', payload);

    // Handle both notification and data-only messages
    const title = payload.notification?.title || payload.data?.title || 'New Notification';
    const body = payload.notification?.body || payload.data?.body || 'You have a new message';

    const notificationOptions = {
        body: body,
        icon: payload.notification?.icon || '/icon-192x192.png',
        badge: '/badge-72x72.png', // Optional
        tag: payload.data?.tag || 'default', // Prevents duplicate notifications
        data: {
            url: payload.data?.url || '/',
            clickAction: payload.data?.click_action || payload.fcm_options?.link || '/'
        },
        requireInteraction: false,
        silent: false
    };

    return self.registration.showNotification(title, notificationOptions);
});

self.addEventListener('notificationclick', (event) => {
    event.notification.close();

    const urlToOpen = event.notification.data?.url || '/';

    event.waitUntil(
        clients.matchAll({ type: 'window' }).then((clientList) => {
            for (const client of clientList) {
                if (client.url === urlToOpen && 'focus' in client) {
                    return client.focus();
                }
            }
            if (clients.openWindow) {
                return clients.openWindow(urlToOpen);
            }
        })
    );
});