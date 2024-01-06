import React from "react";
import { Link } from "react-router-dom";
import { auth } from "../config/firebase.js";
import { signOut } from "firebase/auth";

interface Props {
  // props
  clef: string;
  setClef: React.Dispatch<React.SetStateAction<string>>;
}

const Homepage: React.FC<Props> = ({ clef, setClef }) => {
  const logOut = async () => {
    try {
      await signOut(auth);
    } catch (err) {
      console.error(err);
    }
  };
  return (
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
  );
};

export default Homepage;
