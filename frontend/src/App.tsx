import * as React from "react";
import { CreateSmsDialog } from "./components/CreateSmsDialog";
import { SmsList } from "./components/SmsList";

function App() {
  return (
    <div className="flex w-screen h-screen bg-gradient-to-r from-green-400 to-blue-500">
      <div className="w-1/2 mx-auto mt-5 mb-5 bg-white rounded-md sm:w-5/6">
        <CreateSmsDialog />
        <SmsList />
      </div>
    </div>
  );
}

export default App;
