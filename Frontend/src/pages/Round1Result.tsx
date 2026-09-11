import React, { useEffect, useState } from "react";
import {
  CheckCircle2,
  XCircle,
  CircleHelp,
  Clock3,
  Trophy,
  ArrowRight,
  RotateCcw,
  ShieldCheck,
} from "lucide-react";
import { Link, useNavigate } from "react-router-dom";
import CircuitBackground from "../components/CircuitBackground";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
const RESULT_API = `${API_BASE_URL}/api/round1/result`;

// Define the interface for the expected result data
interface Round1ResultData {
  score: number;
  total: number;
  correct: number;
  incorrect: number;
  unanswered: number;
  timeTaken: string;
}

// Typed helper function outside the component
const formatTime = (seconds: number | string): string => {
  const safeSeconds = Math.max(0, Number(seconds) || 0);

  const minutes = Math.floor(safeSeconds / 60);
  const remainingSeconds = safeSeconds % 60;

  return `${String(minutes).padStart(2, "0")}:${String(
    remainingSeconds
  ).padStart(2, "0")}`;
};

const Round1Result: React.FC = () => {
  const navigate = useNavigate();

  // Typed state variables
  const [result, setResult] = useState<Round1ResultData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    let isMounted = true; 

    const fetchResult = async () => {
      try {
        setLoading(true);
        setError("");

        const token = localStorage.getItem("token");

        if (!token) {
          if (isMounted) navigate("/login", { replace: true });
          return;
        }

        const response = await fetch(RESULT_API, {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });

        if (!response.ok) {
          if (response.status === 401) {
            localStorage.removeItem("token");
            if (isMounted) navigate("/login", { replace: true });
            return;
          }

          throw new Error("Unable to load your Round 01 result.");
        }

        const data = await response.json();

        if (!isMounted) return;

        const total =
          Number(data.total) ||
          Number(data.total_questions) ||
          40;

        const correct = Number(data.correct) || 0;
        const incorrect = Number(data.incorrect) || 0;
        const unanswered = Number(data.unanswered) || 0;

        const score =
          data.score !== undefined
            ? Number(data.score)
            : correct;

        const timeTakenSeconds =
          Number(data.time_taken_seconds) ||
          Number(data.timeTakenSeconds) ||
          0;

        setResult({
          score,
          total,
          correct,
          incorrect,
          unanswered,
          timeTaken: formatTime(timeTakenSeconds),
        });
      } catch (err: any) {
        console.error("Round 1 result error:", err);
        if (isMounted) {
          setError(
            err.message || "Unable to load your result."
          );
        }
      } finally {
        if (isMounted) setLoading(false);
      }
    };

    fetchResult();

    return () => {
      isMounted = false;
    };
  }, [navigate]);

  // Loading screen
  if (loading) {
    return (
      <main className="result-page">
        <CircuitBackground />

        <div className="result-overlay" />

        <div className="result-loading">
          <div className="result-icon">
            <Trophy size={30} />
          </div>

          <h2>Loading Result...</h2>

          <p>
            Please wait while we retrieve your Round 01
            performance.
          </p>
        </div>
      </main>
    );
  }

  // Error screen
  if (error || !result) {
    return (
      <main className="result-page">
        <CircuitBackground />

        <div className="result-overlay" />

        <div className="result-loading">
          <div className="result-icon">
            <XCircle size={30} />
          </div>

          <h2>Unable to Load Result</h2>

          <p>
            {error ||
              "Your result could not be retrieved. Please try again."}
          </p>

          <button
            className="result-primary"
            onClick={() => window.location.reload()}
          >
            TRY AGAIN
            <RotateCcw size={18} />
          </button>
        </div>
      </main>
    );
  }

  const percentage =
    result.total > 0
      ? Math.round((result.correct / result.total) * 100)
      : 0;

  return (
    <main className="result-page">
      <CircuitBackground />

      <div className="result-overlay" />

      {/* HEADER */}
      <header className="result-header">
        <div className="result-brand">
          <span className="brand-dot" />

          <div>
            <div className="brand-small">
              IEEE STUDENT BRANCH · GRAPHIC ERA UNIVERSITY
            </div>

            <div className="brand-name">
              KINDLE JUNIOR <b>5.0</b>
            </div>
          </div>
        </div>

        <div className="result-round">
          ROUND 01 · COMPLETED
        </div>
      </header>

      {/* MAIN */}
      <section className="result-container">
        {/* TOP */}
        <div className="result-heading">
          <div className="result-icon">
            <Trophy size={30} />
          </div>

          <div>
            <div className="result-eyebrow">
              ROUND 01 · PERFORMANCE SUMMARY
            </div>

            <h1>
              Challenge <span>Complete.</span>
            </h1>

            <p>
              Your Round 01 attempt has been recorded
              successfully.
            </p>
          </div>
        </div>

        {/* SCORE CARD */}
        <div className="score-card">
          <div className="score-left">
            <div className="score-label">
              YOUR SCORE
            </div>

            <div className="score-number">
              {result.score}
              <span>/{result.total}</span>
            </div>

            <div className="score-progress">
              <div
                className="score-progress-fill"
                style={{
                  width: `${percentage}%`,
                }}
              />
            </div>

            <div className="score-percentage">
              {percentage}% ACCURACY
            </div>
          </div>

          <div className="score-status">
            <div className="status-icon">
              <ShieldCheck size={24} />
            </div>

            <div>
              <strong>Round Submitted</strong>

              <p>
                Your responses have been securely recorded.
              </p>
            </div>
          </div>
        </div>

        {/* STAT CARDS */}
        <div className="result-stats">
          <div className="result-stat">
            <div className="stat-icon correct">
              <CheckCircle2 size={21} />
            </div>

            <div>
              <span>CORRECT</span>
              <strong>{result.correct}</strong>
            </div>
          </div>

          <div className="result-stat">
            <div className="stat-icon incorrect">
              <XCircle size={21} />
            </div>

            <div>
              <span>INCORRECT</span>
              <strong>{result.incorrect}</strong>
            </div>
          </div>

          <div className="result-stat">
            <div className="stat-icon unanswered">
              <CircleHelp size={21} />
            </div>

            <div>
              <span>UNANSWERED</span>
              <strong>{result.unanswered}</strong>
            </div>
          </div>

          <div className="result-stat">
            <div className="stat-icon time">
              <Clock3 size={21} />
            </div>

            <div>
              <span>TIME TAKEN</span>
              <strong>{result.timeTaken}</strong>
            </div>
          </div>
        </div>

        {/* INFO */}
        <div className="result-info">
          <div className="info-number">
            01
          </div>

          <div>
            <strong>Round 01 Recorded</strong>

            <p>
              Your performance has been saved. Keep this
              result for your reference before proceeding
              to the next stage.
            </p>
          </div>
        </div>

        {/* ACTIONS */}
        <div className="result-actions">
          <Link
            to="/"
            className="result-secondary"
          >
            <RotateCcw size={16} />
            BACK TO HOME
          </Link>

          <Link
            to="/round-2-guidelines"
            className="result-primary"
          >
            CONTINUE TO ROUND 02

            <ArrowRight size={18} />
          </Link>
        </div>

        {/* FOOTER */}
        <div className="result-footer">
          <span>ROUND 01</span>

          <i />

          <span>{result.total} QUESTIONS</span>

          <i />

          <span>60 MINUTES</span>
        </div>
      </section>
    </main>
  );
};

export default Round1Result;