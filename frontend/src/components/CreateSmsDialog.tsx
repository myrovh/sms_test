import { Dialog, Transition } from "@headlessui/react";
import * as React from "react";

export const CreateSmsDialog = () => {
  const [isOpen, setIsOpen] = React.useState(false);

  return (
    <>
        <Dialog
          open={isOpen}
          onClose={() => setIsOpen(false)}
          className="fixed inset-0 z-10 overflow-y-auto"
        >
          <Dialog.Overlay className="fixed inset-0 bg-black opacity-30" />
          <div className="flex items-center justify-center min-h-screen">
            <div className="z-50 max-w-sm mx-auto bg-white rounded">
              <Dialog.Title>Send SMS</Dialog.Title>
              <Dialog.Description>Some stuff</Dialog.Description>
              <input></input>
              <input></input>
            </div>
          </div>
        </Dialog>
      <button
        type="button"
        className="px-4 py-2 text-sm font-medium text-white bg-black rounded-md bg-opacity-20 hover:bg-opacity-30 focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75"
        onClick={() => setIsOpen(!isOpen)}
      >
        Send SMS
      </button>
    </>
  );
};
