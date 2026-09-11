import { useEffect, useMemo, useState, useRef } from "react";
import { AlertTriangle } from "lucide-react";
import Editor from "@monaco-editor/react";
import {
  Play,
  Send,
  Clock3,
  ChevronDown,
  CheckCircle2,
  XCircle,
  Terminal,
  Code2,
  FileText,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  ArrowDownCircle,
} from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";

import CircuitBackground from "../components/CircuitBackground";
import { pb } from "../lib/pb";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const ROUND_TOTAL_SECONDS = 70 * 60;

type Language = "C" | "Python";

type QuestionType = "CASE STUDY" | "DEBUGGING";

interface Question {
  id: number;
  type: QuestionType;
  title: string;
  tags: string[];
  statement: string[];
  input: string;
  output: string;
  constraints: string[];
  publicTestCases: {
    input: string;
    output: string;
  }[];
  starterCode: Record<Language, string>;
}

/* ==========================================================
   STARTER CODE LITERALS (DEBUGGING QUESTIONS - HINTS REMOVED)
========================================================== */

/* Q3: Grade & Attendance Eligibility Checker */
const q3StarterC = `#include <stdio.h>

int main() {
    int attendance, marks;
    scanf("%d %d", &attendance, &marks);
    
    if (attendance >= 75 || marks >= 40) {
        printf("Pass");
    } else if (attendance < 75 && marks >= 80) {
        printf("Special Permission");
    } else {
        printf("Fail");
    }
    return 0;
}`;

const q3StarterPy = `import sys

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue

    attendance, marks = map(int, line.split())

    if attendance >= 75 or marks >= 40:
        print("Pass")
    elif attendance < 75 and marks >= 80:
        print("Special Permission")
    else:
        print("Fail")
`;

/* Q4: Income Tax Slab Calculator */
const q4StarterC = `#include <stdio.h>

int main() {
    double income, tax = 0;
    if (scanf("%lf", &income) == 1) {
        if (income > 500000) {
            tax = income * 0.20;
        } else if (income > 250000) {
            tax = income * 0.05;
        } else {
            tax = 0;
        }
        printf("%.2f\\n", tax);
    }
    return 0;
}`;

const q4StarterPy = `import sys

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue

    income = float(line)
    tax = 0.0

    if income > 500000:
        tax = income * 0.20
    elif income > 250000:
        tax = income * 0.05
    else:
        tax = 0.0

    print(f"{tax:.2f}")
`;

/* Q5: Max of Three Numbers */
const q5StarterC = `#include <stdio.h>

void solve(int a, int b, int c) {
    int max = a;

    if (b > max) {
        max = b;
    }

    if (c < max) {
        max = c;
    }

    printf("%d\\n", max);
}

int main() {
    int a, b, c;

    while (scanf("%d %d %d", &a, &b, &c) == 3) {
        solve(a, b, c);
    }

    return 0;
}`;

const q5StarterPy = `import sys

def solve(a, b, c):
    max_val = a

    if b > max_val:
        max_val = b

    if c < max_val:
        max_val = c

    print(max_val)

if __name__ == "__main__":
    for line in sys.stdin:
        line = line.strip()

        if not line:
            continue

        parts = line.split()

        if len(parts) == 3:
            a, b, c = map(int, parts)
            solve(a, b, c)
`;

/* Q6: Electricity Bill Calculator */
const q6StarterC = `#include <stdio.h>

void solve(int units) {
    float bill = 0;

    if (units <= 100) {
        bill = units * 5;
    } else if (units <= 200) {
        bill = (100 * 5) + ((units - 100) * 7);
    } else {
        bill = (100 * 5) + (100 * 7) + (units * 10);
    }

    printf("%.2f\\n", bill);
}

int main() {
    int units;

    while (scanf("%d", &units) == 1) {
        solve(units);
    }

    return 0;
}`;

const q6StarterPy = `import sys

def solve(units):
    if units <= 100:
        bill = units * 5
    elif units <= 200:
        bill = (100 * 5) + ((units - 100) * 7)
    else:
        bill = (100 * 5) + (100 * 7) + (units * 10)

    print(f"{bill:.2f}")

if __name__ == "__main__":
    for line in sys.stdin:
        line = line.strip()

        if not line:
            continue

        units = int(line)
        solve(units)
`;

/* Q7: Valid Triangle & Type */
const q7StarterC = `#include <stdio.h>

void solve(int a, int b, int c) {
    if (a + b > c && a + c > b && b + c > a) {

        if (a == b && b == c) {
            printf("Equilateral\\n");

        } else if (a == b || b == c) {
            printf("Isosceles\\n");

        } else {
            printf("Scalene\\n");
        }

    } else {
        printf("Invalid\\n");
    }
}

int main() {
    int a, b, c;

    while (scanf("%d %d %d", &a, &b, &c) == 3) {
        solve(a, b, c);
    }

    return 0;
}`;

