// Import the Firebase modules
import { initializeApp } from "https://www.gstatic.com/firebasejs/10.8.1/firebase-app.js";
import {
  getAuth,
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
} from "https://www.gstatic.com/firebasejs/10.8.1/firebase-auth.js";
// Initialize Firebase
const firebaseConfig = {
  apiKey: "AIzaSyC3nYh4teEEPd_ORIjna5crO973wF3lUh0",
  authDomain: "todogo-18674.firebaseapp.com",
  projectId: "todogo-18674",
  storageBucket: "todogo-18674.appspot.com",
  messagingSenderId: "1029922052303",
  appId: "1:1029922052303:web:e4c088694d92caeab0bfba",
  measurementId: "G-9WLJRHQPP2",
};
const app = initializeApp(firebaseConfig);
export { app };

// Signup function
function signup() {
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;
  const auth = getAuth(app);
  createUserWithEmailAndPassword(auth, email, password)
    .then((userCredential) => {
      // The user has been signed up
      const user = userCredential.user;
      // After signup, sign in the user, then get the ID token
      user.getIdToken().then((idToken) => {
        // Send the ID token to the server
        fetch("/user/signup", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: idToken, // Send the ID token in the Authorization header
          },
          body: JSON.stringify({ email: email }), // Send the email in the request body
        });
      });
    })
    .catch((error) => {
      // Handle errors
      console.error(error);
    });
}

function login() {
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;
  const auth = getAuth(app);
  signInWithEmailAndPassword(auth, email, password)
    .then((userCredential) => {
      // The user has been signed in
      const user = userCredential.user;
      console.log("user Credential", userCredential);
      console.log("user", user);
      // After sign in, get the ID token
      user.getIdToken().then((idToken) => {
        // Send the ID token to the server
        fetch("/user/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: idToken, // Send the ID token in the Authorization header
          },
          body: JSON.stringify({ email: email }), // Send the email in the request body
        });
      });
    })
    .catch((error) => {
      // Handle errors
      console.error(error);
    });
}

let signupForm = document.getElementById("signupForm");
if (signupForm) {
  signupForm.addEventListener("submit", function (event) {
    event.preventDefault();
    signup();
  });
}

let loginForm = document.getElementById("loginForm");
if (loginForm) {
  loginForm.addEventListener("submit", function (event) {
    event.preventDefault();
    login();
  });
}
