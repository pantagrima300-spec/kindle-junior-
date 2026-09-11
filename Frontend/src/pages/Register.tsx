import { useState } from "react";
import {
  ArrowLeft,
  LockKeyhole,
  ShieldCheck,
  Mail,
  UserRound,
} from "lucide-react";
import { Link, useNavigate } from "react-router-dom";
import CircuitBackground from "../components/CircuitBackground";

import { pb } from "../lib/pb";

const Register = () => {
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleRegister = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (isLoading) return;

    setError("");
    if (!name.trim() || !email.trim() || !password.trim() || !confirmPassword.trim()) {
      setError("Please enter your Name, Email, Password and Confirm Password.");
      return;
    }
    
    // Simple email regex
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setError("Please enter a valid email address.");
      return;
    }

    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    setIsLoading(true);
    try {
      // Create user
      await pb.collection("users").create({
        name: name,
        email: email,
        password: password,
        passwordConfirm: confirmPassword,
      });

      // Authenticate newly registered user
      await pb.collection("users").authWithPassword(email, password);
      
      // Redirect to competition page or guidelines
      navigate("/round-1-guidelines");
    } catch (err: any) {
      setError(err.message || "Registration failed. Email might already be registered.");
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
          onSubmit={handleRegister}
        >
          <div className="login-icon">
            <ShieldCheck size={25} />
          </div>

          <div className="login-eyebrow">
            PARTICIPANT REGISTRATION
          </div>

          <h1>
            Create <span>Account.</span>
          </h1>

          <p className="login-description">
            Register to participate in the competition.
          </p>

          <div className="input-group">
            <label>
              FULL NAME
            </label>
            <div className="input-wrapper">
              <UserRound size={16} />
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Enter your Full Name"
              />
            </div>
          </div>

          <div className="input-group">
            <label>
              EMAIL
            </label>
            <div className="input-wrapper">
              <Mail size={16} />
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Enter your Email"
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

          <div className="input-group">
            <label>
              CONFIRM PASSWORD
            </label>
            <div className="input-wrapper">
              <LockKeyhole size={16} />
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Confirm your password"
              />
            </div>
          </div>

          {/* ERROR */}
          {error && (
            <div className="login-error">
              {error}
            </div>
          )}

          {/* REGISTER */}
          <button
            type="submit"
            className="login-button"
            disabled={isLoading}
            style={{ marginBottom: '16px' }}
          >
            {isLoading ? "CREATING ACCOUNT..." : "CREATE ACCOUNT"}
            <span>
              →
            </span>
          </button>
          
          <div className="login-security" style={{ justifyContent: 'center' }}>
            <span>
              Already have an account?{" "}
              <Link to="/login" style={{ color: 'white', textDecoration: 'underline' }}>
                Login
              </Link>
            </span>
          </div>

          <div className="login-security" style={{ marginTop: '16px' }}>
            <ShieldCheck size={12} />
            SECURE PARTICIPANT REGISTRATION
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

export default Register;
