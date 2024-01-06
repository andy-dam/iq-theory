import React, { useEffect } from "react";
import { Link } from "react-router-dom";

interface Props {
  // props
  clef: string;
  setClef: React.Dispatch<React.SetStateAction<string>>;
}

const Game: React.FC<Props> = ({ clef, setClef }) => {
  useEffect(() => {
    console.log("clef", clef);
    return () => {
      console.log("unmounting");
    };
  }, [clef]);

  return (
    <>
      <div>{clef}</div>
      <Link to="/home">go back to /home</Link>
    </>
  );
};

export default Game;
