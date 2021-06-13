import React, { useEffect, useState } from "react";

function App() {
  const [count, setCount] = useState(0);
  useEffect(() => {
    const timer = setTimeout(() => setCount(count + 1), 1000);
    return () => clearTimeout(timer);
  }, [count, setCount]);

  return (
    <div className="container mx-auto">
      <header>
        <p className="font-sans">
          Page has been open for <code>{count}</code> seconds.
        </p>
      </header>
    </div>
  );
}

export default App;
