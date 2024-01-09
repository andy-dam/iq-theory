// import { useState } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import { auth } from "./config/firebase";
import { useAuthState } from "react-firebase-hooks/auth";
import Homepage from "./components/Homepage.tsx";
import Landing from "./components/Landing.tsx";
import Login from "./components/Login.tsx";
import Game from "./components/Game.tsx";
import { useState } from "react";
import Signup from "./components/Signup.tsx";
const App: React.FC = () => {
  const [user] = useAuthState(auth);
  const [clef, setClef] = useState<string>("treble");

  return (
    <>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/home" element={<Homepage setClef={setClef} />} />
        <Route
          path="/login"
          element={user ? <Navigate to="/home" /> : <Login />}
        />
        <Route
          path="/signup"
          element={user ? <Navigate to="/home" /> : <Signup />}
        />
        <Route path="/game" element={<Game clef={clef} setClef={setClef} />} />
        <Route path="*" element={<div>404</div>} />
      </Routes>
    </>
  );
};

export default App;
