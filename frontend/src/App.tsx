import * as React from "react";
import { CreateSmsDialog } from "./components/CreateSmsDialog";
import { SmsList } from "./components/SmsList";

function App() {
  return (
    <div className="flex items-center justify-center w-screen h-screen space-x-4 bg-gradient-to-r from-green-400 to-blue-500">
        <CreateSmsDialog />
        <SmsList />
    </div>
  );
}

export default App;
