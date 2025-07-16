"use client";

import IssueSelectItem from "./IssueSelectItem";

export default function MyIssue() {
  const issue_list = [
    [
      "Issue 1",
      "In Progress",
      "Urgent",
      "Alice",
      "2024-01-01",
      "2024-06-01",
      "issue",
    ],
    [
      "Frontend",
      "Completed",
      "Low Priority",
      "Bob",
      "2024-02-01",
      "2024-07-01",
      "issue",
    ],
  ];
  return (
    <div className="flex flex-col p-6 text-white w-9/10">
      <div className="flex flex-row items-center mb-4 gap-6">
        <span className="text-lg font-bold">My Issues</span>
        <button className="px-3 py-2 bg-blue-500 text-xs text-white rounded-md hover:bg-blue-700">
          <i className="fa-solid fa-plus text-[0.5rem] mr-2"></i>
          Create issue
        </button>
      </div>
      <div className="h-150 overflow-y-auto">
        {issue_list.map((issue, index) => (
          <IssueSelectItem
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
