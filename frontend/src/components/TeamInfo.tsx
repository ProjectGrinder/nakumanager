"use client";

export default function TeamInfo() {
  const team_members = [
    ["Admin 1", "Project Manager"],
    ["Admin 2", "Designer"],
    ["Member 1", "Frontend"],
    ["Member 2", "Frontend"],
    ["Member 3", "Backend"],
  ];
  return (
    <div className="flex flex-col p-10 text-white w-1/2">
      <div className="flex-row text-2xl font-bold mb-4">
        <span>Team Members</span>
        <i className="fa-solid fa-gear text-2xl ml-10"></i>
      </div>
      <div className="flex flex-row justify-between mb-4">
        <div className="flex-column items-left justify-start w-100">
          <ul className="flex-column items-start justify-start h-120 mt-6 text-lg overflow-y-auto">
            {team_members.map((member, index) => (
              <li className="pt-4 pb-4" key={index}>
                {member[0]}
              </li>
            ))}
          </ul>
        </div>
        <div className="flex-column items-left justify-start w-100">
          <ul className="flex-column items-start justify-start h-120 mt-6 text-lg overflow-y-auto">
            {team_members.map((member, index) => (
              <li className="pt-4 pb-4" key={index}>
                {member[1]}
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
}
