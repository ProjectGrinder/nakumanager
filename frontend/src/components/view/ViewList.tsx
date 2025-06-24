"use client";

import ViewSelectItem from "./ViewSelectItem";

export default function ViewList() {
  const view_list = [
    ["View 1", "Alice Wonder", 3, "View 1"],
    ["View 2", "Bob Ross", 2, "View 2"],
    ["View 33", "John Johnson", 0, "View 3"],
  ];
  return (
    <div className="flex flex-col p-10 text-white w-3/5">
      <div className="flex-row text-2xl font-bold mb-6">
        <span>All Views</span>
        <i className="fa-solid fa-square-plus text-2xl ml-10"></i>
      </div>
      <div className="h-150 overflow-y-auto">
        {view_list.map((view, index) => (
          <ViewSelectItem
            name={view[0] as string}
            creator={view[1] as string}
            issue_num={view[2] as number}
            destination={view[3] as string}
            key={index}
          />
        ))}
      </div>
    </div>
  );
}
