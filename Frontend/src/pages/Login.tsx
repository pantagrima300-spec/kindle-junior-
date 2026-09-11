import { useState } from "react";
import {
  ArrowLeft,
  LockKeyhole,
  ShieldCheck,
  UserRound,
} from "lucide-react";
import { Link, useNavigate } from "react-router-dom";
import CircuitBackground from "../components/CircuitBackground";

import { pb } from "../lib/pb";

const Login = () => {
  const navigate = useNavigate();

  const [userId, setUserId] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (isLoading) return;

    setError("");
    if (!userId.trim() || !password.trim()) {
      setError("Please enter your User ID and Password.");
      return;
    }
    
    setIsLoading(true);
    try {
      await pb.collection("users").authWithPassword(userId, password);
      navigate("/round-1-guidelines");
    } catch (err: any) {
      setError(err.message || "Invalid credentials");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="login-page">

      <CircuitBackground />

      <div className="login-overlay" />

      <div className="login-top">

        <Link to="/" className="back-button">
          <ArrowLeft size={15} />
          BACK
        </Link>

        <div className="login-event">
          <span />
          KINDLE JUNIOR <b>5.0</b>
        </div>

      </div>
      <section className="login-container">

        <form
          className="login-card"
          onSubmit={handleLogin}
        >

          <div className="login-icon">
            <ShieldCheck size={25} />
          </div>


          <div className="login-eyebrow">
            ROUND 01 · PARTICIPANT ACCESS
          </div>


          <h1>
            Welcome <span>Back.</span>
          </h1>


          <p className="login-description">
            Enter the credentials provided to you
            to access Round 01.
          </p>
          <div className="input-group">

            <label>
              USER ID
            </label>

            <div className="input-wrapper">

              <UserRound size={16} />

              <input
                type="text"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                placeholder="Enter your User ID"
              />

            </div>

          </div>
          <div className="input-group">

            <label>
              PASSWORD
            </label>

            <div className="input-wrapper">

              <LockKeyhole size={16} />

              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter your password"
              />

            </div>

          </div>


          {/* ERROR */}
          {error && (
            <div className="login-error">
              {error}
            </div>
          )}


          {/* LOGIN */}
          <button
            type="submit"
            className="login-button"
            disabled={isLoading}
            style={{ marginBottom: '16px' }}
          >
            {isLoading ? "AUTHENTICATING..." : "ENTER ROUND 01"}

            <span>
              →
            </span>
          </button>
          
          <div className="login-security" style={{ justifyContent: 'center' }}>
            <span>
              Don't have an account?{" "}
              <Link to="/register" style={{ color: 'white', textDecoration: 'underline' }}>
                Register
              </Link>
            </span>
          </div>

          <div className="login-security" style={{ marginTop: '16px' }}>
            <ShieldCheck size={12} />
            SECURE PARTICIPANT ACCESS
          </div>

        </form>


        {/* BOTTOM INFO */}
        <div className="login-meta">

          <span>ROUND 01</span>

          <i />

          <span>40 QUESTIONS</span>

          <i />

          <span>40 MINUTES</span>

        </div>

      </section>

    </main>
  );
};

export default Login;