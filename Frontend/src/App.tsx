
import { BrowserRouter, Routes, Route } from "react-router-dom";

import Home from "./pages/Home";
import Login from "./pages/Login";

import Round1Guidelines from "./pages/Round1Guidelines";
import Round1 from "./pages/Round1";
import Round1Result from "./pages/Round1Result";

import Round2Guidelines from "./pages/Round2Guidlines";
import LanguageSelection from "./pages/LanguageSelection";
import CodingArena from "./pages/CodingArena";
import Round2Result from "./pages/Round2Result";
import Round3Guidelines from "./pages/Round3Guidelines";
import Register from "./pages/Register";
import Round2FinalPage from "./pages/Round2FinalPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>

        {/* HOME */}
        <Route
          path="/"
          element={<Home />}
        />

        {/* REGISTER */}
        <Route
          path="/register"
          element={<Register />}
        />

        {/* LOGIN */}
        <Route
          path="/login"
          element={<Login />}
        />

        {/* LANGUAGE SELECTION — BEFORE ROUND 1 */}
        <Route
          path="/language-selection"
          element={<LanguageSelection />}
        />

        {/* ROUND 1 GUIDELINES */}
        <Route
          path="/round-1-guidelines"
          element={<Round1Guidelines />}
        />

        {/* ROUND 1 */}
        <Route
          path="/round-1"
          element={<Round1 />}
        />

        {/* ROUND 1 RESULT */}
        <Route
          path="/round-1-result"
          element={<Round1Result />}
        />

        {/* ROUND 2 GUIDELINES */}
        <Route
          path="/round-2-guidelines"
          element={<Round2Guidelines />}
        />

        {/* ROUND 2 CODING LANGUAGE SELECTION */}
        <Route
          path="/round-2/language"
          element={<LanguageSelection />}
        />

        {/* ROUND 2 CODING ARENA */}
        <Route
          path="/round-2/coding-arena"
          element={<CodingArena />}
        />

        {/* ROUND 2 RESULT */}
        <Route
          path="/round-2-result"
          element={<Round2Result />}
        />

        {/* ROUND 3 GUIDELINES */}
        <Route
          path="/round-3-guidelines"
          element={<Round3Guidelines />}
        />

        {/* ROUND 2 TIME UP */}
        <Route
          path="/round-2/time-up"
          element={<Round2FinalPage status="time_up" />}
        />

        {/* ROUND 2 COMPLETED */}
        <Route
          path="/round-2/completed"
          element={<Round2FinalPage status="completed" />}
        />

      </Routes>
    </BrowserRouter>
  );
}

export default App;

