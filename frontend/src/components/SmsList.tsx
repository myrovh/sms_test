import { Dialog } from "@headlessui/react";
import * as React from "react";
import useSWR from "swr";

const fetcher = (url: string) => fetch(url).then((r) => r.json());

type Message = {
  id: number;
  outgoing_id: number;
  origin: string;
  destination: string;
  message: string;
};

type List = {
  limit: number;
  messages: Message[];
  ofset: number;
  total: number;
};

function useSmsList(fetch: boolean, limit: number) {
  const { data, error } = useSWR<List>(
    fetch ? `http://${import.meta.env.SNOWPACK_PUBLIC_API_URL}/api/message?limit=${limit}` : null,
    fetcher
  );

 

  return {
    list: data,
    isLoading: !error && !data,
    isError: error,
  };
}

export const SmsList = () => {
  const [isOpen, setIsOpen] = React.useState(false);
  const { list } = useSmsList(isOpen, 10);

  console.log(list);

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
            <Dialog.Title className="text-title">
              Last 10 Messages
            </Dialog.Title>
            <div className="grid grid-cols-1 py-5 divide-y divide-gray-100">
              {list?.messages.map((message) => (
                <div key={message.id} className="px-5 py-1">
                  <p className="text-secondary">
                    {message.origin}&emsp;&emsp;&emsp;🠖&emsp;&emsp;&emsp;
                    {message.destination}
                  </p>
                  <p>message: {message.message}</p>
                </div>
              ))}
            </div>
            <button
              type="button"
              onClick={() => setIsOpen(false)}
              className="m-5 btn-blue"
            >
              Close
            </button>
          </div>
        </div>
      </Dialog>
      <button
        type="button"
        className="btn-gray"
        onClick={() => setIsOpen(!isOpen)}
      >
        View SMS History
      </button>
    </>
  );
};
