import { useState } from "react";

interface PopupProps {
  open: boolean;
  name: string;
  onClose: () => void;
  onSubmit: (value: string) => void;
}

export default function DeleteWorkspacePopup(props: PopupProps) {
  if (!props.open) return null;

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div className="relative flex flex-col bg-gray-700 w-120 rounded-lg shadow-lg p-6">
        <button
          onClick={props.onClose}
          className="absolute top-3 right-3.5 text-gray-200 hover:text-white text-2xl font-bold px-2"
          aria-label="Close"
        >
          ×
        </button>
        <span className="text-base text-gray-200 font-semibold mb-3">
          Delete this workspace?
        </span>
        <span className="text-sm text-gray-300 mb-2">
          Everything related to "{props.name}" will be deleted.
        </span>
        <span className="text-sm text-gray-300 mb-2">
          This action cannot be undone.
        </span>
        <div className="flex justify-end space-x-2 mt-4">
          <button
            onClick={() => {
              props.onSubmit(props.name);
              props.onClose();
            }}
            className="px-8 py-2 rounded-lg bg-red-500 text-white hover:bg-blue-700 text-sm"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}
