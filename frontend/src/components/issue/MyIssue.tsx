"use client";

import IssueSelectItem from "./IssueSelectItem";

export default async function MyIssue() {
  const issue_list = await fetch("http://localhost:8080/api/issues", {
    method: "GET",
  })
    .then((res) => res.json())
    .catch((err) => {
      console.error("Failed to fetch issues:", err);
    });
  return (
    <div className="flex flex-col p-6 text-white w-4/5">
      <div className="flex flex-row items-center mb-4 gap-6">
        <span className="text-lg font-bold">My Issues</span>
        <button className="px-3 py-2 bg-blue-500 text-xs text-white rounded-md hover:bg-blue-700">
          <i className="fa-solid fa-plus text-[0.5rem] mr-2"></i>
          Create issue
        </button>
      </div>
      <div className="h-150 overflow-y-auto">
        {issue_list.map((issue: string[], index: number) => (
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
