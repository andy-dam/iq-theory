import { initializeApp } from "firebase/app";
import { Auth, getAuth, GoogleAuthProvider } from "firebase/auth";
const firebaseConfig = {
  apiKey: "AIzaSyCPiDqt5Uggxy3UbAcENZ6VQX4-IATfLGM",
  authDomain: "iq-theory.firebaseapp.com",
  projectId: "iq-theory",
  storageBucket: "iq-theory.appspot.com",
  messagingSenderId: "199895688692",
  appId: "1:199895688692:web:01f72455789b0f06cb85a1",
  measurementId: "G-FTV503Y5SR"
};

const app = initializeApp(firebaseConfig);
export const auth: Auth = getAuth(app);
export const googleProvider: GoogleAuthProvider = new GoogleAuthProvider(); 