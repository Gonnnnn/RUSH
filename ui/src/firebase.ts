import { initializeApp } from 'firebase/app';
import { getAuth, GoogleAuthProvider } from 'firebase/auth';

// Tested manually, only `apiKey` and `authDomain` are required for Google sign-in.
const firebaseConfig = {
  // Firebase API key is an identifier rather than a secret.
  // https://firebase.google.com/docs/projects/api-keys#general-info
  // Security should be handled by restricting access of the API key.
  // Check https://firebase.google.com/docs/projects/api-keys#faq-required-apis-for-restricted-firebase-api-key
  // to see which APIs are required for Firebase authentication.
  // Monitoring the usage of the API key using services like app-check helps too.
  // ex) https://firebase.google.com/docs/app-check/web/recaptcha-enterprise-provider
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const provider = new GoogleAuthProvider();

export { auth, provider };