const q7StarterPy = `import sys

def solve(a, b, c):
    if (a + b > c) and (a + c > b) and (b + c > a):

        if a == b == c:
            print("Equilateral")

        elif a == b or b == c:
            print("Isosceles")

        else:
            print("Scalene")

    else:
        print("Invalid")

if __name__ == "__main__":
    for line in sys.stdin:
        line = line.strip()

        if not line:
            continue

        parts = line.split()

        if len(parts) == 3:
            a, b, c = map(int, parts)
            solve(a, b, c)
`;

/* ==========================================================
   ROUND 2 QUESTIONS (MATCHING GO BACKEND QUESTION BANK IDs 1-7)
========================================================== */

const round2Questions: Question[] = [
  /* ID 1: Utility Water Tariff Calculator */
  {
    id: 1,
    type: "CASE STUDY",
    title: "Utility Water Tariff Calculator",
    tags: ["CONDITIONALS", "ARITHMETIC", "BILLING"],
    statement: [
      "Calculate the water tariff based on total water consumption.",
      "0-100L is free.",
      "101-300L costs ₹2/L for units over 100.",
      "301-500L costs ₹400 plus ₹5/L for units over 300.",
      "Above 500L costs ₹1400 plus ₹8/L for units over 500, along with a ₹100 surcharge.",
    ],
    input: "Single integer litres.",
    output: "Print the total calculated water tariff.",
    constraints: ["Input is a non-negative integer."],
    publicTestCases: [
      { input: "80", output: "0" },
      { input: "200", output: "200" },
      { input: "400", output: "900" },
      { input: "600", output: "2300" },
    ],
    starterCode: {
      C: `#include <stdio.h>

int main() {
    int litres;
    scanf("%d", &litres);

    // Write your solution here

    return 0;
}`,
      Python: `import sys

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue

    litres = int(line)

    # Write your solution here
`,
    },
  },

  /* ID 2: Flight Baggage Allowance Auditor */
  {
    id: 2,
    type: "CASE STUDY",
    title: "Flight Baggage Allowance Auditor",
    tags: ["CONDITIONALS", "BILLING", "STRINGS"],
    statement: [
      "Calculate the excess baggage charge based on baggage weight and flight type.",
      "Domestic flights allow up to 15kg free, charge ₹500/kg for 16-25kg over the allowance, and charge ₹5000 plus ₹1000/kg over 25kg above that.",
      "International flights allow up to 25kg free and charge ₹1200/kg for every kilogram over 25kg.",
    ],
    input: "Weight integer followed by flight type string (DOMESTIC or INTERNATIONAL).",
    output: "Print the total calculated excess baggage charge.",
    constraints: [
      "0 ≤ weight ≤ 500",
      "Flight type is either DOMESTIC or INTERNATIONAL.",
    ],
    publicTestCases: [
      { input: "12 DOMESTIC", output: "0" },
      { input: "20 DOMESTIC", output: "2500" },
      { input: "28 DOMESTIC", output: "8000" },
      { input: "30 INTERNATIONAL", output: "6000" },
    ],
    starterCode: {
      C: `#include <stdio.h>
#include <string.h>

int main() {
    int weight;
    char type[20];

    scanf("%d %s", &weight, type);

    // Write your solution here

    return 0;
}`,
      Python: `weight, flight_type = input().split()
weight = int(weight)

# Write your solution here
`,
    },
  },

  /* ID 3: Grade & Attendance Eligibility Checker */
  {
    id: 3,
    type: "DEBUGGING",
    title: "Grade & Attendance Eligibility Checker",
    tags: ["CONDITIONALS", "LOGIC"],
    statement: [
      "Attendance >= 75% AND Marks >= 40 results in 'Pass'.",
      "Attendance < 75% BUT Marks >= 80 results in 'Special Permission'.",
      "Otherwise, results in 'Fail'. Fix the logical bug in the condition.",
    ],
    input: "Two space-separated integers: attendance and marks.",
    output: "Print 'Pass', 'Special Permission', or 'Fail'.",
    constraints: ["0 ≤ attendance ≤ 100", "0 ≤ marks ≤ 100"],
    publicTestCases: [
      { input: "80 50", output: "Pass" },
      { input: "60 85", output: "Special Permission" },
      { input: "60 50", output: "Fail" },
      { input: "75 40", output: "Pass" },
    ],
    starterCode: {
      C: q3StarterC,
      Python: q3StarterPy,
    },
  },

  /* ID 4: Income Tax Slab Calculator */
  {
    id: 4,
    type: "DEBUGGING",
    title: "Income Tax Slab Calculator",
    tags: ["CONDITIONALS", "MATH", "BILLING"],
    statement: [
      "Calculate income tax based on marginal slabs:",
      "0 to 2,50,000 = 0%",
      "2,50,001 to 5,00,000 = 5% on amount over 2,50,000",
      "Above 5,00,000 = ₹12,500 + 20% on amount over 5,00,000.",
      "Fix the flat-rate calculation bug.",
    ],
    input: "A single float value representing annual income.",
    output: "Print total calculated tax formatted to 2 decimal places.",
    constraints: ["0 ≤ income ≤ 10,00,00,000"],
    publicTestCases: [
      { input: "200000", output: "0.00" },
      { input: "400000", output: "7500.00" },
      { input: "600000", output: "32500.00" },
      { input: "500000", output: "12500.00" },
    ],
    starterCode: {
      C: q4StarterC,
      Python: q4StarterPy,
    },
  },

  /* ID 5: Max of Three Numbers */
  {
    id: 5,
    type: "DEBUGGING",
    title: "Max of Three Numbers",
    tags: ["CONDITIONALS", "COMPARISON"],
    statement: [
      "Find the maximum of three integers. Fix the comparison bug in the code.",
    ],
    input: "Three integers separated by spaces.",
    output: "Print the maximum integer value.",
    constraints: ["-10⁹ ≤ a, b, c ≤ 10⁹"],
    publicTestCases: [
      { input: "3 7 5", output: "7" },
      { input: "10 10 10", output: "10" },
      { input: "1 2 3", output: "3" },
      { input: "-1 -2 -3", output: "-1" },
    ],
    starterCode: {
      C: q5StarterC,
      Python: q5StarterPy,
    },
  },

  /* ID 6: Electricity Bill Calculator */
  {
    id: 6,
    type: "DEBUGGING",
    title: "Electricity Bill Calculator",
    tags: ["CONDITIONALS", "MATH"],
    statement: [
      "Calculate the electricity bill based on units consumed. Fix the calculation bug for high consumption.",
    ],
    input: "A single integer units.",
    output: "Print total bill formatted to 2 decimal places.",
    constraints: ["0 ≤ units ≤ 100000"],
    publicTestCases: [
      { input: "50", output: "250.00" },
      { input: "150", output: "850.00" },
      { input: "250", output: "2200.00" },
      { input: "100", output: "500.00" },
    ],
    starterCode: {
      C: q6StarterC,
      Python: q6StarterPy,
    },
  },

  /* ID 7: Valid Triangle & Type */
  {
    id: 7,
    type: "DEBUGGING",
    title: "Valid Triangle & Type",
    tags: ["CONDITIONALS", "GEOMETRY"],
    statement: [
      "Check if a triangle is valid and classify it as Equilateral, Isosceles, or Scalene. Fix the logic so that invalid side combinations and types are properly identified.",
    ],
    input: "Three side lengths a, b, and c.",
    output: "Print 'Equilateral', 'Isosceles', 'Scalene', or 'Invalid'.",
    constraints: ["1 ≤ a, b, c ≤ 10000"],
    publicTestCases: [
      { input: "5 5 5", output: "Equilateral" },
      { input: "1 2 10", output: "Invalid" },
      { input: "3 4 5", output: "Scalene" },
      { input: "5 5 8", output: "Isosceles" },
    ],
    starterCode: {
      C: q7StarterC,
      Python: q7StarterPy,
    },
  },
];

