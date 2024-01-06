import React, { useState, useEffect } from "react";
import { User, onAuthStateChanged, getAuth } from "firebase/auth";
import { Link, useNavigate } from "react-router-dom";
import { auth } from "../config/firebase.js";
import { signOut } from "firebase/auth";

interface Props {
  // props
  setClef: React.Dispatch<React.SetStateAction<string>>;
}

const Homepage: React.FC<Props> = ({ setClef }) => {
  const [authCheckCompleted, setAuthCheckCompleted] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const navigate = useNavigate();

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
    if (authCheckCompleted) {
      if (!user) {
        navigate("/login");
      }
    }
  }, [authCheckCompleted, user]);
  const logOut = async () => {
    try {
      await signOut(auth);
    } catch (err) {
      console.error(err);
    }
  };
  return (
    <>
      {!authCheckCompleted ? (
        <div></div>
      ) : (
        <>
          <div className="flex flex-col gap-10">
            <Link
              to="/game"
              className="flex rounded-full bg-blue-400 max-w-28 justify-center items-center"
              onClick={() => {
                setClef("treble");
              }}
            >
              go to treble
            </Link>
            <Link
              to="/game"
              className="flex rounded-full bg-blue-400 max-w-28 justify-center items-center"
              onClick={() => {
                setClef("bass");
              }}
            >
              go to bass
            </Link>
          </div>
          <Link to="/">go back</Link>
          <button
            className="flex rounded-full bg-blue-400 max-w-28 justify-center items-center"
            onClick={logOut}
          >
            Log Out
          </button>
        </>
      )}
    </>
  );
};

export default Homepage;
