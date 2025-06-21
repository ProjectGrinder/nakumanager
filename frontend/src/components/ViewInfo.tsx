"use client";

import SelectableItem from "./SelectableItem";

export default function ViewInfo() {
  const issue_list = [
    [
      "Issue 1",
      "In Progress",
      "Urgent",
      "Alice",
      "2024-01-01",
      "2024-06-01",
      "Issue 1",
    ],
    [
      "Frontend",
      "Completed",
      "Low Priority",
      "Bob",
      "2024-02-01",
      "2024-07-01",
      "Issue 2",
    ],
  ];
  return (
    <div className="flex flex-col p-10 text-white w-9/10">
      <div className="flex-row text-2xl font-bold mb-8">
        <span>View 1</span>
        <i className="fa-solid fa-gear text-2xl ml-10"></i>
      </div>
      <div className="text-xl font-normal mb-6">
        <span>Creator: Alice Wonder</span>
      </div>
      <div className="h-150 overflow-y-auto">
        {issue_list.map((issue, index) => (
          <SelectableItem
            name={issue[0]}
            status={issue[1]}
            priority={issue[2]}
            assigned={issue[3]}
            startDate={issue[4]}
            endDate={issue[5]}
            destination={issue[6]}
            key={index}
          />
        ))}
      </div>
    </div>
  );
}