/* ==========================================================
   LANGUAGE CONFIG
========================================================== */

const languageConfig: Record<
  Language,
  { monaco: string }
> = {
  C: {
    monaco: "c",
  },
  Python: {
    monaco: "python",
  },
};

/* ==========================================================
   CODING ARENA
========================================================== */

const CodingArena = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const problemContentRef = useRef<HTMLDivElement>(null);

  const [attemptId, setAttemptId] = useState<string | null>(
    location.state?.attempt_id || null
  );

  const [currentQuestion, setCurrentQuestion] = useState(0);
  const [totalQuestions, setTotalQuestions] = useState(7);
  const [code, setCode] = useState("");
  const [questionDetails, setQuestionDetails] = useState<any>(null);

  const language =
    (questionDetails?.selected_language as Language) || "C";

  const [timeLeft, setTimeLeft] = useState(ROUND_TOTAL_SECONDS);
  const [activeTab, setActiveTab] = useState<"tests" | "output">("tests");
  const [isRunning, setIsRunning] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [loading, setLoading] = useState(true);
  const [solvedQuestions, setSolvedQuestions] = useState<number[]>([]);

  const [testResults, setTestResults] = useState<
    {
      name: string;
      passed: boolean;
      error?: string;
      expected?: string;
      actual?: string;
    }[]
  >([]);

  /* ==========================================================
     FULLSCREEN STATE
  ========================================================== */

  const [isFullscreen, setIsFullscreen] = useState(false);
  const [fullscreenRequired, setFullscreenRequired] = useState(true);

  const enterFullscreen = async () => {
    try {
      if (!document.fullscreenElement) {
        await document.documentElement.requestFullscreen();
      }
      setIsFullscreen(true);
      setFullscreenRequired(false);
    } catch (error) {
      console.error("Fullscreen request failed:", error);
      setIsFullscreen(false);
      setFullscreenRequired(true);
    }
  };

  useEffect(() => {
    const handleFullscreenChange = () => {
      const active = !!document.fullscreenElement;
      setIsFullscreen(active);
      setFullscreenRequired(!active);
    };

    document.addEventListener("fullscreenchange", handleFullscreenChange);
    handleFullscreenChange();

    return () => {
      document.removeEventListener("fullscreenchange", handleFullscreenChange);
    };
  }, []);

  /* ==========================================================
     TIMER
  ========================================================== */

  useEffect(() => {
    let timer: number | undefined;
    const localStartTime = Date.now();

    const updateLocalTimer = () => {
      const elapsedSeconds = Math.floor(
        (Date.now() - localStartTime) / 1000
      );
      const remainingSeconds = Math.max(
        0,
        ROUND_TOTAL_SECONDS - elapsedSeconds
      );

      setTimeLeft(remainingSeconds);

      if (remainingSeconds <= 0) {
        navigate("/round-2/time-up", { replace: true });
      }
    };

    const serverDeadline = questionDetails?.server_deadline;

    if (serverDeadline) {
      const updateServerTimer = () => {
        const deadline = new Date(serverDeadline).getTime();
        const now = Date.now();
        const diff = Math.max(0, Math.ceil((deadline - now) / 1000));

        setTimeLeft(diff);

        if (diff <= 0) {
          navigate("/round-2/time-up", { replace: true });
        }
      };

      updateServerTimer();
      timer = window.setInterval(updateServerTimer, 1000);

      return () => {
        if (timer !== undefined) window.clearInterval(timer);
      };
    }

    updateLocalTimer();
    timer = window.setInterval(updateLocalTimer, 1000);

    return () => {
      if (timer !== undefined) window.clearInterval(timer);
    };
  }, [questionDetails?.server_deadline, navigate]);

  /* ==========================================================
     STATUS CHECK
  ========================================================== */

  useEffect(() => {
    if (questionDetails?.status === "expired") {
      navigate("/round-2/time-up", { replace: true });
    }

    if (questionDetails?.status === "submitted") {
      navigate("/round-2/completed", { replace: true });
    }
  }, [questionDetails?.status, navigate]);

  const normalizeQuestionIndex = (
    backendQuestionNumber: number,
    total: number
  ) => {
    const number = Number(backendQuestionNumber);
    if (!Number.isFinite(number)) return 0;
    const index = number >= 1 ? number - 1 : 0;
    return Math.max(0, Math.min(index, Math.max(0, total - 1)));
  };

  /* ==========================================================
     FETCH STATE
  ========================================================== */

  useEffect(() => {
    let cancelled = false;

    const fetchState = async () => {
      try {
        setLoading(true);
        let currentAttemptId: string | null = attemptId;

        const compRes = await fetch(`${API_BASE_URL}/api/competitions`);
        if (!compRes.ok) throw new Error(`Failed to fetch competitions: ${compRes.status}`);

        const comps = await compRes.json();
        if (!comps || comps.length === 0) throw new Error("No active competition found.");

        const compId = comps[0].id;

        if (!currentAttemptId) {
          const attemptRes = await fetch(`${API_BASE_URL}/api/round2/attempt`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: pb.authStore.token,
            },
            body: JSON.stringify({
              competition_id: compId,
              initial_language: "C",
            }),
          });

          if (!attemptRes.ok) {
            const errorText = await attemptRes.text();
            throw new Error(`Failed to create Round 2 attempt: ${attemptRes.status} ${errorText}`);
          }

          const attempt = await attemptRes.json();
          if (!attempt?.id) throw new Error("Invalid Round 2 attempt response.");

          currentAttemptId = attempt.id;
          setAttemptId(currentAttemptId);
        }

        if (!currentAttemptId) throw new Error("Round 2 attempt ID is missing.");
        if (cancelled) return;

        const stateRes = await fetch(
          `${API_BASE_URL}/api/round2/state?attempt_id=${encodeURIComponent(currentAttemptId)}`,
          { headers: { Authorization: pb.authStore.token } }
        );

        if (!stateRes.ok) {
          const errorText = await stateRes.text();
          throw new Error(`Failed to fetch Round 2 state: ${stateRes.status} ${errorText}`);
        }

        const stateData = await stateRes.json();
        if (cancelled) return;

        if (stateData.status === "submitted") {
          navigate("/round-2/completed", { replace: true });
          return;
        }

        if (stateData.status === "expired") {
          navigate("/round-2/time-up", { replace: true });
          return;
        }

        const backendTotal = Number(stateData.total_questions);
        const safeTotal = Number.isFinite(backendTotal) && backendTotal > 0
          ? backendTotal
          : round2Questions.length;

        setTotalQuestions(safeTotal);

        const normalizedIndex = normalizeQuestionIndex(
          stateData.current_question,
          safeTotal
        );

        setCurrentQuestion(normalizedIndex);

        const solved: number[] = [];
        if (stateData.answers && typeof stateData.answers === "object") {
          for (const qId of Object.keys(stateData.answers)) {
            const parsed = Number(qId);
            if (Number.isFinite(parsed)) solved.push(parsed);
          }
        }
        setSolvedQuestions(solved);

        const qRes = await fetch(
          `${API_BASE_URL}/api/round2/question?attempt_id=${encodeURIComponent(currentAttemptId)}`,
          { headers: { Authorization: pb.authStore.token } }
        );

        if (!qRes.ok) {
          const errorText = await qRes.text();
          throw new Error(`Failed to fetch question: ${qRes.status} ${errorText}`);
        }

        const qData = await qRes.json();
        if (cancelled) return;

        setQuestionDetails(qData);

        const langToUse = (qData?.selected_language || "C") as Language;
        const frontendQuestion = round2Questions[normalizedIndex];

        if (!frontendQuestion) {
          throw new Error(`Question ${normalizedIndex + 1} is not available.`);
        }

        setCode(frontendQuestion.starterCode[langToUse]);
        setTestResults([]);
        setActiveTab("tests");
        setLoading(false);
      } catch (err: any) {
        if (cancelled) return;
        console.error("Round 2 state error:", err);
        setLoading(false);
        alert(err?.message || "Unable to load Round 2.");
      }
    };

    fetchState();
    return () => {
      cancelled = true;
    };
  }, [attemptId, navigate]);

  /* ==========================================================
     SCROLL & TIMER HELPERS
  ========================================================== */

  const formatTime = (seconds: number) => {
    const safeSeconds = Math.max(0, seconds);
    const minutes = Math.floor(safeSeconds / 60);
    const remainingSeconds = safeSeconds % 60;
    return `${String(minutes).padStart(2, "0")}:${String(remainingSeconds).padStart(2, "0")}`;
  };

  const handleScrollDown = () => {
    problemContentRef.current?.scrollBy({ top: 300, behavior: "smooth" });
  };

  const handleScrollUp = () => {
    problemContentRef.current?.scrollBy({ top: -300, behavior: "smooth" });
  };

  /* ==========================================================
     QUESTION NAVIGATION & SWITCHING
  ========================================================== */

  const switchQuestionTo = (index: number) => {
  if (index < 0 || index >= round2Questions.length) return;

  const targetQuestion = round2Questions[index];
  setCurrentQuestion(index);

  setQuestionDetails((prev: any) => ({
    ...prev,
    question_id: targetQuestion.id,
    current_question: index + 1,
  }));

  setCode(targetQuestion.starterCode[language] || targetQuestion.starterCode["C"]);
  setTestResults([]);
  setActiveTab("tests");

  if (problemContentRef.current) {
    problemContentRef.current.scrollTop = 0;
  }
};

  const loadNextQuestion = async () => {
    if (!attemptId || isSubmitting || timeLeft <= 0) return;

    setIsSubmitting(true);

    try {
      const activeQuestion = round2Questions[currentQuestion];

      const res = await fetch(`${API_BASE_URL}/api/round2/submit`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: pb.authStore.token,
        },
        body: JSON.stringify({
          attempt_id: attemptId,
          question_id: Number(activeQuestion.id),
          language,
          code,
        }),
      });

      const data = await res.json().catch(() => null);

      if (!res.ok) {
        if (data?.error === "time up") {
          navigate("/round-2/time-up", { replace: true });
          return;
        }
        throw new Error(data?.error || `Unable to move to next question: ${res.status}`);
      }

      if (
        data?.status === "submitted" ||
        data?.status === "completed" ||
        currentQuestion >= totalQuestions - 1
      ) {
        navigate("/round-2/completed", { replace: true });
        return;
      }

      const nextIndex = currentQuestion + 1;
      if (nextIndex < round2Questions.length) {
        switchQuestionTo(nextIndex);
      } else {
        navigate("/round-2/completed", { replace: true });
      }
    } catch (err: any) {
      console.error("Next question error:", err);
      alert("Unable to move to next question: " + (err?.message || "Unknown error"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const switchQuestion = (index: number) => {
    switchQuestionTo(index);
  };

  const handlePrevious = () => {
    if (currentQuestion > 0) {
      switchQuestionTo(currentQuestion - 1);
    }
  };

  const handleNext = () => {
    if (currentQuestion < round2Questions.length - 1) {
      switchQuestionTo(currentQuestion + 1);
    } else {
      loadNextQuestion();
    }
  };

  /* ==========================================================
     RUN CODE
  ========================================================== */

const handleRun = async () => {
  if (!attemptId || isRunning) return;

  const activeQuestion = round2Questions[currentQuestion];

  setActiveTab("tests");
  setTestResults([]);
  setIsRunning(true);

  try {
    const res = await fetch(`${API_BASE_URL}/api/round2/run`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: pb.authStore.token,
      },
      body: JSON.stringify({
        attempt_id: attemptId,
        question_id: Number(activeQuestion.id), // Ensure explicit number mapping
        language,
        code,
      }),
    });

    if (!res.ok) {
      const errData = await res.json().catch(() => null);
      throw new Error(errData?.error || `Run failed: ${res.status}`);
    }

    const data = await res.json();
    setTestResults(data.tests || []);
  } catch (err: any) {
    alert("Error running code: " + (err?.message || "Unknown error"));
  } finally {
    setIsRunning(false);
  }
};

  /* ==========================================================
     SUBMIT CODE
  ========================================================== */

  const handleGlobalSubmitCode = async () => {
    await loadNextQuestion();
  };

  const handleExit = () => {
    navigate("/round-2/completed");
  };

  const progress = useMemo(() => {
    if (totalQuestions <= 0) return 0;
    return Math.round(((currentQuestion + 1) / totalQuestions) * 100);
  }, [currentQuestion, totalQuestions]);

  if (fullscreenRequired || !isFullscreen) {
    return (
      <main
        className="coding-arena-page"
        style={{
          height: "100vh",
          width: "100vw",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          position: "relative",
          overflow: "hidden",
          background: "#02070d",
          color: "white",
        }}
      >
        <CircuitBackground />
        <div
          style={{
            position: "absolute",
            inset: 0,
            background: "rgba(0, 5, 12, 0.82)",
            zIndex: 1,
          }}
        />
        <div
          style={{
            position: "relative",
            zIndex: 10,
            width: "min(520px, 90%)",
            textAlign: "center",
            padding: "40px",
            border: "1px solid rgba(40, 150, 255, 0.25)",
            background: "rgba(5, 18, 30, 0.94)",
            boxShadow: "0 0 50px rgba(0, 130, 255, 0.12)",
          }}
        >
          <Code2 size={44} style={{ margin: "0 auto 20px", color: "#2196f3" }} />
          <h1 style={{ fontSize: "26px", fontWeight: 700, marginBottom: "12px" }}>
            Full Screen Required
          </h1>
          <p style={{ color: "#8fa3b8", lineHeight: 1.7, marginBottom: "28px" }}>
            Round 2 Coding Arena can only be accessed in full screen mode.
            <br />
            Please enter full screen to continue.
          </p>
          <button
            onClick={enterFullscreen}
            style={{
              width: "100%",
              padding: "14px 20px",
              border: "none",
              borderRadius: "6px",
              background: "#2196f3",
              color: "#ffffff",
              fontWeight: 700,
              cursor: "pointer",
              letterSpacing: "0.08em",
            }}
          >
            ENTER FULL SCREEN
          </button>
        </div>
      </main>
    );
  }

  if (loading) {
    return (
      <main
        className="coding-arena-page"
        style={{
          height: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <CircuitBackground />
        <div style={{ zIndex: 10, color: "white" }}>Loading Question...</div>
      </main>
    );
  }

  const frontendQ = round2Questions[currentQuestion];

  if (!frontendQ) {
    return (
      <main
        className="coding-arena-page"
        style={{
          height: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          color: "white",
        }}
      >
        <CircuitBackground />
        <div style={{ zIndex: 10, textAlign: "center" }}>
          <AlertTriangle size={48} style={{ color: "#f44336", margin: "0 auto 16px" }} />
          <h2>Question loading error</h2>
          <p style={{ color: "#8a9bac", marginTop: "10px" }}>
            Invalid question index selected.
          </p>
        </div>
      </main>
    );
  }

  return (
    <main
      className="coding-arena-page"
      style={{
        height: "100vh",
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
      }}
    >
      <CircuitBackground />
      <div className="arena-overlay" />

      {/* HEADER */}
      <header className="arena-header">
        <div className="arena-brand">
          <div className="arena-brand-mark">
            <Code2 size={20} />
          </div>
          <div>
            <div className="arena-brand-title">
              KINDLE JUNIOR <span>5.0</span>
            </div>
            <div className="arena-brand-subtitle">
              IEEE STUDENT BRANCH · GRAPHIC ERA UNIVERSITY
            </div>
          </div>
        </div>

        <div className="arena-center-title">
          <span className="round-label">ROUND 02</span>
          <span className="arena-divider">/</span>
          <span>CODING ARENA</span>
        </div>

        <div className="arena-header-right">
          <button className="exit-button" onClick={handleExit}>
            <XCircle size={16} />
            EXIT
          </button>

          <div className={`arena-timer ${timeLeft <= 300 ? "timer-warning" : ""}`}>
            <Clock3 size={17} />
            <div>
              <small>TIME LEFT</small>
              <strong>{formatTime(timeLeft)}</strong>
            </div>
          </div>
        </div>
      </header>

      {/* QUESTION NAVIGATION */}
      <div className="question-navigation">
        <div className="question-nav-inner">
          <div className="question-nav-title">CHALLENGES ({language})</div>

          <div className="question-buttons">
            {round2Questions.map((item, index) => {
              const active = index === currentQuestion;
              const solved = solvedQuestions.includes(item.id);

              return (
                <button
                  key={item.id}
                  className={`question-button ${active ? "question-active" : ""} ${
                    solved ? "question-solved" : ""
                  }`}
                  onClick={() => switchQuestion(index)}
                >
                  <span className="question-button-number">
                    {String(index + 1).padStart(2, "0")}
                  </span>

                  {solved && <CheckCircle2 size={12} />}
                </button>
              );
            })}
          </div>

          <div className="question-progress">
            <span>{currentQuestion + 1}</span> / {totalQuestions}
            <div className="progress-bar">
              <div style={{ width: `${progress}%` }} />
            </div>
          </div>
        </div>
      </div>

      {/* MAIN WORKSPACE */}
      <section
        className="arena-workspace"
        style={{
          flex: 1,
          display: "flex",
          minHeight: 0,
          overflow: "hidden",
        }}
      >
        {/* PROBLEM PANEL */}
        <aside
          className="problem-panel"
          style={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            height: "100%",
            overflow: "hidden",
          }}
        >
          <div className="problem-header">
            <div className="problem-number">
              QUESTION {String(currentQuestion + 1).padStart(2, "0")}
            </div>

            <div className="up-down-navigation" style={{ display: "flex", gap: "6px" }}>
              <button
                className="icon-action"
                title="Scroll Statement Up"
                onClick={handleScrollUp}
                style={{
                  cursor: "pointer",
                  display: "flex",
                  alignItems: "center",
                  padding: "4px 8px",
                  border: "1px solid rgba(255, 255, 255, 0.15)",
                  background: "rgba(255, 255, 255, 0.08)",
                  color: "#fff",
                }}
              >
                <ChevronUp size={16} />
              </button>

              <button
                className="icon-action"
                title="Scroll Statement Down"
                onClick={handleScrollDown}
                style={{
                  cursor: "pointer",
                  display: "flex",
                  alignItems: "center",
                  padding: "4px 8px",
                  border: "1px solid rgba(255, 255, 255, 0.15)",
                  background: "rgba(255, 255, 255, 0.08)",
                  color: "#fff",
                }}
              >
                <ChevronDown size={16} />
              </button>
            </div>
          </div>

          <div
            className="problem-content"
            ref={problemContentRef}
            style={{
              overflowY: "auto",
              flex: 1,
              paddingBottom: "2rem",
            }}
          >
            <div className="problem-category">
              CHALLENGE ({language})
            </div>

            <h1>{frontendQ.title}</h1>

            <div className="problem-tags">
              {frontendQ.tags.map((tag) => (
                <span key={tag}>{tag}</span>
              ))}
            </div>

            <section className="statement-section">
              <h3>
                <FileText size={15} /> PROBLEM STATEMENT
              </h3>
              {frontendQ.statement.map((paragraph, index) => (
                <p key={index}>{paragraph}</p>
              ))}
            </section>

            <section className="statement-section">
              <h3>INPUT</h3>
              <div className="code-block">
                <code>{frontendQ.input}</code>
              </div>
            </section>

            <section className="statement-section">
              <h3>OUTPUT</h3>
              <p>{frontendQ.output}</p>
            </section>

            <section className="statement-section">
              <h3>CONSTRAINTS</h3>
              <ul>
                {frontendQ.constraints.map((constraint) => (
                  <li key={constraint}>{constraint}</li>
                ))}
              </ul>
            </section>

            {frontendQ.publicTestCases?.length > 0 && (
              <section className="statement-section">
                <h3>VALIDATION TEST CASES</h3>
                <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                  {frontendQ.publicTestCases.map((tc, idx) => (
                    <div
                      key={idx}
                      style={{
                        padding: "8px 12px",
                        background: "rgba(255, 255, 255, 0.04)",
                        borderRadius: "4px",
                        fontSize: "0.85rem",
                      }}
                    >
                      <div>
                        <strong>Input:</strong> <code>{tc.input}</code>
                      </div>
                      <div>
                        <strong>Expected:</strong> <code>{tc.output}</code>
                      </div>
                    </div>
                  ))}
                </div>
              </section>
            )}

            <div style={{ marginTop: "1rem", marginBottom: "1rem" }}>
              <button
                onClick={handleScrollDown}
                style={{
                  width: "100%",
                  padding: "8px",
                  background: "rgba(255,255,255,0.05)",
                  border: "1px dashed rgba(255,255,255,0.2)",
                  color: "#aaa",
                  borderRadius: "4px",
                  cursor: "pointer",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  gap: "6px",
                  fontSize: "0.8rem",
                }}
              >
                <ArrowDownCircle size={14} /> SCROLL DOWN FOR MORE
              </button>
            </div>
          </div>
        </aside>

        {/* IDE PANEL */}
        <section
          className="ide-panel"
          style={{
            flex: 1.2,
            display: "flex",
            flexDirection: "column",
            height: "100%",
            overflow: "hidden",
          }}
        >
          <div className="ide-toolbar" />

          <div className="editor-container" style={{ flex: 1, minHeight: 0 }}>
            <Editor
              height="100%"
              language={languageConfig[language].monaco}
              value={code}
              onChange={(value) => setCode(value ?? "")}
              theme="vs-dark"
              options={{
                fontSize: 14,
                minimap: { enabled: true },
                automaticLayout: true,
                padding: { top: 18 },
                smoothScrolling: true,
                cursorSmoothCaretAnimation: "on",
                scrollBeyondLastLine: false,
                wordWrap: "on",
                lineNumbers: "on",
                renderLineHighlight: "all",
                tabSize: 4,
              }}
            />
          </div>

          {/* TERMINAL */}
          <div className="bottom-panel" style={{ height: "160px", flexShrink: 0 }}>
            <div className="bottom-tabs">
              <button
                className={activeTab === "tests" ? "bottom-tab active" : "bottom-tab"}
                onClick={() => setActiveTab("tests")}
              >
                <Terminal size={15} /> TEST CASES
              </button>

              <button
                className={activeTab === "output" ? "bottom-tab active" : "bottom-tab"}
                onClick={() => setActiveTab("output")}
              >
                OUTPUT
              </button>
            </div>

            <div className="test-content">
              {activeTab === "tests" && (
                <>
                  {testResults.length === 0 ? (
                    <div className="empty-tests">
                      <Terminal size={18} /> Run your code to test your solution.
                    </div>
                  ) : (
                    <div className="test-results">
                      {testResults.map((test, index) => (
                        <div
                          className="test-result"
                          key={index}
                          style={{
                            flexDirection: "column",
                            alignItems: "flex-start",
                            gap: "8px",
                          }}
                        >
                          <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                            {test.passed ? (
                              <CheckCircle2 size={17} className="passed-icon" />
                            ) : (
                              <XCircle size={17} className="failed-icon" />
                            )}
                            <span>Test Case {index + 1}</span>
                            <strong>{test.passed ? "PASSED" : "FAILED"}</strong>
                          </div>

                          {!test.passed &&
                            test.error !== "HIDDEN" &&
                            test.expected !== "HIDDEN" && (
                              <div
                                style={{
                                  fontSize: "12px",
                                  color: "#ccc",
                                  marginLeft: "27px",
                                }}
                              >
                                {test.error ? (
                                  <div style={{ color: "#ff6b6b" }}>
                                    Error: {test.error}
                                  </div>
                                ) : null}
                                <div>
                                  <strong>Expected:</strong> {test.expected}
                                </div>
                                <div>
                                  <strong>Actual:</strong> {test.actual}
                                </div>
                              </div>
                            )}
                        </div>
                      ))}
                    </div>
                  )}
                </>
              )}

              {activeTab === "output" && (
                <div className="output-content">
                  <span>Run your code to display output logs here.</span>
                </div>
              )}
            </div>
          </div>

          {/* ACTION BAR */}
          <div className="arena-action-bar" style={{ flexShrink: 0 }}>
            <div className="submission-info">
              <span>LANG</span>
              <strong>{language}</strong>
              <span>PROGRESS</span>
              <strong>
                {currentQuestion + 1}/{totalQuestions}
              </strong>
            </div>

            <div className="action-buttons">
              <button
                className="previous-button"
                onClick={handlePrevious}
                disabled={currentQuestion === 0}
              >
                <ChevronLeft size={16} /> PREV
              </button>

              <button
                className="run-button"
                onClick={handleRun}
                disabled={isRunning || isSubmitting}
              >
                <Play size={16} />
                {isRunning ? "RUNNING..." : "RUN CODE"}
              </button>

              <button
                className="next-button"
                onClick={handleNext}
                disabled={isSubmitting || isRunning || timeLeft <= 0}
              >
                {isSubmitting ? "SAVING..." : "NEXT"}
                <ChevronRight size={16} />
              </button>

              <button
                className="submit-button"
                onClick={handleGlobalSubmitCode}
                disabled={isSubmitting || timeLeft <= 0}
                style={{
                  backgroundColor: "#10b981",
                  color: "#ffffff",
                  display: "flex",
                  alignItems: "center",
                  gap: "8px",
                  padding: "0.6rem 1.2rem",
                  borderRadius: "6px",
                  fontWeight: 700,
                  marginLeft: "12px",
                  cursor: isSubmitting ? "not-allowed" : "pointer",
                  border: "none",
                  boxShadow: "0 0 12px rgba(16, 185, 129, 0.4)",
                  transition: "all 0.2s ease-in-out",
                  opacity: isSubmitting ? 0.6 : 1,
                }}
              >
                <Send size={16} />
                {isSubmitting ? "SUBMITTING..." : "SUBMIT CODE"}
              </button>
            </div>
          </div>
        </section>
      </section>
    </main>
  );
};

export default CodingArena;