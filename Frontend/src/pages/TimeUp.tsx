
import { useNavigate } from "react-router-dom";

import "./CodingArena.css";

const TimeUp = () => {
  const navigate = useNavigate();

  return (
    <main className="coding-arena-page" style={{ height: "100vh", display: "flex", alignItems: "center", justifyContent: "center", flexDirection: "column" }}>
      
      <div style={{ zIndex: 10, color: "white", textAlign: "center", backgroundColor: "rgba(0,0,0,0.8)", padding: "40px", borderRadius: "10px", border: "1px solid #FF3366" }}>
        <h1 style={{ color: "#FF3366", fontSize: "3rem", marginBottom: "20px", textShadow: "0 0 10px rgba(255,51,102,0.5)" }}>TIME'S UP!</h1>
        <p style={{ fontSize: "1.2rem", marginBottom: "15px" }}>Your 70-minute Round 2 time has ended.</p>
        <p style={{ fontSize: "1rem", color: "#ccc", marginBottom: "30px" }}>Your progress and final recorded time have been saved.</p>
        <button 
          onClick={() => navigate("/")}
          style={{ padding: "10px 20px", backgroundColor: "#00f0ff", color: "black", border: "none", borderRadius: "5px", cursor: "pointer", fontWeight: "bold" }}
        >
          Return to Dashboard
        </button>
      </div>
    </main>
  );
};

export default TimeUp;
