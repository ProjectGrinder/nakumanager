"use client";

import ViewSelectItem from "./ViewSelectItem";

export default async function ViewList() {
  const view_list = await fetch("http://localhost:8080/api/views/:id", {
    method: "GET",
  })
    .then((res) => res.json())
    .catch((err) => {
      console.error("Failed to fetch projects:", err);
    });
  return (
    <div className="flex flex-col p-6 text-white w-4/5">
      <div className="flex flex-row items-center mb-4 gap-6">
        <span className="text-lg font-bold">All Views</span>
        <button className="px-3 py-2 bg-blue-500 text-xs text-white rounded-md hover:bg-blue-700">
          <i className="fa-solid fa-plus text-[0.5rem] mr-2"></i>
          Create view
        </button>
      </div>
      <div className="h-150 overflow-y-auto">
        {view_list.map((view: any, index: number) => (
          <ViewSelectItem
            name={view[0] as string}
            creator={view[1] as string}
            issue_list={view[2] as string[]}
            destination={view[3] as string}
            key={index}
          />
        ))}
      </div>
    </div>
  );
}
