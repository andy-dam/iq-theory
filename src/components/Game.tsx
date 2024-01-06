import { User, onAuthStateChanged, getAuth } from "firebase/auth";
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

interface Props {
  // props
  clef: string;
  setClef: React.Dispatch<React.SetStateAction<string>>;
}

const Game: React.FC<Props> = ({ clef, setClef }) => {
  const [countdown, setCountdown] = useState<number>(3);
  const [user, setUser] = useState<User | null>(null);
  const navigate = useNavigate();
  const [authCheckCompleted, setAuthCheckCompleted] = useState(false);

  // Next two useEffects are checks for if the user is logged in
  useEffect(() => {
    const auth = getAuth();
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setUser(user);
      setAuthCheckCompleted(true);
    });

    // Cleanup subscription on unmount
    return () => unsubscribe();
  }, []);

  useEffect(() => {
    if(!clef) {
      navigate("/home");
    }
    if (authCheckCompleted) {
      if (!user) {
        navigate("/home");
      } 
    }
  }, [authCheckCompleted, user]);

  // This useEffect is a countdown timer
  useEffect(() => {
    const interval = setInterval(() => {
      setCountdown((prev) => prev - 1);
    }, 1000);
    return () => {
      clearInterval(interval);
    };
  }, []);


  //refactor these stuff to be components 
  return (
    <>
      {countdown > 0 ? (
        <div>
          <h1>Quiz will start in {countdown} seconds...</h1>
        </div>
      ) : (
        <div>
          <div>{clef}</div>
        </div>
      )}
      
    </>
  );
};
export default Game;
