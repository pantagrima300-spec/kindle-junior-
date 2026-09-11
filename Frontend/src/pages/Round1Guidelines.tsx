import {
  ArrowRight,
  Clock3,
  Code2,
  FileQuestion,
  ShieldCheck,
  LockKeyhole,
  CheckCircle2,
  XCircle,
} from "lucide-react";

import { useState } from "react";
import { useNavigate } from "react-router-dom";
import CircuitBackground from "../components/CircuitBackground";

const Round1Guidelines = () => {
  const navigate = useNavigate();

  const [accessCode, setAccessCode] = useState("");
  const [verified, setVerified] = useState(false);
  const [error, setError] = useState("");
  const [selectedLanguage, setSelectedLanguage] =
    useState<"c" | "python" | null>(null);

  const VALID_CODE = "201945";

  const verifyCode = () => {
    if (accessCode.trim() === VALID_CODE) {
      setVerified(true);
      setError("");
      setSelectedLanguage(null);
    } else {
      setVerified(false);
      setSelectedLanguage(null);
      setError("Invalid access code. Please enter the correct code.");
    }
  };

  const handleProceed = () => {
    if (verified && selectedLanguage) {
      navigate("/round-1", {
        state: {
          language: selectedLanguage,
        },
      });
    }
  };

  return (
    <main className="guidelines-page">
      <CircuitBackground />
      <div className="guidelines-overlay" />

      <header className="guidelines-top">
        <div className="guidelines-brand">
          <div className="guidelines-brand-mark">
            <ShieldCheck size={17} />
          </div>

          <div>
            <strong>IEEE STUDENT BRANCH</strong>
            <span>GRAPHIC ERA UNIVERSITY</span>
          </div>
        </div>

        <div className="guidelines-event">
          KINDLE JUNIOR <b>5.0</b>
        </div>
      </header>

      <section className="guidelines-container">
        <div className="guidelines-heading">
          <span className="guidelines-eyebrow">
            ROUND 01 · BEFORE YOU BEGIN
          </span>

          <h1>
            Know the
            <span> Rules.</span>
          </h1>

          <p>
            Read the following instructions carefully before
            starting your coding challenge.
          </p>
        </div>

        <div className="guidelines-card">
          <div className="guidelines-card-header">
            <div>
              <span>KINDLE JUNIOR 5.0</span>
              <h2>ROUND 01</h2>
            </div>

            <div className="round-number">01</div>
          </div>

          <div className="guidelines-grid">
            <div className="guideline-item">
              <div className="guideline-icon">
                <FileQuestion size={19} />
              </div>

              <div>
                <strong>60 Questions</strong>
                <p>You will answer 60 multiple-choice questions.</p>
              </div>
            </div>

            <div className="guideline-item">
              <div className="guideline-icon">
                <Clock3 size={19} />
              </div>

              <div>
                <strong>60 Minutes</strong>
                <p>
                  You have a total of 60 minutes to complete the round.
                </p>
              </div>
            </div>

            <div className="guideline-item">
              <div className="guideline-icon">
                <Code2 size={19} />
              </div>

              <div>
                <strong>Programming Domains</strong>
                <p>Questions will cover C, python and basic aptitude</p>
              </div>
            </div>

            <div className="guideline-item">
              <div className="guideline-icon">
                <Clock3 size={19} />
              </div>

              <div>
                <strong>45 seconds Per Question</strong>
                <p>
                  The round is designed around approximately 45 seconds
                  per question.
                </p>
              </div>
            </div>
          </div>

          <div className="guidelines-warning">
            <ShieldCheck size={17} />

            <div>
              <strong>Read Carefully</strong>
              <p>
                Once submitted, answers may not be changed. Manage your
                time wisely.
              </p>
            </div>
          </div>

          <div className="guidelines-access-code">
            <div className="guidelines-access-heading">
              <div className="guidelines-access-icon">
                <LockKeyhole size={19} />
              </div>

              <div>
                <span>SECURE ACCESS</span>
                <h3>Enter Code to Start Round 01</h3>
              </div>
            </div>

            <p className="guidelines-access-description">
              Enter the access code provided by the organisers to unlock
              Round 01.
            </p>

            <div className="guidelines-code-input-wrapper">
              <input
                type="text"
                value={accessCode}
                onChange={(e) => {
                  setAccessCode(e.target.value);
                  setError("");
                  setVerified(false);
                  setSelectedLanguage(null);
                }}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    verifyCode();
                  }
                }}
                placeholder="ENTER ACCESS CODE"
                maxLength={6}
                inputMode="numeric"
                autoComplete="off"
              />

              <button
                type="button"
                onClick={verifyCode}
                disabled={accessCode.length === 0}
                className="guidelines-verify-button"
              >
                VERIFY
              </button>
            </div>

            {verified && (
              <div className="guidelines-code-success">
                <CheckCircle2 size={17} />

                <span>
                  Access code verified. Round 01 is unlocked.
                </span>
              </div>
            )}

            {error && (
              <div className="guidelines-code-error">
                <XCircle size={17} />

                <span>{error}</span>
              </div>
            )}
          </div>

          {/* LANGUAGE SELECTION */}
          {verified && (
            <div className="round-language-selection">
              <div className="round-language-heading">
                <Code2 size={19} />

                <div>
                  <span>PROGRAMMING LANGUAGE</span>
                  <h3>Choose Your Language</h3>
                  <p>
                    Select the language you want to use for Round 01.
                  </p>
                </div>
              </div>

              <div className="round-language-options">
                <button
                  type="button"
                  onClick={() => setSelectedLanguage("c")}
                  className={`round-language-option ${
                    selectedLanguage === "c"
                      ? "round-language-option-selected"
                      : ""
                  }`}
                >
                  <span className="round-language-icon c-language-icon">
                    C
                  </span>

                  <span className="round-language-info">
                    <strong>C</strong>
                    <small>C Programming</small>
                  </span>

                  {selectedLanguage === "c" && (
                    <CheckCircle2 size={19} />
                  )}
                </button>

                <button
                  type="button"
                  onClick={() => setSelectedLanguage("python")}
                  className={`round-language-option ${
                    selectedLanguage === "python"
                      ? "round-language-option-selected"
                      : ""
                  }`}
                >
                  <span className="round-language-icon python-language-icon">
                    🐍
                  </span>

                  <span className="round-language-info">
                    <strong>Python</strong>
                    <small>Python Programming</small>
                  </span>

                  {selectedLanguage === "python" && (
                    <CheckCircle2 size={19} />
                  )}
                </button>
              </div>
            </div>
          )}

          <div className="guidelines-card-footer">
            <div className="question-domains">
              {selectedLanguage === "c"
                ? "C SELECTED"
                : selectedLanguage === "python"
                ? "PYTHON SELECTED"
                : "C · PYTHON"}
            </div>

            <button
              type="button"
              onClick={handleProceed}
              disabled={!verified || !selectedLanguage}
              className={`proceed-button ${
                !verified || !selectedLanguage
                  ? "proceed-button-disabled"
                  : ""
              }`}
            >
              <span>
                {!verified
                  ? "ROUND 01 LOCKED"
                  : !selectedLanguage
                  ? "CHOOSE LANGUAGE"
                  : "PROCEED TO ROUND 01"}
              </span>

              {verified && selectedLanguage ? (
                <ArrowRight size={17} />
              ) : (
                <LockKeyhole size={17} />
              )}
            </button>
          </div>
        </div>

        <div className="guidelines-meta">
          <span>PARTICIPANT ACCESS</span>
          <i />
          <span>ROUND 01</span>
          <i />
          <span>60 QUESTIONS</span>
          <i />
          <span>60 MINUTES</span>
        </div>
      </section>
    </main>
  );
};

export default Round1Guidelines;
