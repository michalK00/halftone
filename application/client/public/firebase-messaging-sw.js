// public/firebase-messaging-sw.js
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
    const notificationTitle = payload.notification?.title || 'New Notification';
    const notificationOptions = {
        body: payload.notification?.body || 'You have a new message',
        icon: '/icon-192x192.png',
        data: {
            url: payload.data?.url || '/',
        }
    };

    self.registration.showNotification(notificationTitle, notificationOptions);
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