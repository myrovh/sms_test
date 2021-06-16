import { Dialog } from "@headlessui/react";
import * as React from "react";
import useSWR from "swr";

const fetcher = (url: string) => fetch(url).then((r) => r.json());

export const SmsList = () => {
  const [isOpen, setIsOpen] = React.useState(false);
  const { data } = useSWR(
    isOpen ? "http://localhost:8080/api/message?limit=10" : null,
    fetcher
  );

  console.log(data);

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
            <Dialog.Title>Last 10 Messages</Dialog.Title>
            <Dialog.Description>Some stuff</Dialog.Description>
            <div>test</div>
            <input></input>
          </div>
        </div>
      </Dialog>
      <button
        type="button"
        className="px-4 py-2 text-sm font-medium text-white bg-black rounded-md bg-opacity-20 hover:bg-opacity-30 focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75"
        onClick={() => setIsOpen(!isOpen)}
      >
        View SMS History
      </button>
    </>
  );
};
