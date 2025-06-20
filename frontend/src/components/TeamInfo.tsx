"use client";

export default function TeamInfo() {
  const team_members = [
    ["Admin 1", "Project Manager"],
    ["Admin 2", "Admin"],
    ["Member 1", "Member"],
    ["Member 2", "Member"],
    ["Member 3", "Member"],
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
      <div className="flex flex-row align-center gap-10 mt-10 ml-10">
        <button className="bg-gray-100 text-gray-700 px-8 py-2 rounded hover:bg-gray-300">
          Projects
        </button>
        <button className="bg-gray-100 text-gray-700 px-8 py-2 rounded hover:bg-gray-300">
          Issues
        </button>
        <button className="bg-gray-100 text-gray-700 px-8 py-2 rounded hover:bg-gray-300">
          Views
        </button>
      </div>
    </div>
  );
}
