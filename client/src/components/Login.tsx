import { auth, googleProvider } from "../config/firebase.js";
import { signInWithEmailAndPassword, signInWithPopup } from "firebase/auth";
import { useRef, useState, useEffect } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";
import Logo from "../assets/logos/black-box.svg";
import { User, onAuthStateChanged, getAuth } from "firebase/auth";

const Login = () => {
  const email: React.MutableRefObject<string> = useRef("");
  const password: React.MutableRefObject<string> = useRef("");
  const [loading, setLoading] = useState<boolean>(false);
  const [home, setHome] = useState<boolean>(false);
  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    try {
      await signInWithEmailAndPassword(auth, email.current, password.current);
    } catch (err) {
      console.error(err);
      alert("Incorrect email or password");
    }
    setLoading(false);
  };

  const signInWithGoogle = async () => {
    try {
      await signInWithPopup(auth, googleProvider);
    } catch (err) {
      console.error(err);
    }
  };
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

  return (
    <>
      {!authCheckCompleted ? (
        <div></div>
      ) : (
        <div className="flex flex-col justify-center items-center text-center">
          <img
            src={Logo}
            alt="logo"
            className="cursor-pointer w-56 mx-auto my-10"
            onClick={() => setHome(true)}
          />
          {home && <Navigate to="/" />}
          <form
            className="flex flex-col max-w-[75rem] my-0"
            onSubmit={handleSubmit}
          >
            <h1 className="justify-center items-center mx-auto text-3xl pb-3 w-56">
              Login
            </h1>
            <div className="flex flex-col mx-auto">
              <input
                className="border-2 border-gray-500 rounded-md p-2 m-2 w-52 max-w-56 placeholder:text-gray-500 "
                placeholder="Email"
                onChange={(e) => (email.current = e.target.value)}
              />
              <input
                className="border-2 border-gray-500 rounded-md p-2 m-2 mb-4 w-52 max-w-56 placeholder:text-gray-500"
                placeholder="Password"
                onChange={(e) => (password.current = e.target.value)}
                type="password"
              />
            </div>
            <div className="flex flex-col items-center">
              <button
                className=" text-white w-52  bg-gray-800 rounded-md p-2 m-2 max-w-52 hover:opacity-50"
                type="submit"
                disabled={loading}
              >
                Log In
              </button>
              <button
                className="text-white bg-gray-800 w-52 rounded-md p-2 mb-4 max-w-56 hover:opacity-50"
                onClick={signInWithGoogle}
                disabled={loading}
                type="button"
              >
                Sign In With Google
              </button>
              <h2 className="max-w-64">
                Don&rsquo;t have an account? Log in with Google or Sign up
                below!
              </h2>
              <Link to="/signup" className="">
                <button
                  className="border-2 border-gray-400 rounded-md p-2 m-2 max-w-52 hover:opacity-50"
                  disabled={loading}
                >
                  Go to Sign Up
                </button>
              </Link>
            </div>
          </form>
        </div>
      )}
    </>
  );
};

export default Login;
