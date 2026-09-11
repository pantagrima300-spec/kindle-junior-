import React from "react";
import { CheckCircle2, AlertCircle } from "lucide-react";
import CircuitBackground from "../components/CircuitBackground";
import "./Round2Final.css";

interface Round2FinalPageProps {
  status: "completed" | "time_up";
}

const Round2FinalPage: React.FC<Round2FinalPageProps> = ({ status }) => {
  const isTimeUp = status === "time_up";

  return (
    <div className="round2-final-container">
      <div className="round2-final-background">
        <CircuitBackground />
      </div>

      <div className="round2-final-card">
        <div className="round2-final-header">
          KINDLE JUNIOR 5.0
        </div>

        <div className="round2-final-icon">
          {isTimeUp ? (
            <AlertCircle size={64} color="#ff3366" style={{ filter: "drop-shadow(0 0 12px rgba(255,51,102,0.6))" }} />
          ) : (
            <CheckCircle2 size={64} color="#00f0ff" style={{ filter: "drop-shadow(0 0 12px rgba(0,240,255,0.6))" }} />
          )}
        </div>

        {isTimeUp ? (
          <>
            <h1 className="round2-final-title time-up">TIME'S UP!</h1>
            <p className="round2-final-text">Thank you for your participation.</p>
            <p className="round2-final-subtext">Results will be announced soon.</p>
          </>
        ) : (
          <>
            <h1 className="round2-final-title completed">
              THANK YOU FOR YOUR<br />PARTICIPATION
            </h1>
            <p className="round2-final-subtext" style={{ marginTop: "1rem" }}>
              Results will be announced soon.
            </p>
          </>
        )}
      </div>
    </div>
  );
};

export default Round2FinalPage;
