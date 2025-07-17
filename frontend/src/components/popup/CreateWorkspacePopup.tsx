import { useState } from "react";

interface PopupProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (value: string) => void;
}

export default function CreateWorkspacePopup(props: PopupProps) {
  const [name, setName] = useState("");

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
          Enter new workspace name
        </span>
        <input
          type="text"
          className="w-full bg-gray-200 border-2 border-gray-400 text-gray-700 px-4 py-2 text-sm rounded-lg outline-none"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <div className="flex justify-end space-x-2 mt-8">
          <button
            onClick={() => {
              props.onSubmit(name);
              props.onClose();
            }}
            className="px-8 py-2 rounded-lg bg-blue-500 text-white hover:bg-blue-700 text-sm"
          >
            Create
          </button>
        </div>
      </div>
    </div>
  );
}
