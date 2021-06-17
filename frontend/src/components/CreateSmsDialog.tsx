import * as React from "react";
import { Dialog } from "@headlessui/react";
import { useForm } from "react-hook-form";
import Schema, { Type, string } from "computed-types";
import { computedTypesResolver } from "@hookform/resolvers/computed-types";
import type { SubmitHandler } from "react-hook-form";

type SmsRequest = {
  /**
   * Destination mobile number. 3-15 digits
   */
  destination?: string;
  /**
   * A collection of destination mobile numbers. 3-15 digits
   */
  destinations?: string[];
  /**
   * The SMS message. If longer than 160 characters (GSM) or 70 characters
   * (Unicode), splits into multiple SMS
   */
  message?: string;
  /**
   * Where the SMS appears to come from. 3-11 characters A-Za-z0-9 if
   * alphanumeric; 3-15 digits if numeric (if set, set sharedPool to null)
   */
  origin?: string;
};

const smsMultipartThreshold = 70;

const SmsSchema = Schema({
  destination: string
    .trim()
    .normalize()
    .between(3, 15, "Mobile number must be between 3 and 15 digits"),
  origin: string
    .trim()
    .normalize()
    .between(3, 15, "Mobile number must be between 3 and 15 digits")
    .optional(),
  message: string.trim().normalize().min(1).error("Messages field is required"),
});

type Sms = Type<typeof SmsSchema>;

type Status = "FLIGHT" | "SUCCESS" | "FAILURE" | "NONE";

export const CreateSmsDialog = () => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [status, setStatus] = React.useState<Status>("NONE");
  const {
    register,
    handleSubmit,
    watch,
    reset,
    formState: { errors },
  } = useForm<Sms>({
    resolver: computedTypesResolver(SmsSchema),
  });
  const watchMessage = watch("message", "");

  const onSubmit: SubmitHandler<Sms> = async (data) => {
    setStatus("FLIGHT");
    const response = await fetch(`http://${import.meta.env.SNOWPACK_PUBLIC_API_URL}/api/message`, {
      method: "POST",
      mode: "cors",
      cache: "no-cache",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      setStatus("FAILURE");
    }
    setStatus("SUCCESS");
  };

  const onReset = (close?: boolean) => {
    reset();
    setStatus("NONE");
    if (close) {
      setIsOpen(false);
    }
  };

  return (
    <>
      <Dialog
        open={isOpen}
        onClose={() => setIsOpen(false)}
        className="fixed inset-0 z-10 overflow-y-auto"
      >
        <Dialog.Overlay className="fixed inset-0 bg-black opacity-30" />
        <div className="flex items-center justify-center min-h-screen">
          <div className="z-50 max-w-sm mx-auto bg-white rounded w-72">
            <Dialog.Title className="text-title">Send SMS</Dialog.Title>
            <form onSubmit={handleSubmit(onSubmit)}>
              <div className="grid grid-cols-1 gap-6 mx-5 mb-5">
                <label className="block">
                  <span className="text-gray-700">From</span>
                  <input
                    {...register("origin")}
                    type="tel"
                    className="input-default"
                  />
                </label>
                {errors.origin?.message ? (
                  <p className="text-error">{errors.origin.message}</p>
                ) : null}
                <label className="block">
                  <span className="text-gray-700">To</span>
                  <input
                    {...register("destination")}
                    type="tel"
                    className="input-default"
                  />
                </label>
                {errors.destination?.message ? (
                  <p className="text-error">{errors.destination.message}</p>
                ) : null}
                <label className="block">
                  <span className="text-gray-700">Message</span>
                  <textarea
                    {...register("message")}
                    rows={3}
                    className="input-default"
                  />
                </label>
                <div className="flex">
                  <span className="flex-grow text-error">
                    {errors.message?.message}
                  </span>
                  <span className="text-sm text-secondary">
                    {watchMessage.length}/{smsMultipartThreshold}
                  </span>
                </div>
                {watchMessage.length > smsMultipartThreshold ? (
                  <p className="p-3 font-bold text-white rounded-md bg-gradient-to-r from-pink-500 to-yellow-500">
                    The sms will be sent as a multipart message
                  </p>
                ) : null}
              </div>
              <div className="flex justify-between p-4 space-x-4">
                <button onClick={() => onReset(true)} className="btn-blue">
                  Cancel
                </button>
                {status === "FLIGHT" || status === "NONE" ? (
                  <button
                    type="submit"
                    className={status === "FLIGHT" ? "btn-disabled" : "btn-blue"}
                    disabled={status === "FLIGHT"}
                  >
                    Submit
                  </button>
                ) : (
                  <button
                    type="button"
                    onClick={() => onReset()}
                    className={status === "SUCCESS" ? "btn-green" : "btn-red"}
                  >
                    {status === "SUCCESS" ? "Success" : "Failure"}
                  </button>
                )}
              </div>
            </form>
          </div>
        </div>
      </Dialog>
      <button
        type="button"
        className="btn-gray"
        onClick={() => {
          reset();
          setIsOpen(!isOpen);
        }}
      >
        Send SMS
      </button>
    </>
  );
};
