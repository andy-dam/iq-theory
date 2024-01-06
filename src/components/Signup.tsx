import { auth } from "../config/firebase.js";
import { createUserWithEmailAndPassword } from "firebase/auth";
import { useRef, useState, useEffect } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";
import { User, onAuthStateChanged, getAuth } from "firebase/auth";
import Logo from "../assets/logos/black-box.svg";

const Signup = () => {
  const email = useRef("");
  const password = useRef("");
  const confirmPassword = useRef("");
  const [home, setHome] = useState(false);
  const [loading, setLoading] = useState(false);
  const [authCheckCompleted, setAuthCheckCompleted] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (password.current !== confirmPassword.current) {
      alert("Passwords do not match");
      return;
    }
    try {
      setLoading(true);
      await createUserWithEmailAndPassword(
        auth,
        email.current,
        password.current
      );
    } catch (err: unknown) {
      console.error(err);
    }
    setLoading(false);
  };
  
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
              Signup
            </h1>
            <div className="flex flex-col mx-auto">
              <input
                className="border-2 border-gray-500 rounded-md p-2 m-2 w-52 max-w-56 placeholder:text-gray-500"
                placeholder="Email"
                onChange={(e) => (email.current = e.target.value)}
              />
              <input
                className="border-2 border-gray-500 rounded-md p-2 m-2 w-52 max-w-56 placeholder:text-gray-500"
                placeholder="Password"
                onChange={(e) => (password.current = e.target.value)}
                type="password"
              />
              <input
                className="border-2 border-gray-500 rounded-md p-2 m-2 mb-4 w-52 max-w-56 placeholder:text-gray-500"
                placeholder="Confirm Password"
                onChange={(e) => (confirmPassword.current = e.target.value)}
                type="password"
              />
            </div>
            <div className="max-w-64">
              <button
                className=" text-white w-52  bg-gray-800 rounded-md p-2 m-2 max-w-52 hover:opacity-50"
                type="submit"
                disabled={loading}
              >
                Sign Up
              </button>
              <h2 className="">
                Already have an account or want to use Google? Log in below!
              </h2>
              <Link to="/login" className="">
                <button
                  className=" border-2 border-gray-400 rounded-md p-2 m-2 max-w-52 hover:opacity-50"
                  disabled={loading}
                >
                  Go to Log In
                </button>
              </Link>
            </div>
          </form>
        </div>
      )}
    </>
  );
};

export default Signup;
