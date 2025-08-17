import React, { useState } from "react";
import { Link } from "react-router-dom";
import { IoMenuSharp, IoCloseSharp } from "react-icons/io5";
import Logo from "../assets/logos/white-horizontal.svg";
const Landing: React.FC = () => {
  const [showMenu, setShowMenu] = useState<boolean>(false);
  const [menuIcon, setMenuIcon] = useState<boolean>(false);
  return (
    <>
      <header className="bg-gray-500 font-mono">
        <nav className="flex justify-evenly items-center w-[92%] mx-auto py-3">
          <div>
            <img src={Logo} className="w-56" />
          </div>
          <div
            className={
              "text-white text-3xl md:static absolute md:min-h-fit  min-h-[20vh] left-0 w-full flex items-center px-5 md:w-auto " +
              (showMenu ? "top-[9%]" : "top-[-100%]")
            }
          >
            <ul className="flex md:flex-row flex-col md:items-center md:gap-[6vw] gap-8">
              <li className="hover:text-gray-300">
                <Link className="" to="/login">
                  About
                </Link>
              </li>
              <li className={"hover:text-gray-300"}>
                <Link to="/login">Team</Link>
              </li>
            </ul>
          </div>
          <div className="flex items-center gap-6">
            <Link
              to="/login"
              className="text-xl bg-gray-800 text-white px-5 py-2 rounded-full hover:bg-gray-600"
            >
              Sign In
            </Link>
            <button
              className="text-3xl cursor-pointer md:hidden"
              onClick={() => {
                setShowMenu(!showMenu);
                setMenuIcon(!menuIcon);
              }}
            >
              {menuIcon ? <IoCloseSharp /> : <IoMenuSharp />}
            </button>
          </div>
        </nav>
      </header>
    </>
  );
};

export default Landing;
