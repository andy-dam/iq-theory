import { User, onAuthStateChanged, getAuth } from "firebase/auth";
import { set } from "firebase/database";
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

interface Props {
  // props
  clef: string;
  setClef: React.Dispatch<React.SetStateAction<string>>;
}

type Note = {
  note: string;
  imgs: string[];
};

type Question = {
  img: string;
  note: string;
  correctAnswer: string;
  answers: string[];
};

const Game: React.FC<Props> = ({ clef, setClef }) => {
  //TODO change back to 3
  const [data, setData] = useState<Note[]>([]);
  const [countdown, setCountdown] = useState<number>(3);
  const [user, setUser] = useState<User | null>(null);
  const [authCheckCompleted, setAuthCheckCompleted] = useState(false);
  const [question, setQuestion] = useState<Question | null>(null);
  const [userAnswer, setUserAnswer] = useState<string>("");
  const [score, setScore] = useState<number>(0);
  const navigate = useNavigate();

  const fetchQuestion = () => {
    const randomIndex: number = Math.floor(Math.random() * data.length);
    const curNote = data[randomIndex];
    const correctAnswer = curNote.note;
    const img = curNote.imgs[Math.floor(Math.random() * curNote.imgs.length)];
    const options = ["A", "B", "C", "D", "E", "F", "G"].filter(
      (note) => note !== correctAnswer
    );
    const wrongAnswers: string[] = [];
    while (wrongAnswers.length < 3) {
      const randomIndex = Math.floor(Math.random() * options.length);
      const randomNote = options[randomIndex];
      if (!wrongAnswers.includes(randomNote)) {
        wrongAnswers.push(randomNote);
        options.filter((note) => note !== randomNote);
      }
    }
    const answers = [...wrongAnswers, correctAnswer];
    //shuffle answers
    answers.sort(() => Math.random() - 0.5);
    console.log(answers);
    console.log(correctAnswer);
    console.log(img);
    setQuestion({ img, note: correctAnswer, correctAnswer, answers });
  };

  // This useEffect is for fetching the data
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`src/quizbank/${clef}.json`);
        const data = await response.json();
        setData(data);
      } catch (error: unknown) {
        console.error(error);
      }
    };
    fetchData();
  }, []);

  // This useEffect is for setting the question
  useEffect(() => {
    if (data.length > 0) {
      fetchQuestion();
    }
  }, [data]);

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

  //handle user answer
  const handleAnswer = (e: React.SyntheticEvent<HTMLButtonElement>) => {
    const answer = e.currentTarget.value;
    setUserAnswer(answer);
    if (answer === question?.correctAnswer) {
      setScore((prev) => prev + 1);
      
    }
    fetchQuestion();
  };

  useEffect(() => {
    if (!clef) {
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
          <div>{score}</div>
          <div><img src={`src/assets/notes/${clef}//${question?.img}.png`}/></div>
          <ul>
            {question?.answers.map((answer) => (
              <li key={`${answer}`}>
                <button className="border w-16 bg-red-100" value={answer} onClick={e => handleAnswer(e)}>
                  {answer}
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}
    </>
  );
};
export default Game;
