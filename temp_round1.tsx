import { useEffect, useState, useCallback, useRef } from "react";
import {
  Clock3,
  ChevronRight,
  CheckCircle2,
  AlertTriangle,
  Maximize,
} from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";
import CircuitBackground from "../components/CircuitBackground";
import { pb } from "../lib/pb";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const TOTAL_QUESTIONS = 60;
const QUESTION_TIME = 60;

type APIQuestionResponse = {
  attempt_id: string;
  question_number: number;
  total_questions: number;
  question_id: number;
  domain: string;
  question: string;
  options: string[];
  server_deadline: string;
  question_deadline?: string;
};

const Round1 = () => {
  const navigate = useNavigate, useLocation();

  /* ---------------------------------------------------------
   * STATE
   * --------------------------------------------------------- */

  const [attemptId, setAttemptId] = useState<string | null>(null);

  const [needsLanguageSelection, setNeedsLanguageSelection] = useState(false);
  const [selectedLanguage, setSelectedLanguage] = useState<"C" | "Python" | null>(null);

  const [currentQuestion, setCurrentQuestion] =
    useState<APIQuestionResponse | null>(null);

  const [selectedAnswer, setSelectedAnswer] =
    useState<number | null>(null);

  const [timeLeft, setTimeLeft] =
    useState<number>(QUESTION_TIME);

  const [isSubmitting, setIsSubmitting] =
    useState(false);

  const [isExpired, setIsExpired] =
    useState(false);

  const [loading, setLoading] =
    useState(false);

  const [errorMessage, setErrorMessage] =
    useState("");

  /* ---------------------------------------------------------
   * FULLSCREEN
   * --------------------------------------------------------- */

  const [isFullscreen, setIsFullscreen] =
    useState(!!document.fullscreenElement);

  const [fullscreenRequested, setFullscreenRequested] =
    useState(false);

  /* ---------------------------------------------------------
   * REFS
   * --------------------------------------------------------- */

  const submittingRef =
    useRef(false);

  /*
   * Prevent duplicate POST /attempt calls.
   *
   * This is especially important with React StrictMode
   * during development.
   */
  const initStartedRef =
    useRef(false);

  /* ---------------------------------------------------------
   * FULLSCREEN CHANGE LISTENER
   * --------------------------------------------------------- */

  useEffect(() => {
    const handleFullscreenChange = () => {
      const fullscreen =
        !!document.fullscreenElement;

      setIsFullscreen(fullscreen);
    };

    document.addEventListener(
      "fullscreenchange",
      handleFullscreenChange
    );

    return () => {
      document.removeEventListener(
        "fullscreenchange",
        handleFullscreenChange
      );
    };
  }, []);

  /* ---------------------------------------------------------
   * ENTER FULLSCREEN
   * --------------------------------------------------------- */

  const enterFullscreen = async () => {
    try {
      if (!document.fullscreenElement) {
        await document.documentElement.requestFullscreen();
      }

      setFullscreenRequested(true);
    } catch (error) {
      console.error(
        "Fullscreen request failed:",
        error
      );

      setErrorMessage(
        "Fullscreen permission was not granted. Please click the button again."
      );
    }
  };

  /* ---------------------------------------------------------
   * 1. START / RESUME ATTEMPT
   *
   * IMPORTANT:
   * Attempt creation starts ONLY after fullscreen.
   * --------------------------------------------------------- */

  const initAttempt = useCallback(async (langToSubmit?: string) => {
    try {
      setLoading(true);
      setErrorMessage("");

      const token =
        pb.authStore.token;

      if (!token) {
        throw new Error(
          "Authentication required. Please login again."
        );
      }

      /* -----------------------------------------------------
       * FETCH COMPETITION
       * ----------------------------------------------------- */

      const compRes = await fetch(
        `${API_BASE_URL}/api/competitions`,
        {
          headers: {
            Authorization: token,
          },
        }
      );

      if (!compRes.ok) {
        throw new Error(
          `Failed to fetch competitions: ${compRes.status}`
        );
      }

      const competitions =
        await compRes.json();

      if (
        !competitions ||
        competitions.length === 0
      ) {
        throw new Error(
          "No active competitions found on server"
        );
      }

      const compId =
        competitions[0].id;

      /* -----------------------------------------------------
       * START / RESUME ATTEMPT
       * ----------------------------------------------------- */

      const attemptRes =
        await fetch(
          `${API_BASE_URL}/api/round1/attempt`,
          {
            method: "POST",

            headers: {
              "Content-Type":
                "application/json",

              Authorization: token,
            },

            body: JSON.stringify({
              competition_id:
                compId,
              language: langToSubmit || "",
            }),
          }
        );


      /* -----------------------------------------------------
       * EXPIRED / FORBIDDEN
       * ----------------------------------------------------- */

      if (attemptRes.status === 403) {
        setErrorMessage(
          "Your Round 1 attempt is not available or has expired."
        );

        setIsExpired(true);
        setLoading(false);

        return;
      }

      /* -----------------------------------------------------
       * LANGUAGE REQUIRED
       * ----------------------------------------------------- */
       
      if (attemptRes.status === 400) {
          const errText = await attemptRes.text();
          if (errText.includes("invalid language") || errText.includes("language_required")) {
              setNeedsLanguageSelection(true);
              setLoading(false);
              return;
          }
      }

      /* -----------------------------------------------------
       * OTHER ERRORS
       * ----------------------------------------------------- */

      if (!attemptRes.ok) {
        const errText =
          await attemptRes.text();

        throw new Error(
          `Failed to start attempt: ${attemptRes.status} ${errText}`
        );
      }

      const attemptData =
        await attemptRes.json();

      /* -----------------------------------------------------
       * ALREADY COMPLETED
       * ----------------------------------------------------- */

      if (
        attemptData.status ===
          "submitted" ||
        attemptData.status ===
          "expired"
      ) {
        const timeMins =
          Math.floor(
            (attemptData.time_taken_seconds ||
              0) / 60
          );

        const timeSecs =
          Math.floor(
            (attemptData.time_taken_seconds ||
              0) % 60
          );

        const answeredCount =
          Object.keys(
            attemptData.answers || {}
          ).length;

        const correctCount =
          attemptData.correct_answers ||
          0;

        navigate(
          "/round-1-result",
          {
            replace: true,

            state: {
              score:
                attemptData.score ||
                correctCount ||
                0,

              total:
                TOTAL_QUESTIONS,

              correct:
                correctCount,

              incorrect:
                answeredCount -
                correctCount,

              unanswered:
                TOTAL_QUESTIONS -
                answeredCount,

              timeTaken:
                `${String(
                  timeMins
                ).padStart(2, "0")}:` +
                `${String(
                  timeSecs
                ).padStart(2, "0")}`,
            },
          }
        );

        return;
      }

      /* -----------------------------------------------------
       * VALIDATE ATTEMPT
       * ----------------------------------------------------- */

      if (!attemptData.id) {
        throw new Error(
          "Invalid attempt response from server"
        );
      }

      /*
       * IMPORTANT:
       *
       * Backend should create NEW attempts with:
       *
       * current_question = 0
       *
       * because:
       *
       * 0 = Q1
       * 1 = Q2
       * 2 = Q3
       * ...
       * 59 = Q60
       */

      console.log(
        "Round 1 attempt:",
        attemptData.id
      );

      console.log(
        "Backend current question:",
        attemptData.current_question
      );

      setAttemptId(
        attemptData.id
      );

      setLoading(false);
    } catch (err: any) {
      console.error(
        "Round 1 initialization error:",
        err
      );

      setErrorMessage(
        err?.message ||
          "Unable to start Round 1."
      );

      setIsExpired(true);
      setLoading(false);
    }
  }, [navigate]);

  useEffect(() => {
    if (!isFullscreen) return;

    if (initStartedRef.current) return;
    initStartedRef.current = true;
    initAttempt();
  }, [isFullscreen, initAttempt]);

  /* ---------------------------------------------------------
   * 2. FETCH CURRENT QUESTION
   * --------------------------------------------------------- */

  const fetchQuestion =
    useCallback(
      async (aId: string) => {
        try {
          setLoading(true);

          const token =
            pb.authStore.token;

          if (!token) {
            throw new Error(
              "Authentication required."
            );
          }

          const res =
            await fetch(
              `${API_BASE_URL}/api/round1/question?attempt_id=${encodeURIComponent(
                aId
              )}`,
              {
                headers: {
                  Authorization:
                    token,
                },
              }
            );

          /* -------------------------------------------------
           * NO RELOAD
           * ------------------------------------------------- */

          if (res.status === 403) {
            setErrorMessage(
              "Your Round 1 session has expired or is no longer available."
            );

            setIsExpired(true);
            setLoading(false);

            return;
          }

          if (!res.ok) {
            const errText =
              await res.text();

            throw new Error(
              `Failed to fetch question: ${res.status} ${errText}`
            );
          }

          const data: APIQuestionResponse =
            await res.json();

          if (
            !data ||
            !data.question
          ) {
            throw new Error(
              "Invalid question received from server"
            );
          }

          console.log(
            "Question received:",
            data.question_number + 1,
            "/",
            data.total_questions
          );

          /* -------------------------------------------------
           * ATTEMPT VALIDATION
           * ------------------------------------------------- */

          if (
            data.attempt_id &&
            data.attempt_id !== aId
          ) {
            console.warn(
              "Attempt ID mismatch:",
              data.attempt_id,
              aId
            );
          }

          /* -------------------------------------------------
           * SET QUESTION
           * ------------------------------------------------- */

          setCurrentQuestion(data);

          setSelectedAnswer(null);

          /*
           * Initial value.
           *
           * Backend question_deadline will immediately
           * correct this.
           */
          setTimeLeft(
            QUESTION_TIME
          );

          setLoading(false);
        } catch (err: any) {
          console.error(
            "Question fetch error:",
            err
          );

          setErrorMessage(
            err?.message ||
              "Failed to fetch question."
          );

          setIsExpired(true);
          setLoading(false);
        }
      },
      []
    );

  /* ---------------------------------------------------------
   * FETCH QUESTION WHEN ATTEMPT IS READY
   * --------------------------------------------------------- */

  useEffect(() => {
    if (!attemptId) return;

    fetchQuestion(attemptId);
  }, [
    attemptId,
    fetchQuestion,
  ]);

  /* ---------------------------------------------------------
   * 3. EXACT 60 SECOND QUESTION TIMER
   * --------------------------------------------------------- */

  useEffect(() => {
    if (
      !currentQuestion ||
      isExpired
    ) {
      return;
    }

    const deadlineString =
      currentQuestion.question_deadline;

    /*
     * Fallback timer.
     */
    if (!deadlineString) {
      setTimeLeft(
        QUESTION_TIME
      );
    }

    const updateTimer = () => {
      if (submittingRef.current) {
        return;
      }

      /* -----------------------------------------------------
       * BACKEND AUTHORITATIVE DEADLINE
       * ----------------------------------------------------- */

      if (deadlineString) {
        const deadline =
          new Date(
            deadlineString
          ).getTime();

        const now =
          Date.now();

        const diff =
          Math.max(
            0,
            Math.ceil(
              (deadline - now) /
                1000
            )
          );

        setTimeLeft(diff);

        if (
          diff <= 0 &&
          !submittingRef.current
        ) {
          setTimeout(() => {
            handleSubmitAnswer();
          }, 0);
        }

        return;
      }

      /* -----------------------------------------------------
       * LOCAL FALLBACK
       * ----------------------------------------------------- */

      setTimeLeft(
        (previous) => {
          const next =
            previous - 1;

          if (
            next <= 0 &&
            !submittingRef.current
          ) {
            setTimeout(() => {
              handleSubmitAnswer();
            }, 0);

            return 0;
          }

          return next;
        }
      );
    };

    updateTimer();

    const interval =
      window.setInterval(
        updateTimer,
        1000
      );

    return () => {
      window.clearInterval(
        interval
      );
    };
  }, [
    currentQuestion,
    isExpired,
  ]);

  /* ---------------------------------------------------------
   * 4. SELECT ANSWER
   * --------------------------------------------------------- */

  const selectAnswer = (
    optionIndex: number
  ) => {
    if (timeLeft <= 0) return;
    if (isExpired) return;
    if (isSubmitting) return;

    setSelectedAnswer(
      optionIndex
    );
  };

  /* ---------------------------------------------------------
   * 5. SUBMIT ANSWER
   * --------------------------------------------------------- */

  async function handleSubmitAnswer() {
    if (!attemptId) return;
    if (!currentQuestion) return;

    if (
      submittingRef.current
    ) {
      return;
    }

    if (isExpired) return;

    submittingRef.current = true;

    setIsSubmitting(true);

    try {
      const answer =
        selectedAnswer !== null
          ? selectedAnswer
          : -1;

      const token =
        pb.authStore.token;

      const res =
        await fetch(
          `${API_BASE_URL}/api/round1/answer`,
          {
            method: "POST",

            headers: {
              "Content-Type":
                "application/json",

              Authorization:
                token,
            },

            body: JSON.stringify({
              attempt_id:
                attemptId,

              question_id:
                currentQuestion.question_id,

              answer,
            }),
          }
        );

      /* -----------------------------------------------------
       * SESSION EXPIRED
       * ----------------------------------------------------- */

      if (res.status === 403) {
        setErrorMessage(
          "Your Round 1 session has expired."
        );

        setIsExpired(true);

        return;
      }

      if (!res.ok) {
        const errText =
          await res.text();

        throw new Error(
          `Failed to submit answer: ${res.status} ${errText}`
        );
      }

      /* -----------------------------------------------------
       * CHECK LAST QUESTION
       * ----------------------------------------------------- */

      const isLastQuestion =
        currentQuestion.question_number >=
        currentQuestion.total_questions - 1;

      if (isLastQuestion) {
        navigate(
          "/round-1-result",
          {
            replace: true,
          }
        );

        return;
      }

      /* -----------------------------------------------------
       * NEXT QUESTION
       * ----------------------------------------------------- */

      await fetchQuestion(
        attemptId
      );
    } catch (err: any) {
      console.error(
        "Submit answer error:",
        err
      );

      setErrorMessage(
        err?.message ||
          "Failed to submit answer."
      );

      setIsExpired(true);
    } finally {
      submittingRef.current =
        false;

      setIsSubmitting(false);
    }
  }

  /* ---------------------------------------------------------
   * 6. FULLSCREEN GATE
   *
   * Round 1 will NOT initialize until fullscreen.
   * --------------------------------------------------------- */

  if (!isFullscreen) {
    return (
      <main
        className="round-page flex items-center justify-center text-white"
        style={{
          minHeight: "100vh",
          width: "100%",
          position: "relative",
          overflow: "hidden",
        }}
      >
        <CircuitBackground />

        <div
          className="z-10 text-center"
          style={{
            position: "relative",
            zIndex: 10,
            textAlign: "center",
            maxWidth: "520px",
            padding: "40px",
            margin: "auto",
          }}
        >
          <div
            style={{
              width: "72px",
              height: "72px",
              border: "1px solid rgba(50, 160, 255, 0.4)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              margin: "0 auto 24px",
              borderRadius: "12px",
              background:
                "rgba(10, 35, 65, 0.45)",
            }}
          >
            <Maximize
              size={30}
              color="#249cff"
            />
          </div>

          <div
            style={{
              fontSize: "11px",
              letterSpacing: "3px",
              color: "#249cff",
              marginBottom: "12px",
            }}
          >
            KINDLE JUNIOR 5.0
          </div>

          <h1
            style={{
              fontSize: "32px",
              fontWeight: 700,
              marginBottom: "14px",
            }}
          >
            Fullscreen Required
          </h1>

          <p
            style={{
              color: "#8a9bac",
              lineHeight: 1.7,
              marginBottom: "28px",
            }}
          >
            Round 1 can only be attempted
            in fullscreen mode. Please enter
            fullscreen to continue.
          </p>

          <button
            type="button"
            onClick={enterFullscreen}
            style={{
              display: "inline-flex",
              alignItems: "center",
              justifyContent: "center",
              gap: "10px",
              minWidth: "220px",
              height: "48px",
              padding: "0 24px",
              border:
                "1px solid rgba(36, 156, 255, 0.6)",
              background:
                "rgba(20, 90, 150, 0.22)",
              color: "#ffffff",
              cursor: "pointer",
              fontSize: "12px",
              fontWeight: 600,
              letterSpacing: "1.5px",
            }}
          >
            <Maximize size={17} />
            ENTER FULLSCREEN
          </button>

          {fullscreenRequested &&
            !isFullscreen && (
              <p
                style={{
                  marginTop: "18px",
                  fontSize: "12px",
                  color: "#f0a04b",
                }}
              >
                Please allow fullscreen mode
                in your browser.
              </p>
            )}
        </div>
      </main>
    );
  }

  /* ---------------------------------------------------------
   * 7. EXPIRED / ERROR SCREEN
   * --------------------------------------------------------- */

  if (isExpired) {
    return (
      <main
        className="quiz-page flex items-center justify-center text-white"
        style={{
          minHeight: "100vh",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <CircuitBackground />

        <div
          className="z-10 text-center"
          style={{
            textAlign: "center",
            zIndex: 10,
          }}
        >
          <AlertTriangle
            size={64}
            className="mx-auto mb-4 text-red-500"
            style={{
              marginBottom: "20px",
              color: "#f44336",
            }}
          />

          <h1
            className="text-3xl font-bold mb-2"
            style={{
              fontSize: "32px",
              marginBottom: "10px",
            }}
          >
            Round 1 Complete or Expired
          </h1>

          <p
            className="mb-6"
            style={{
              fontSize: "16px",
              color: "#8a9bac",
            }}
          >
            {errorMessage ||
              "Your attempt has been finalised."}
          </p>

          <button
            className="quiz-action-next confirm-submit px-6 py-2 bg-blue-600 rounded"
            style={{
              marginTop: "30px",
              width: "auto",
              padding: "0 20px",
            }}
            onClick={() =>
              navigate(
                "/round-1-result",
                {
                  replace: true,
                }
              )
            }
          >
            Go to Results
          </button>
        </div>
      </main>
    );
  }

  /* ---------------------------------------------------------
   * 7.5 LANGUAGE SELECTION SCREEN
   * --------------------------------------------------------- */

  if (needsLanguageSelection) {
    return (
      <main
        className="round-page flex items-center justify-center text-white"
        style={{
          minHeight: "100vh",
          width: "100%",
          position: "relative",
          overflow: "hidden",
        }}
      >
        <CircuitBackground />

        <div
          className="z-10 text-center"
          style={{
            position: "relative",
            zIndex: 10,
            textAlign: "center",
            maxWidth: "520px",
            padding: "40px",
            margin: "auto",
          }}
        >
          <div
            style={{
              fontSize: "11px",
              letterSpacing: "3px",
              color: "#249cff",
              marginBottom: "12px",
            }}
          >
            KINDLE JUNIOR 5.0
          </div>

          <h1
            style={{
              fontSize: "32px",
              fontWeight: 700,
              marginBottom: "14px",
            }}
          >
            Select Your Language
          </h1>

          <p
            style={{
              color: "#8a9bac",
              lineHeight: 1.7,
              marginBottom: "32px",
            }}
          >
            Please choose your preferred programming language to begin Round 1.
          </p>

          <div style={{ display: 'flex', gap: '20px', justifyContent: 'center' }}>
            <button
              type="button"
              onClick={() => {
                setSelectedLanguage("C");
                setNeedsLanguageSelection(false);
                initAttempt("C");
              }}
              style={{
                display: "inline-flex",
                alignItems: "center",
                justifyContent: "center",
                minWidth: "150px",
                height: "52px",
                padding: "0 24px",
                border: "1px solid rgba(36, 156, 255, 0.6)",
                background: "rgba(20, 90, 150, 0.22)",
                color: "#ffffff",
                cursor: "pointer",
                fontSize: "14px",
                fontWeight: 600,
                letterSpacing: "1.5px",
                borderRadius: "4px"
              }}
            >
              C PROGRAMMING
            </button>
            <button
              type="button"
              onClick={() => {
                setSelectedLanguage("Python");
                setNeedsLanguageSelection(false);
                initAttempt("Python");
              }}
              style={{
                display: "inline-flex",
                alignItems: "center",
                justifyContent: "center",
                minWidth: "150px",
                height: "52px",
                padding: "0 24px",
                border: "1px solid rgba(36, 156, 255, 0.6)",
                background: "rgba(20, 90, 150, 0.22)",
                color: "#ffffff",
                cursor: "pointer",
                fontSize: "14px",
                fontWeight: 600,
                letterSpacing: "1.5px",
                borderRadius: "4px"
              }}
            >
              PYTHON
            </button>
          </div>
        </div>
      </main>
    );
  }

  /* ---------------------------------------------------------
   * 8. LOADING
   * --------------------------------------------------------- */

  if (
    loading ||
    !currentQuestion
  ) {
    return (
      <main className="round-page flex items-center justify-center text-white">
        <CircuitBackground />

        <div
          className="z-10"
          style={{
            position: "relative",
            zIndex: 10,
          }}
        >
          Loading Question...
        </div>
      </main>
    );
  }

  /* ---------------------------------------------------------
   * 9. QUESTION DISPLAY
   * --------------------------------------------------------- */

  /*
   * Backend is ZERO-BASED:
   *
   * 0 → Q1
   * 1 → Q2
   * ...
   * 59 → Q60
   *
   * Therefore +1 here is CORRECT.
   */

  const questionNumber =
    currentQuestion.question_number + 1;

  const totalQuestions =
    currentQuestion.total_questions ||
    TOTAL_QUESTIONS;

  const timerDanger =
    timeLeft <= 10;

  const progress =
    Math.min(
      100,
      Math.max(
        0,
        ((QUESTION_TIME -
          timeLeft) /
          QUESTION_TIME) *
          100
      )
    );

  /* ---------------------------------------------------------
   * RENDER
   * --------------------------------------------------------- */

  return (
    <main className="round-page">
      <CircuitBackground />

      <div className="round-overlay" />

      {/* -----------------------------------------------------
          HEADER
      ----------------------------------------------------- */}

      <header className="round-header">
        <div className="round-brand">
          <span className="round-brand-dot" />

          <div>
            <div className="round-brand-small">
              KINDLE JUNIOR
            </div>

            <div className="round-brand-name">
              5.0
            </div>
          </div>
        </div>

        <div className="round-title">
          ROUND 01 · RAPID FIRE {selectedLanguage ? `· ${selectedLanguage.toUpperCase()}` : ''}
        </div>

        <div
          className={`round-timer ${
            timerDanger
              ? "timer-danger"
              : ""
          }`}
        >
          <Clock3 size={17} />

          <div>
            <span>TIME LEFT</span>

            <strong>
              00:
              {String(
                timeLeft
              ).padStart(2, "0")}
            </strong>
          </div>
        </div>
      </header>

      {/* -----------------------------------------------------
          MAIN QUESTION
      ----------------------------------------------------- */}

      <section className="round-container">
        <div className="question-meta">
          <div>
            <span className="question-label">
              QUESTION
            </span>

            <strong>
              {String(
                questionNumber
              ).padStart(2, "0")}

              <small>
                {" "}
                /{" "}
                {String(
                  totalQuestions
                ).padStart(2, "0")}
              </small>
            </strong>
          </div>

          <div className="question-language">
            {currentQuestion.domain?.toLowerCase() === "python" ? (
              <>
                <span className="language-icon python-icon">🐍</span>
                <span>Python</span>
              </>
            ) : currentQuestion.domain?.toLowerCase() === "c" ? (
              <>
                <span className="language-icon c-icon">C</span>
                <span>C</span>
              </>
            ) : (
              <span>{currentQuestion.domain}</span>
            )}
          </div>
        </div>

        {/* QUESTION NUMBER PROGRESS */}

        <div className="question-progress">
          <div
            style={{
              width: `${
                (questionNumber /
                  totalQuestions) *
                100
              }%`,
            }}
          />
        </div>

        {/* QUESTION CARD */}

        <div className="question-card">
          <div className="question-card-top">
            <span>
              CHALLENGE{" "}
              {questionNumber}
            </span>

            <span>
              60 SEC
            </span>
          </div>

          <h1>
            {currentQuestion.question}
          </h1>

          {/* OPTIONS */}

          <div className="options-container">
            {currentQuestion.options.map(
              (
                option,
                index
              ) => {
                const selected =
                  selectedAnswer ===
                  index;

                return (
                  <button
                    key={index}
                    type="button"
                    onClick={() =>
                      selectAnswer(
                        index
                      )
                    }
                    disabled={
                      timeLeft <= 0 ||
                      isSubmitting
                    }
                    className={`option-card ${
                      selected
                        ? "option-selected"
                        : ""
                    }`}
                  >
                    <span className="option-number">
                      {String.fromCharCode(
                        65 + index
                      )}
                    </span>

                    <span className="option-text">
                      {option}
                    </span>

                    {selected && (
                      <CheckCircle2
                        size={18}
                        className="option-check"
                      />
                    )}
                  </button>
                );
              }
            )}
          </div>

          {/* TIMER BAR */}

          <div className="question-timer-bar">
            <div
              className={
                timerDanger
                  ? "timer-progress-danger"
                  : ""
              }
              style={{
                width: `${progress}%`,
              }}
            />
          </div>

          {/* ACTION */}

          <div className="question-action">
            <div className="rapid-note">
              <AlertTriangle
                size={14}
              />

              <span>
                Time-out automatically
                moves to the next
                question.
              </span>
            </div>

            <button
              className="next-question-button"
              onClick={
                handleSubmitAnswer
              }
              disabled={
                isSubmitting ||
                timeLeft <= 0
              }
            >
              {questionNumber >=
              totalQuestions
                ? "SUBMIT ROUND"
                : "NEXT QUESTION"}

              <ChevronRight
                size={18}
              />
            </button>
          </div>
        </div>

        {/* FOOTER */}

        <div className="round-footer">
          <span>
            {questionNumber} OF{" "}
            {totalQuestions}{" "}
            QUESTIONS
          </span>

          <i />

          <span>
            60 SECONDS EACH
          </span>

          <i />

          <span>
            NO BACKTRACKING
          </span>
        </div>
      </section>

      {/* -----------------------------------------------------
          FULLSCREEN EXIT WARNING
          This appears immediately if participant presses Esc.
      ----------------------------------------------------- */}

      {!isFullscreen && (
        <div
          style={{
            position: "fixed",
            inset: 0,
            zIndex: 99999,
            background:
              "rgba(0, 0, 0, 0.96)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            textAlign: "center",
            color: "#fff",
          }}
        >
          <div>
            <Maximize
              size={48}
              color="#249cff"
              style={{
                margin: "0 auto 20px",
              }}
            />

            <h2
              style={{
                fontSize: "28px",
                marginBottom: "12px",
              }}
            >
              Fullscreen Required
            </h2>

            <p
              style={{
                color: "#8a9bac",
                marginBottom: "24px",
              }}
            >
              Please return to fullscreen
              to continue Round 1.
            </p>

            <button
              type="button"
              onClick={enterFullscreen}
              style={{
                padding:
                  "12px 28px",
                background:
                  "rgba(20, 90, 150, 0.3)",
                border:
                  "1px solid rgba(36, 156, 255, 0.6)",
                color: "#fff",
                cursor: "pointer",
                letterSpacing:
                  "1.5px",
                fontWeight: 600,
              }}
            >
              RETURN TO FULLSCREEN
            </button>
          </div>
        </div>
      )}
    </main>
  );
};

export default Round1;