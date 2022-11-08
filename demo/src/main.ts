import './style.css'

import { initializeApp } from "firebase/app";
import { initializeAppCheck, ReCaptchaV3Provider, getToken } from "firebase/app-check";

const firebaseConfig = {
  apiKey: "AIzaSyDmIi4AHW2UoUDUq5rFfovSaXS5LX8EeQw",
  authDomain: "holidays-1170.firebaseapp.com",
  projectId: "holidays-1170",
  storageBucket: "holidays-1170.appspot.com",
  messagingSenderId: "394606266940",
  appId: "1:394606266940:web:9caf08b5e1210ed304d139",
};

// Initialize Firebase.
const app = initializeApp(firebaseConfig);

// Initialize AppCheck.
const appCheck = initializeAppCheck(app, {
  provider: new ReCaptchaV3Provider(import.meta.env.VITE_RECAPTCHA_KEY),
  isTokenAutoRefreshEnabled: true
});

// Grab an AppCheck token.
const appCheckToken = await getToken(appCheck).then(t => t.token);

// Call our backend to convert the AppCheck token into a JWT.
const jwt = await fetch(import.meta.env.VITE_TOKEN_BACKEND, {
  headers: {
    'X-Firebase-AppCheck': appCheckToken,
  }
}).then((data) => data.text());

// Call the Routes API!
const response = await fetch('https://routes.googleapis.com/directions/v2:computeRoutes', {
  method: 'POST',
  headers: {
      'Authorization': `Bearer ${jwt}`, // Pass our JWT!
      'Content-Type': 'application/json',
      'X-Goog-FieldMask': 'routes.duration,routes.distanceMeters,routes.polyline.encodedPolyline', 
  },
  body: JSON.stringify({
    "origin":{
      "location":{
        "latLng":{
          "latitude": 37.419734,
          "longitude": -122.0827784
        }
      }
    },
    "destination":{
      "location":{
        "latLng":{
          "latitude": 37.417670,
          "longitude": -122.079595
        }
      }
    },
    "travelMode": "DRIVE"}),
});

console.log(response);
