import { useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { CheckCircle2, AlertCircle } from "lucide-react";
import CircuitBackground from "../components/CircuitBackground";
import { pb } from "../lib/pb";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const Completed = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [status, setStatus] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchState = async () => {
      try {
        
        
        // If we don't have an ID, we'll just query the backend to see if there's an active one by attempting to start one, but we shouldn't create a new one.
        // Actually, just fetch it via the state endpoint. If no ID is passed, we might have to get it from the last known state.
        // The easiest way is to just fetch the user's attempt list from pocketbase directly.
        const compRes = await fetch(`${API_BASE_URL}/api/competitions`);
        const comps = await compRes.json();
        const compId = comps[0].id;
        
        // Find existing attempt for this user
        const records = await pb.collection("attempts_round2").getFullList({
            filter: `competition = "${compId}"`,
            sort: "-created",
        });
        
        if (records.length === 0) {
            navigate("/");
            return;
        }
        
        const attempt = records[0];
        
        if (attempt.status === "in_progress") {
            navigate("/round-2/coding-arena");
            return;
        }

        setStatus(attempt.status);
        setLoading(false);
      } catch (err) {
        console.error(err);
        navigate("/");
      }
    };
    fetchState();
  }, [navigate, location]);

  if (loading) {
      return (
        <main className="coding-arena-page" style={{ height: "100vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
          <CircuitBackground />
          <div style={{ color: "white" }}>Loading...</div>
        </main>
      );
  }

  const isTimeUp = status === "expired" || status === "expired";

  return (
    <main className="coding-arena-page" style={{ height: "100vh", display: "flex", alignItems: "center", justifyContent: "center", flexDirection: "column", fontFamily: "'Inter', sans-serif" }}>
      <CircuitBackground />
      <div style={{ 
          zIndex: 10, 
          color: "white", 
          textAlign: "center", 
          backgroundColor: "rgba(10, 15, 30, 0.8)", 
          padding: "50px", 
          borderRadius: "15px", 
          border: "1px solid rgba(255, 255, 255, 0.1)",
          backdropFilter: "blur(10px)",
          boxShadow: "0 0 30px rgba(0, 240, 255, 0.1)",
          maxWidth: "500px",
          width: "90%"
      }}>
        
        <div style={{ display: "flex", justifyContent: "center", marginBottom: "20px" }}>
            {isTimeUp ? (
                <AlertCircle size={60} color="#FF3366" style={{ filter: "drop-shadow(0 0 10px rgba(255,51,102,0.5))" }} />
            ) : (
                <CheckCircle2 size={60} color="#00f0ff" style={{ filter: "drop-shadow(0 0 10px rgba(0,240,255,0.5))" }} />
            )}
        </div>

        {isTimeUp ? (
            <h1 style={{ color: "#FF3366", fontSize: "2.5rem", marginBottom: "20px", textShadow: "0 0 10px rgba(255,51,102,0.5)", letterSpacing: "2px", fontWeight: "800" }}>TIME'S UP!</h1>
        ) : (
            <h1 style={{ color: "#00f0ff", fontSize: "2rem", marginBottom: "20px", textShadow: "0 0 10px rgba(0,240,255,0.5)", letterSpacing: "2px", fontWeight: "800", lineHeight: "1.3" }}>
                THANK YOU FOR YOUR<br/>PARTICIPATION
            </h1>
        )}

        {isTimeUp && (
            <p style={{ fontSize: "1.1rem", marginBottom: "15px", color: "#e0e0e0" }}>Thank you for your participation.</p>
        )}
        
        <p style={{ fontSize: "1.1rem", color: "#a0a0a0", marginBottom: "40px", marginTop: isTimeUp ? "0" : "15px" }}>Results will be announced soon.</p>
        
        <div style={{ marginTop: "30px", paddingTop: "20px", borderTop: "1px solid rgba(255,255,255,0.1)", letterSpacing: "4px", fontSize: "0.9rem", color: "#666", fontWeight: "600" }}>
            KINDLE JUNIOR 5.0
        </div>
      </div>
    </main>
  );
};

export default Completed;
