import { React, useState, useEffect } from "react";
import ExpItem from "../ExpItem/ExpItem";
import Button from "../Button/Button";
import axios from "axios";

function ExpressionsList() {
  const [expressions, setExp] = useState([]);

  async function fetchExps() {
    const resp = await axios.get("http://localhost:8080/expressions/");
    console.log(resp)
    setExp(resp.data);
  }

  useEffect(() => {
    fetchExps();
  }, []);

  return (
    <div className="ExpList">
      <ul>
        {expressions.length !== 0 ? (
          <div>
            {expressions.map((exp) => (
              <li key={exp.id}>
                <ExpItem
                  expression={exp.expression}
                  status={exp.status}
                  start_date={exp.start_date}
                />
              </li>
            ))}{" "}
          </div>
        ) : (
          <div style={{ textAlign: "center" }}>
            <br></br>
            <h2>You have no expressions yet</h2>
            <h3>Let's count!</h3>
            <br></br>
            <a href="/">
              <Button>Calculator</Button>
            </a>
          </div>
        )}
      </ul>
    </div>
  );
}

export default ExpressionsList;